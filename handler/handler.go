package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.Home)
	r.GET("/resources/:id", h.Get)
}

func (h *Handler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, Gin!",
	})
}

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "엄마를 위한 이상형 월드컵",
		"games": []gin.H{
			{
				"id":    1,
				"name":  "음식 월드컵",
				"image": "/static/images/food.jpg",
			},
			{
				"id":    2,
				"name":  "영화 월드컵",
				"image": "/static/images/movie.jpg",
			},
			{
				"id":    3,
				"name":  "여행지 월드컵",
				"image": "/static/images/travel.jpg",
			},
		},
	})
}
