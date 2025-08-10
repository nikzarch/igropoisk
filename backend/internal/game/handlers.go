package game

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetGameByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	game, err := h.service.GetGameById(id)
	if err != nil {
		if err.Error() == "no such game" {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (h *Handler) GetAllGames(c *gin.Context) {
	games, err := h.service.GetAllGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (h *Handler) AddGame(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	game := Game{Name: req.Name}
	fmt.Println(req.Name)
	err := h.service.AddGame(game)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) DeleteGameByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	err = h.service.DeleteGameById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete game"})
		return
	}

	c.Status(http.StatusNoContent)
}
