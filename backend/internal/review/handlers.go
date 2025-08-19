package review

import (
	"errors"
	"github.com/gin-gonic/gin"
	"igropoisk_backend/internal/middleware"
	"igropoisk_backend/internal/user"
	"net/http"
	"strconv"
)

type AddReviewRequest struct {
	Content string    `json:"content"`
	GameID  int       `json:"-"`
	Rating  int       `json:"rating"`
	User    user.User `json:"-"`
}

type Handler struct {
	reviewService Service
}

func NewHandler(service Service) *Handler {
	return &Handler{reviewService: service}
}

func validateRequest(req AddReviewRequest) error {
	if req.Rating <= 0 || req.Rating > 10 {
		return errors.New("Rating must be between 0 and 10")
	}
	return nil
}

func (h *Handler) AddReview(c *gin.Context) {
	var req AddReviewRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := validateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, ok := c.Request.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not found"})
	}
	username, _ := c.Request.Context().Value(middleware.UserNameKey).(string)
	req.User = user.User{ID: userID, Name: username}
	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}
	req.GameID = gameID
	err = h.reviewService.AddReview(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetReviewsByGameID(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}
	reviews, err := h.reviewService.GetReviewsByGameID(c.Request.Context(), gameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
