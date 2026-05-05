package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleGeneralManager Role = "general_manager"
	RoleDepartmentManager Role = "department_manager"
	RoleEmployee Role = "employee"
)

type Department struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null;unique" json:"name"`
	Phone       string     `json:"phone"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Username        string     `gorm:"not null;unique" json:"username"`
	Password        string     `gorm:"not null" json:"-"`
	Name            string     `gorm:"not null" json:"name"`
	Gender          string     `json:"gender"`
	BirthDate       *time.Time `json:"birth_date"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Address         string     `json:"address"`
	Role            Role       `gorm:"not null" json:"role"`
	DepartmentID    *uint      `json:"department_id"`
	Department      Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position        string     `json:"position"`
	JoinDate        *time.Time `json:"join_date"`
	Status          string     `gorm:"default:active" json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type AttendanceType string

const (
	AttendanceTypeClockIn AttendanceType = "clock_in"
	AttendanceTypeClockOut AttendanceType = "clock_out"
	AttendanceTypeLeave AttendanceType = "leave"
	AttendanceTypeBusinessTrip AttendanceType = "business_trip"
	AttendanceTypeOvertime AttendanceType = "overtime"
)

type AttendanceStatus string

const (
	AttendanceStatusPending AttendanceStatus = "pending"
	AttendanceStatusApproved AttendanceStatus = "approved"
	AttendanceStatusRejected AttendanceStatus = "rejected"
)

type Attendance struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	UserID      uint             `gorm:"not null" json:"user_id"`
	User        User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        AttendanceType   `gorm:"not null" json:"type"`
	Date        time.Time        `gorm:"not null" json:"date"`
	Time        string           `json:"time"`
	Status      AttendanceStatus `gorm:"default:pending" json:"status"`
	Reason      string           `json:"reason"`
	Remark      string           `json:"remark"`
	ApproverID  *uint            `json:"approver_id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type WorkTask struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Title          string     `gorm:"not null" json:"title"`
	Description    string     `json:"description"`
	AssigneeID     uint       `gorm:"not null" json:"assignee_id"`
	Assignee       User       `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	AssignorID     uint       `gorm:"not null" json:"assignor_id"`
	Assignor       User       `gorm:"foreignKey:AssignorID" json:"assignor,omitempty"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Priority       string     `gorm:"default:medium" json:"priority"`
	Status         string     `gorm:"default:pending" json:"status"`
	Progress       int        `gorm:"default:0" json:"progress"`
	Evaluation     string     `json:"evaluation"`
	EvaluationScore *int      `json:"evaluation_score"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Information struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Content     string     `json:"content"`
	Type        string     `gorm:"not null" json:"type"`
	AuthorID    uint       `gorm:"not null" json:"author_id"`
	Author      User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	IsPublic    bool       `gorm:"default:true" json:"is_public"`
	DepartmentID *uint      `json:"department_id"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Document struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	FilePath    string     `gorm:"not null" json:"file_path"`
	FileSize    int64      `json:"file_size"`
	FileType    string     `json:"file_type"`
	UploaderID  uint       `gorm:"not null" json:"uploader_id"`
	Uploader    User       `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
	IsPublic    bool       `gorm:"default:true" json:"is_public"`
	DepartmentID *uint      `json:"department_id"`
	DownloadCount int       `gorm:"default:0" json:"download_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Message struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SenderID    uint       `gorm:"not null" json:"sender_id"`
	Sender      User       `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	ReceiverID  uint       `gorm:"not null" json:"receiver_id"`
	Receiver    User       `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Subject     string     `gorm:"not null" json:"subject"`
	Content     string     `json:"content"`
	IsRead      bool       `gorm:"default:false" json:"is_read"`
	ReplyToID   *uint      `json:"reply_to_id"`
	ParentID    *uint      `json:"parent_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type WorkLog struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Date        time.Time  `gorm:"not null" json:"date"`
	WorkSummary string     `json:"work_summary"`
	WorkContent string     `json:"work_content"`
	TomorrowPlan string    `json:"tomorrow_plan"`
	Problems    string     `json:"problems"`
	TaskID      *uint      `json:"task_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type WorkReport struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title       string     `gorm:"not null" json:"title"`
	Content     string     `json:"content"`
	Type        string     `json:"type"`
	TaskID      *uint      `json:"task_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Department{},
		&User{},
		&Attendance{},
		&WorkTask{},
		&Information{},
		&Document{},
		&Message{},
		&WorkLog{},
		&WorkReport{},
	)
	if err != nil {
		return err
	}
	return nil
}

func InitDefaultData(db *gorm.DB) error {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	defaultDepartment := Department{
		Name:        "总经办",
		Phone:       "000-00000000",
		Description: "总经理办公室",
	}
	if err := db.Create(&defaultDepartment).Error; err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	defaultUser := User{
		Username:     "admin",
		Password:     string(hashedPassword),
		Name:         "系统管理员",
		Role:         RoleGeneralManager,
		DepartmentID: &defaultDepartment.ID,
		Status:       "active",
	}
	if err := db.Create(&defaultUser).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
