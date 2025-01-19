package handler

import (
	"favorite_world_cup/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	game *service.Game
}

func NewHandler() *Handler {
	game := service.NewGame()
	return &Handler{game: game}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.Home)
	r.GET("/name", h.ShowNameForm)
	r.POST("/start", h.StartGame)
	r.POST("/game/select", h.SelectCandidate)
}

func (h *Handler) Home(c *gin.Context) {
	data := map[string]interface{}{
		"title": "엄마를 위한 이상형 월드컵",
		"games": []map[string]interface{}{
			{
				"id":    1,
				"name":  "시작하기",
				"image": "/static/images/hoona.webp",
			},
		},
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func (h *Handler) ShowNameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "name.html", nil)
}

func (h *Handler) StartGame(c *gin.Context) {
	firstname := c.DefaultPostForm("firstname", "")
	var name string
	if firstname != "" {
		name = firstname + "씨"
	} else {
		name = "엄마"
	}
	data := h.game.InitGame(name)
	c.HTML(http.StatusOK, "game.html", data)
}

func (h *Handler) SelectCandidate(c *gin.Context) {
	var request struct {
		SelectedID int `json:"selectedId"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.game.ProcessSelection(request.SelectedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
