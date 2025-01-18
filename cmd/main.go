package main

import (
	"favorite_world_cup/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin 라우터 생성
	r := gin.Default()

	// HTML 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// 정적 파일 제공
	r.Static("/static", "./static")

	h := handler.NewHandler()
	h.RegisterRoutes(r)

	// 서버 실행 (기본 포트: 8080)
	r.Run()
}
