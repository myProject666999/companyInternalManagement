package routes

import (
	"companyInternalManagement/handlers"
	"companyInternalManagement/middlewares"
	"companyInternalManagement/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/change-password", middlewares.JWTAuth(), handlers.ChangePassword)
			auth.POST("/logout", middlewares.JWTAuth(), handlers.Logout)
			auth.GET("/profile", middlewares.JWTAuth(), handlers.GetProfile)
		}

		generalManager := api.Group("")
		generalManager.Use(middlewares.JWTAuth(), middlewares.RoleAuth(models.RoleGeneralManager))
		{
			departments := generalManager.Group("/departments")
			{
				departments.GET("", handlers.ListDepartments)
				departments.GET("/:id", handlers.GetDepartment)
				departments.POST("", handlers.CreateDepartment)
				departments.PUT("/:id", handlers.UpdateDepartment)
				departments.DELETE("/:id", handlers.DeleteDepartment)
			}

			users := generalManager.Group("/users")
			{
				users.GET("", handlers.ListUsers)
				users.GET("/:id", handlers.GetUser)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
				users.POST("/:id/promote", handlers.PromoteToManager)
			}

			info := generalManager.Group("/information")
			{
				info.GET("", handlers.ListInformation)
				info.GET("/:id", handlers.GetInformation)
				info.POST("", handlers.CreateInformation)
				info.PUT("/:id", handlers.UpdateInformation)
				info.DELETE("/:id", handlers.DeleteInformation)
			}

			documents := generalManager.Group("/documents")
			{
				documents.GET("", handlers.ListDocuments)
				documents.GET("/:id", handlers.GetDocument)
				documents.POST("", handlers.UploadDocument)
				documents.GET("/:id/download", handlers.DownloadDocument)
				documents.DELETE("/:id", handlers.DeleteDocument)
			}
		}

		deptManager := api.Group("")
		deptManager.Use(middlewares.JWTAuth(), middlewares.RoleAuth(models.RoleDepartmentManager))
		{
			tasks := deptManager.Group("/tasks")
			{
				tasks.GET("", handlers.ListTasks)
				tasks.GET("/:id", handlers.GetTask)
				tasks.POST("", handlers.CreateTask)
				tasks.PUT("/:id", handlers.UpdateTask)
				tasks.DELETE("/:id", handlers.DeleteTask)
				tasks.POST("/:id/evaluate", handlers.EvaluateTask)
			}

			deptAttendance := deptManager.Group("/department-attendance")
			{
				deptAttendance.GET("", handlers.ListDepartmentAttendance)
				deptAttendance.GET("/:id", handlers.GetDepartmentAttendance)
			}

			deptUsers := deptManager.Group("/department-users")
			{
				deptUsers.GET("", handlers.ListDepartmentUsers)
				deptUsers.GET("/:id", handlers.GetDepartmentUser)
				deptUsers.POST("", handlers.CreateDepartmentUser)
				deptUsers.PUT("/:id", handlers.UpdateDepartmentUser)
				deptUsers.DELETE("/:id", handlers.DeleteDepartmentUser)
			}

			deptInfo := deptManager.Group("/department-information")
			{
				deptInfo.GET("", handlers.ListDepartmentInformation)
				deptInfo.GET("/:id", handlers.GetDepartmentInformation)
				deptInfo.POST("", handlers.CreateDepartmentInformation)
				deptInfo.PUT("/:id", handlers.UpdateDepartmentInformation)
				deptInfo.DELETE("/:id", handlers.DeleteDepartmentInformation)
			}

			deptDocs := deptManager.Group("/department-documents")
			{
				deptDocs.GET("", handlers.ListDepartmentDocuments)
				deptDocs.GET("/:id", handlers.GetDepartmentDocument)
				deptDocs.POST("", handlers.UploadDepartmentDocument)
				deptDocs.GET("/:id/download", handlers.DownloadDepartmentDocument)
				deptDocs.DELETE("/:id", handlers.DeleteDepartmentDocument)
			}
		}

		employee := api.Group("")
		employee.Use(middlewares.JWTAuth(), middlewares.RoleAuth(models.RoleEmployee))
		{
			attendance := employee.Group("/my-attendance")
			{
				attendance.GET("", handlers.ListMyAttendance)
				attendance.GET("/:id", handlers.GetMyAttendance)
				attendance.POST("", handlers.CreateAttendance)
				attendance.PUT("/:id", handlers.UpdateAttendance)
				attendance.DELETE("/:id", handlers.DeleteAttendance)
			}

			employeeDocs := employee.Group("/employee-documents")
			{
				employeeDocs.GET("", handlers.ListEmployeeDocuments)
				employeeDocs.GET("/:id", handlers.GetEmployeeDocument)
				employeeDocs.POST("", handlers.UploadEmployeeDocument)
				employeeDocs.GET("/:id/download", handlers.DownloadEmployeeDocument)
				employeeDocs.DELETE("/:id", handlers.DeleteEmployeeDocument)
			}

			companyInfo := employee.Group("/company-information")
			{
				companyInfo.GET("", handlers.ListCompanyInformation)
				companyInfo.GET("/:id", handlers.GetCompanyInformation)
			}

			workLogs := employee.Group("/work-logs")
			{
				workLogs.GET("", handlers.ListWorkLogs)
				workLogs.GET("/:id", handlers.GetWorkLog)
				workLogs.POST("", handlers.CreateWorkLog)
				workLogs.PUT("/:id", handlers.UpdateWorkLog)
				workLogs.DELETE("/:id", handlers.DeleteWorkLog)
			}

			workReports := employee.Group("/work-reports")
			{
				workReports.GET("", handlers.ListWorkReports)
				workReports.GET("/:id", handlers.GetWorkReport)
				workReports.POST("", handlers.CreateWorkReport)
			}
		}

		messages := api.Group("/messages")
		messages.Use(middlewares.JWTAuth())
		{
			messages.GET("/inbox", handlers.ListInbox)
			messages.GET("/outbox", handlers.ListOutbox)
			messages.GET("/:id", handlers.GetMessage)
			messages.POST("", handlers.SendMessage)
			messages.POST("/:id/reply", handlers.ReplyMessage)
			messages.DELETE("/:id", handlers.DeleteMessage)
		}

		shared := api.Group("")
		shared.Use(middlewares.JWTAuth())
		{
			shared.GET("/all-users", handlers.ListAllUsers)
			shared.GET("/all-departments", handlers.ListAllDepartments)
		}
	}
}
