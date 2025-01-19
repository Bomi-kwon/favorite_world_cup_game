package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Candidate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	FailStage string `json:"failStage"`
	isUsed    bool   // JSON에서 제외
}

// 후보자 전체 목록 가져오기
func (g *Game) pullCandidates() []*Candidate {
	candidates, err := g.loadCelebrityList()
	if err != nil {
		log.Printf("후보자 목록 로드 실패: %v", err)
		return nil
	}

	// 32개 이하면 그대로 반환
	if len(candidates) <= 32 {
		return candidates
	}

	// 셔플 후 32명만 선택
	g.shuffleCandidates(candidates)
	return candidates[:32]
}

func (g *Game) loadCelebrityList() ([]*Candidate, error) {
	file, err := os.Open("static/celebrities.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var candidates []*Candidate

	// 첫 줄(헤더) 건너뛰기
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	// 데이터 읽기
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // 파일 끝에 도달하면 종료
		}
		if err != nil {
			log.Printf("레코드 읽기 실패: %v", err)
			continue
		}

		candidate, err := g.createCandidate(record)
		if err != nil {
			log.Printf("후보자 생성 실패: %v", err)
			continue
		}

		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func (g *Game) createCandidate(record []string) (*Candidate, error) {
	celebId, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("ID 변환 실패: %v", err)
	}

	celebName := record[2]
	imageURL, err := g.searchImage(celebName)
	if err != nil {
		return nil, fmt.Errorf("이미지 검색 실패: %v", err)
	}

	log.Printf("후보자 생성: Name=%s, Image=%s", celebName, imageURL)

	return &Candidate{
		ID:        celebId,
		Name:      celebName,
		Image:     imageURL,
		FailStage: "",
		isUsed:    false,
	}, nil
}

func (g *Game) shuffleCandidates(candidates []*Candidate) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(candidates) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}
}

// 라운드 시작 시 모든 후보의 isUsed 초기화
func (r *Round) resetCandidatesStatus() {
	for i := range r.candidates {
		r.candidates[i].isUsed = false
	}
}
