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

func ListMyAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	var attendances []models.Attendance
	query := config.DB.Where("user_id = ?", userID)

	if attendanceType := c.Query("type"); attendanceType != "" {
		query = query.Where("type = ?", attendanceType)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

func GetMyAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
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

	if attendance.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func CreateAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	attendance.UserID = userID
	attendance.Status = models.AttendanceStatusPending

	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance record"})
		return
	}

	c.JSON(http.StatusCreated, attendance)
}

func UpdateAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	var attendance models.Attendance
	if err := config.DB.First(&attendance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	if attendance.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateData models.Attendance
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := config.DB.Model(&attendance).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance record"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func DeleteAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	var attendance models.Attendance
	if err := config.DB.First(&attendance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	if attendance.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := config.DB.Delete(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance record deleted successfully"})
}

func ListEmployeeDocuments(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var documents []models.Document
	query := config.DB.Preload("Uploader")

	query = query.Where("is_public = ? OR department_id = ? OR uploader_id = ?", true, user.DepartmentID, userID)

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

func GetEmployeeDocument(c *gin.Context) {
	GetDocument(c)
}

func UploadEmployeeDocument(c *gin.Context) {
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

func DownloadEmployeeDocument(c *gin.Context) {
	DownloadDocument(c)
}

func DeleteEmployeeDocument(c *gin.Context) {
	DeleteDocument(c)
}

func ListCompanyInformation(c *gin.Context) {
	var infoList []models.Information
	query := config.DB.Preload("Author").Where("is_public = ?", true)

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

func GetCompanyInformation(c *gin.Context) {
	GetInformation(c)
}

func ListWorkLogs(c *gin.Context) {
	userID := c.GetUint("user_id")

	var workLogs []models.WorkLog
	query := config.DB.Where("user_id = ?", userID)

	if workSummary := c.Query("work_summary"); workSummary != "" {
		query = query.Where("work_summary LIKE ?", "%"+workSummary+"%")
	}

	query = query.Order("date DESC")

	if err := query.Find(&workLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work logs"})
		return
	}

	c.JSON(http.StatusOK, workLogs)
}

func GetWorkLog(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work log ID"})
		return
	}

	var workLog models.WorkLog
	if err := config.DB.Preload("User").First(&workLog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work log not found"})
		return
	}

	if workLog.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, workLog)
}

func CreateWorkLog(c *gin.Context) {
	userID := c.GetUint("user_id")

	var workLog models.WorkLog
	if err := c.ShouldBindJSON(&workLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	workLog.UserID = userID

	if err := config.DB.Create(&workLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work log"})
		return
	}

	c.JSON(http.StatusCreated, workLog)
}

func UpdateWorkLog(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work log ID"})
		return
	}

	var workLog models.WorkLog
	if err := config.DB.First(&workLog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work log not found"})
		return
	}

	if workLog.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateData models.WorkLog
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := config.DB.Model(&workLog).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work log"})
		return
	}

	c.JSON(http.StatusOK, workLog)
}

func DeleteWorkLog(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work log ID"})
		return
	}

	var workLog models.WorkLog
	if err := config.DB.First(&workLog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work log not found"})
		return
	}

	if workLog.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := config.DB.Delete(&workLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete work log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work log deleted successfully"})
}

func ListWorkReports(c *gin.Context) {
	userID := c.GetUint("user_id")

	var workReports []models.WorkReport
	query := config.DB.Where("user_id = ?", userID).Order("created_at DESC")

	if err := query.Find(&workReports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work reports"})
		return
	}

	c.JSON(http.StatusOK, workReports)
}

func GetWorkReport(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work report ID"})
		return
	}

	var workReport models.WorkReport
	if err := config.DB.Preload("User").First(&workReport, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work report not found"})
		return
	}

	if workReport.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, workReport)
}

func CreateWorkReport(c *gin.Context) {
	userID := c.GetUint("user_id")

	var workReport models.WorkReport
	if err := c.ShouldBindJSON(&workReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	workReport.UserID = userID

	if err := config.DB.Create(&workReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work report"})
		return
	}

	c.JSON(http.StatusCreated, workReport)
}
