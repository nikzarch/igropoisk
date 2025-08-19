package game

import (
	"errors"
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

	game, err := h.service.GetGameByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (h *Handler) GetAllGames(c *gin.Context) {
	games, err := h.service.GetAllGames(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get games"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"games": games})
}

type AddGameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Genre       string `json:"genre"`
}

func (h *Handler) AddGame(c *gin.Context) {
	var req AddGameRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateAddGameRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddGame(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func validateAddGameRequest(request AddGameRequest) (err error) {
	switch {
	case len(request.Name) == 0:
		err = errors.New("name is required")
	case len(request.Description) == 0:
		err = errors.New("description is required")
	case len(request.ImageURL) == 0:
		err = errors.New("image_url is required")
	case len(request.Genre) == 0:
		err = errors.New("genre is required")
	}
	return err
}

func (h *Handler) DeleteGameByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	err = h.service.DeleteGameByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete game"})
		return
	}

	c.Status(http.StatusNoContent)
}
