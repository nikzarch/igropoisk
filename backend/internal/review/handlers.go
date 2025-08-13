package review

import (
	"errors"
	"github.com/gin-gonic/gin"
	"igropoisk_backend/internal/user"
	"net/http"
	"strconv"
)

type AddReviewRequest struct {
	Content string    `json:"content"`
	GameId  int       `json:"-"`
	Score   int       `json:"Score"`
	User    user.User `json:"-"`
}

type Handler struct {
	reviewService Service
}

func NewHandler(service Service) *Handler {
	return &Handler{reviewService: service}
}

func validateRequest(req AddReviewRequest) error {
	if req.Score <= 0 || req.Score > 10 {
		return errors.New("Score must be between 0 and 10")
	}
	if req.GameId <= 0 {
		return errors.New("invalid game id")
	}
	return nil
}

func (h *Handler) AddReview(c *gin.Context) {
	req := new(AddReviewRequest)
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateRequest(*req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.User = user.User{Id: userId.(int), Name: username.(string)}
	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}
	req.GameId = gameID
	err = h.reviewService.AddReview(*req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetReviewsByGameId(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}
	var reviews []*Review
	reviews, err = h.reviewService.GetReviewsByGameId(gameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}
