package handlers

import (
	"net/http"
	"strconv"

	"companyInternalManagement/config"
	"companyInternalManagement/models"

	"github.com/gin-gonic/gin"
)

func ListInbox(c *gin.Context) {
	userID := c.GetUint("user_id")

	var messages []models.Message
	query := config.DB.Preload("Sender").Preload("Receiver").Where("receiver_id = ?", userID)

	if search := c.Query("search"); search != "" {
		query = query.Where("subject LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inbox messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func ListOutbox(c *gin.Context) {
	userID := c.GetUint("user_id")

	var messages []models.Message
	query := config.DB.Preload("Sender").Preload("Receiver").Where("sender_id = ?", userID)

	query = query.Order("created_at DESC")

	if err := query.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch outbox messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := config.DB.Preload("Sender").Preload("Receiver").First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	if message.SenderID != userID && message.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if message.ReceiverID == userID && !message.IsRead {
		message.IsRead = true
		config.DB.Save(&message)
	}

	c.JSON(http.StatusOK, message)
}

func SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if message.ReceiverID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send message to yourself"})
		return
	}

	message.SenderID = userID
	message.IsRead = false

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func ReplyMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var parentMessage models.Message
	if err := config.DB.First(&parentMessage, parentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent message not found"})
		return
	}

	if parentMessage.SenderID != userID && parentMessage.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var replyMessage models.Message
	if err := c.ShouldBindJSON(&replyMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	replyMessage.SenderID = userID
	if parentMessage.SenderID == userID {
		replyMessage.ReceiverID = parentMessage.ReceiverID
	} else {
		replyMessage.ReceiverID = parentMessage.SenderID
	}
	replyMessage.ReplyToID = &parentMessage.ID
	replyMessage.ParentID = &parentMessage.ID
	replyMessage.IsRead = false

	if err := config.DB.Create(&replyMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reply"})
		return
	}

	c.JSON(http.StatusCreated, replyMessage)
}

func DeleteMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := config.DB.First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	if message.SenderID != userID && message.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := config.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
