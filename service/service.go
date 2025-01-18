package service

type Service struct {
	candidates []Candidate
	games      map[string]*Game
}

type Candidate struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Game struct {
	Username     string
	Candidates   []Candidate
	CurrentRound int
	Progress     float64
}

func NewService() *Service {
	return &Service{
		candidates: []Candidate{
			{ID: 1, Name: "후보1", Image: "/static/images/1.jpg"},
			{ID: 2, Name: "후보2", Image: "/static/images/2.jpg"},
			// ... 32개의 후보 추가
		},
		games: make(map[string]*Game),
	}
}

func (s *Service) InitGame(username string) map[string]interface{} {
	game := &Game{
		Username:     username,
		Candidates:   s.candidates,
		CurrentRound: 32,
		Progress:     0,
	}
	s.games[username] = game

	return map[string]interface{}{
		"username":     username,
		"currentRound": "32강",
		"progress":     0,
		"candidates": []Candidate{
			game.Candidates[0],
			game.Candidates[1],
		},
	}
}

func (s *Service) ProcessSelection(selectedID int) (map[string]interface{}, error) {
	// TODO: 선택된 후보 처리 및 다음 라운드 진행 로직 구현
	return map[string]interface{}{
		"finished": false,
	}, nil
}
