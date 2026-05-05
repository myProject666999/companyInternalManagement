package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"companyInternalManagement/config"
	"companyInternalManagement/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var tasks []models.WorkTask
	query := config.DB.Preload("Assignee").Preload("Assignor")

	if user.DepartmentID != nil {
		query = query.Joins("JOIN users ON users.id = work_tasks.assignee_id").
			Where("users.department_id = ?", user.DepartmentID)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.WorkTask
	if err := config.DB.Preload("Assignee").Preload("Assignor").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")

	var task models.WorkTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	task.AssignorID = userID
	task.Progress = 0

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.WorkTask
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var updateData models.WorkTask
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := config.DB.Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.WorkTask
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func EvaluateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.WorkTask
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var evaluationData struct {
		Evaluation     string `json:"evaluation"`
		EvaluationScore *int  `json:"evaluation_score"`
	}

	if err := c.ShouldBindJSON(&evaluationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	task.Evaluation = evaluationData.Evaluation
	task.EvaluationScore = evaluationData.EvaluationScore
	task.Status = "completed"

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func ListDepartmentAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.DepartmentID == nil {
		c.JSON(http.StatusOK, []models.Attendance{})
		return
	}

	var attendances []models.Attendance
	query := config.DB.Preload("User").Joins("JOIN users ON users.id = attendances.user_id").
		Where("users.department_id = ?", user.DepartmentID)

	if attendanceType := c.Query("type"); attendanceType != "" {
		query = query.Where("attendances.type = ?", attendanceType)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("attendances.status = ?", status)
	}

	query = query.Order("attendances.created_at DESC")

	if err := query.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

func GetDepartmentAttendance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	var attendance models.Attendance
	if err := config.DB.Preload("User").First(&attendance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func ListDepartmentUsers(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.DepartmentID == nil {
		c.JSON(http.StatusOK, []models.User{})
		return
	}

	var users []models.User
	query := config.DB.Preload("Department").Where("department_id = ?", user.DepartmentID)

	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetDepartmentUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.DB.Preload("Department").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateDepartmentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var currentUser models.User
	if err := config.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Current user not found"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Name:         req.Name,
		Gender:       req.Gender,
		Phone:        req.Phone,
		Email:        req.Email,
		Address:      req.Address,
		Role:         models.Role(req.Role),
		DepartmentID: currentUser.DepartmentID,
		Position:     req.Position,
		Status:       "active",
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set password"})
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func UpdateDepartmentUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteDepartmentUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func ListDepartmentInformation(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var infoList []models.Information
	query := config.DB.Preload("Author")

	query = query.Where("is_public = ? OR department_id = ?", true, user.DepartmentID)

	if infoType := c.Query("type"); infoType != "" {
		query = query.Where("type = ?", infoType)
	}

	if title := c.Query("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&infoList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch information"})
		return
	}

	c.JSON(http.StatusOK, infoList)
}

func GetDepartmentInformation(c *gin.Context) {
	GetInformation(c)
}

func CreateDepartmentInformation(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var info models.Information
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	info.AuthorID = userID
	info.DepartmentID = user.DepartmentID
	info.ViewCount = 0

	if err := config.DB.Create(&info).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create information"})
		return
	}

	c.JSON(http.StatusCreated, info)
}

func UpdateDepartmentInformation(c *gin.Context) {
	UpdateInformation(c)
}

func DeleteDepartmentInformation(c *gin.Context) {
	DeleteInformation(c)
}

func ListDepartmentDocuments(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var documents []models.Document
	query := config.DB.Preload("Uploader")

	query = query.Where("is_public = ? OR department_id = ?", true, user.DepartmentID)

	if title := c.Query("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, documents)
}

func GetDepartmentDocument(c *gin.Context) {
	GetDocument(c)
}

func UploadDepartmentDocument(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	isPublic := c.PostForm("is_public") == "true"

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	uploadDir := "./uploads/documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	document := models.Document{
		Title:         title,
		Description:   description,
		FilePath:      filePath,
		FileSize:      header.Size,
		FileType:      ext,
		UploaderID:    userID,
		IsPublic:      isPublic,
		DepartmentID:  user.DepartmentID,
		DownloadCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document information"})
		return
	}

	c.JSON(http.StatusCreated, document)
}

func DownloadDepartmentDocument(c *gin.Context) {
	DownloadDocument(c)
}

func DeleteDepartmentDocument(c *gin.Context) {
	DeleteDocument(c)
}
