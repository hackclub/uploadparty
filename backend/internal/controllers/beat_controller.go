package controllers

import (
	"net/http"
	"strconv"
	"uploadparty/internal/services"

	"github.com/gin-gonic/gin"
)

type BeatController struct {
	beatService *services.BeatService
}

func NewBeatController(beatService *services.BeatService) *BeatController {
	return &BeatController{
		beatService: beatService,
	}
}

func (bc *BeatController) Upload(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file required"})
		return
	}
	defer file.Close()

	var req services.UploadBeatRequest
	req.UserID = userID.(string)
	req.Title = c.PostForm("title")
	req.Description = c.PostForm("description")
	req.Genre = c.PostForm("genre")
	req.Key = c.PostForm("key")
	req.Tags = c.PostForm("tags")

	if bpm := c.PostForm("bpm"); bpm != "" {
		if parsedBPM, err := strconv.Atoi(bpm); err == nil {
			req.BPM = parsedBPM
		}
	}

	if isPublic := c.PostForm("is_public"); isPublic == "false" {
		req.IsPublic = false
	} else {
		req.IsPublic = true
	}

	beat, err := bc.beatService.Upload(&req, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, beat)
}

func (bc *BeatController) GetBeat(c *gin.Context) {
	beatID := c.Param("id")

	beat, err := bc.beatService.GetByID(beatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Beat not found"})
		return
	}

	c.JSON(http.StatusOK, beat)
}

func (bc *BeatController) GetBeats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	genre := c.Query("genre")
	sortBy := c.DefaultQuery("sort", "created_at")

	beats, total, err := bc.beatService.GetAll(page, limit, genre, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"beats": beats,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (bc *BeatController) GetUserBeats(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	beats, total, err := bc.beatService.GetByUserID(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"beats": beats,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (bc *BeatController) LikeBeat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	beatID := c.Param("id")

	err := bc.beatService.LikeBeat(userID.(string), beatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beat liked successfully"})
}

func (bc *BeatController) UnlikeBeat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	beatID := c.Param("id")

	err := bc.beatService.UnlikeBeat(userID.(string), beatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beat unliked successfully"})
}

func (bc *BeatController) AddComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req services.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID.(string)
	req.BeatID = c.Param("id")

	comment, err := bc.beatService.AddComment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (bc *BeatController) GetComments(c *gin.Context) {
	beatID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, total, err := bc.beatService.GetComments(beatID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
