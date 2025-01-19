package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	ROUND_32 = "32강"
	ROUND_16 = "16강"
	ROUND_8  = "8강"
	ROUND_4  = "4강"
	FINAL    = "결승"
)

// 한 게임의 전체 진행 상태
type Game struct {
	userName        string
	totalCandidates []*Candidate
	currentRound    *Round
	status          string     // "ready", "in_progress", "finished" 등 게임 상태
	finalWinner     *Candidate // 최종 우승자
}

// 현재 라운드 정보
type Round struct {
	candidates    []*Candidate // 현재 라운드의 후보들
	roundNumber   string       // ROUND_32, ROUND_16 등
	totalMatches  int          // 라운드의 총 매치 수
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
		totalCandidates: make([]*Candidate, 0, 32),
		status:          "ready",
	}
}

// 게임 초기화
func (g *Game) InitGame(username string) {
	g.userName = username

	// 후보자 목록 가져오기
	candidates := g.pullCandidates()
	if len(candidates) < 32 { // 32명 미만인지 체크
		log.Printf("후보자가 충분하지 않습니다: %d명", len(candidates))
		return
	}

	// 전체 후보 설정 (정확히 32명만 사용)
	g.totalCandidates = candidates[:32] // 32명만 선택

	// 32강 라운드 초기화
	g.currentRound = &Round{
		candidates:   g.totalCandidates, // 32명의 후보
		roundNumber:  ROUND_32,
		totalMatches: 16, // 32명이 대결하려면 16번의 매치 필요
		winners:      make([]*Candidate, 0, 16),
		currentBattle: &Battle{ // 첫 Battle 초기화
			MatchNumber: 0,
		},
	}

	g.status = "in_progress"
}

// 대결 시작
func (g *Game) StartBattle() map[string]interface{} {
	r := g.currentRound
	// 아직 대결하지 않은 후보들 중에서 랜덤하게 2명 선택
	remainingCandidates := r.getRemainingCandidates()
	if len(remainingCandidates) < 2 {
		log.Printf("남은 후보가 충분하지 않습니다: %d명", len(remainingCandidates))
		return nil
	}

	battle, err := r.getRandomPair(remainingCandidates)
	if err != nil {
		return nil
	}

	// 대결 중인 후보들은 isUsed 체크
	battle.Candidate1.isUsed = true
	battle.Candidate2.isUsed = true

	r.currentBattle = battle

	response := g.makeGameResponse()
	return response
}

// 아직 대결하지 않은 후보들 가져오기
func (r *Round) getRemainingCandidates() []*Candidate {
	b := r.currentBattle
	if b.MatchNumber == 0 {
		// 라운드 첫 시작이면 전체 후보 초기화
		r.resetCandidatesStatus()
		// 라운드 전체 후보 반환
		return r.candidates
	}

	// 이미 대결한 후보들 제외
	var remaining []*Candidate
	for _, c := range r.candidates {
		if !c.isUsed {
			remaining = append(remaining, c)
		}
	}
	return remaining
}

// 선택 처리
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

	// 승자를 winners에 추가하고 매치 수 증가
	r.winners = append(r.winners, winner)
	loser.FailStage = r.roundNumber
	b.MatchNumber++ // 매치 수를 여기서 증가
	log.Printf("%s %d / %d 에서 승자 : %s", r.roundNumber, b.MatchNumber, r.totalMatches, winner.Name)

	// 현재 라운드가 끝났는지 확인
	if b.MatchNumber >= r.totalMatches {
		// 다음 라운드로 진행
		if err := g.moveToNextRound(r.winners); err != nil {
			return nil, err
		}
	}

	// 다음 배틀 시작
	return g.StartBattle(), nil
}

func (g *Game) moveToNextRound(winners []*Candidate) error {
	switch g.currentRound.roundNumber {
	case ROUND_32:
		g.currentRound = &Round{
			roundNumber:  ROUND_16,
			totalMatches: 8,
			winners:      winners,
			currentBattle: &Battle{
				MatchNumber: 0,
			},
		}
	case ROUND_16:
		g.currentRound = &Round{
			roundNumber:  ROUND_8,
			totalMatches: 4,
			winners:      winners,
			currentBattle: &Battle{
				MatchNumber: 0,
			},
		}
	case ROUND_8:
		g.currentRound = &Round{
			roundNumber:  ROUND_4,
			totalMatches: 2,
			winners:      winners,
			currentBattle: &Battle{
				MatchNumber: 0,
			},
		}
	case ROUND_4:
		g.currentRound = &Round{
			roundNumber:  FINAL,
			totalMatches: 1,
			winners:      winners,
			currentBattle: &Battle{
				MatchNumber: 0,
			},
		}
	case FINAL:
		g.status = "finished"
		g.finalWinner = g.currentRound.winners[0]
		return nil
	}

	// 새 라운드의 첫 배틀 시작
	if g.status != "finished" {
		g.StartBattle()
	}
	return nil
}

func (g *Game) makeGameResponse() map[string]interface{} {
	r := g.currentRound
	if r == nil {
		log.Printf("게임 상태 이상: currentRound is nil")
		return nil
	}
	b := r.currentBattle

	return map[string]interface{}{
		"username":      g.userName,
		"currentRound":  r.roundNumber,
		"matchNumber":   b.MatchNumber,
		"totalMatches":  r.totalMatches,
		"currentBattle": b,
		"status":        g.status,
	}
}

func (r *Round) getRandomPair(candidates []*Candidate) (*Battle, error) {
	if len(candidates) < 2 {
		return nil, fmt.Errorf("후보가 부족합니다: %d명", len(candidates))
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 두 개의 서로 다른 인덱스를 한번에 선택
	idx1 := random.Intn(len(candidates))
	idx2 := random.Intn(len(candidates) - 1)

	// idx2가 idx1 이상이면 하나 증가시켜서 중복 방지
	if idx2 >= idx1 {
		idx2++
	}

	battle := &Battle{
		Candidate1:  candidates[idx1],
		Candidate2:  candidates[idx2],
		MatchNumber: r.currentBattle.MatchNumber,
	}

	return battle, nil
}
