package service

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	ROUND_32 = "ROUND_32"
	ROUND_16 = "ROUND_16"
	ROUND_8  = "ROUND_8"
	ROUND_4  = "ROUND_4"
	FINAL    = "FINAL"
)

// 한 게임의 전체 진행 상태
type Game struct {
	userName     string
	candidates   []*Candidate
	currentRound *Round
	status       string     // "ready", "in_progress", "finished" 등 게임 상태
	finalWinner  *Candidate // 최종 우승자
}

// 현재 라운드 정보
type Round struct {
	roundNumber   string       // ROUND_32, ROUND_16 등
	totalMatches  int          // 현재 라운드의 총 매치 수
	matchesPlayed int          // 현재까지 진행된 매치 수
	winners       []*Candidate // 현재 라운드를 통과한 후보들
	currentBattle *Battle      // 현재 대결 중인 두 후보
}

// 현재 진행중인 1:1 대결
type Battle struct {
	Candidate1  *Candidate `json:"candidate1"`       // 왼쪽 참가자
	Candidate2  *Candidate `json:"candidate2"`       // 오른쪽 참가자
	MatchNumber int        `json:"matchNumber"`      // 현재 라운드에서 몇 번째 매치인지
	Winner      *Candidate `json:"winner,omitempty"` // 승자 (결정된 경우)
}

func NewGame() *Game {
	return &Game{
		candidates: make([]*Candidate, 0, 32),
		status:     "ready",
	}
}

func (g *Game) InitGame(username string) map[string]interface{} {
	g.userName = username
	g.candidates = g.pullCandidates()

	// 32강 라운드 초기화
	g.currentRound = &Round{
		roundNumber:   ROUND_32,
		totalMatches:  16,
		matchesPlayed: 0,
		winners:       make([]*Candidate, 0, 16),
	}

	g.status = "in_progress"
	g.startBattle(g.currentRound)

	return g.makeGameResponse()
}

func (g *Game) startBattle(r *Round) *Round {
	// 아직 대결하지 않은 후보들 중에서 랜덤하게 2명 선택
	remainingCandidates := g.getRemainingCandidates()
	pair := g.getRandomPair(remainingCandidates)

	// 대결 중인 후보들은 isUsed 체크
	pair[0].isUsed = true
	pair[1].isUsed = true

	// 현재 배틀 정보 설정
	r.currentBattle = &Battle{
		Candidate1:  pair[0],
		Candidate2:  pair[1],
		MatchNumber: r.matchesPlayed + 1, // 현재 라운드의 몇 번째 매치인지
	}

	return r
}

// 아직 대결하지 않은 후보들 가져오기
func (g *Game) getRemainingCandidates() []*Candidate {
	if g.currentRound.matchesPlayed == 0 {
		// 라운드 첫 시작이면 전체 후보 반환
		g.resetCandidatesStatus()
		return g.candidates
	}

	// 이미 대결한 후보들 제외
	var remaining []*Candidate
	for _, c := range g.candidates {
		if !c.isUsed {
			remaining = append(remaining, c)
		}
	}
	return remaining
}

func (g *Game) ProcessSelection(selectedID int) (map[string]interface{}, error) {
	r := g.currentRound
	b := r.currentBattle
	// 현재 진행 중인 배틀이 없는 경우
	if b == nil {
		return nil, fmt.Errorf("진행 중인 대결이 없습니다")
	}

	// 선택된 후보가 현재 대결 중인 후보가 맞는지 확인
	var winner, loser *Candidate
	if b.Candidate1.ID == selectedID {
		winner = b.Candidate1
		loser = b.Candidate2
	} else if b.Candidate2.ID == selectedID {
		winner = b.Candidate2
		loser = b.Candidate1
	} else {
		return nil, fmt.Errorf("잘못된 선택입니다")
	}

	r.winners = append(r.winners, winner)
	r.matchesPlayed++
	loser.FailStage = r.roundNumber

	// 현재 라운드가 끝났는지 확인
	if g.currentRound.matchesPlayed >= g.currentRound.totalMatches {
		// 다음 라운드로 진행
		if err := g.moveToNextRound(r.winners); err != nil {
			return nil, err
		}
	} else {
		// 다음 배틀 시작
		g.startBattle(g.currentRound)
	}

	return g.makeGameResponse(), nil
}

func (g *Game) moveToNextRound(winners []*Candidate) error {
	switch g.currentRound.roundNumber {
	case ROUND_32:
		g.currentRound = &Round{
			roundNumber:   ROUND_16,
			totalMatches:  8,
			matchesPlayed: 0,
			winners:       winners,
		}
	case ROUND_16:
		g.currentRound = &Round{
			roundNumber:   ROUND_8,
			totalMatches:  4,
			matchesPlayed: 0,
			winners:       winners,
		}
	case ROUND_8:
		g.currentRound = &Round{
			roundNumber:   ROUND_4,
			totalMatches:  2,
			matchesPlayed: 0,
			winners:       winners,
		}
	case ROUND_4:
		g.currentRound = &Round{
			roundNumber:   FINAL,
			totalMatches:  1,
			matchesPlayed: 0,
			winners:       winners,
		}
	case FINAL:
		g.status = "finished"
		g.finalWinner = g.currentRound.winners[0]
		return nil
	}

	// 새 라운드의 첫 배틀 시작
	g.startBattle(g.currentRound)

	return nil
}

func (g *Game) makeGameResponse() map[string]interface{} {
	return map[string]interface{}{
		"username":      g.userName,
		"currentRound":  g.currentRound.roundNumber,
		"matchNumber":   g.currentRound.matchesPlayed + 1,
		"totalMatches":  g.currentRound.totalMatches,
		"currentBattle": g.currentRound.currentBattle,
	}
}

func (g *Game) getRandomPair(candidates []*Candidate) []*Candidate {
	if len(candidates) < 2 {
		return candidates
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 두 개의 서로 다른 인덱스를 한번에 선택
	idx1 := r.Intn(len(candidates))
	idx2 := r.Intn(len(candidates) - 1)

	// idx2가 idx1 이상이면 하나 증가시켜서 중복 방지
	if idx2 >= idx1 {
		idx2++
	}

	return []*Candidate{candidates[idx1], candidates[idx2]}
}
