package service

import (
	"encoding/json"
	"favorite_world_cup/config"
	"fmt"
	"net/http"
	"net/url"
)

// Kakao API 응답 구조체
type KakaoImageSearchResponse struct {
	Meta struct {
		TotalCount    int  `json:"total_count"`
		PageableCount int  `json:"pageable_count"`
		IsEnd         bool `json:"is_end"`
	} `json:"meta"`
	Documents []struct {
		Collection      string `json:"collection"`
		ThumbnailURL    string `json:"thumbnail_url"`
		ImageURL        string `json:"image_url"`
		Width           int    `json:"width"`
		Height          int    `json:"height"`
		DisplaySitename string `json:"display_sitename"`
		DocURL          string `json:"doc_url"`
		DateTime        string `json:"datetime"`
	} `json:"documents"`
}

func (g *Game) searchImage(query string) (string, error) {
	baseURL := "https://dapi.kakao.com/v2/search/image"
	params := url.Values{
		"query": {query},
		"sort":  {"accuracy"},
		"page":  {"1"},
		"size":  {"1"},
	}

	req, err := g.createRequest(baseURL, params)
	if err != nil {
		return "", err
	}

	imageURL, err := g.sendRequest(req)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

func (g *Game) createRequest(baseURL string, params url.Values) (*http.Request, error) {
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %v", err)
	}

	req.Header.Set("Authorization", config.KAKAO_API_KEY)
	return req, nil
}

func (g *Game) sendRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 필요한 필드만 파싱
	var result struct {
		Documents []struct {
			ImageURL string `json:"image_url"`
		} `json:"documents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("응답 파싱 실패: %v", err)
	}

	if len(result.Documents) == 0 {
		return "", fmt.Errorf("검색 결과 없음")
	}

	return result.Documents[0].ImageURL, nil
}
