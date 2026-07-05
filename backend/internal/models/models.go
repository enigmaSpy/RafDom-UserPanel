package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Surname      string    `gorm:"type:varchar(100)"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Phone        string    `gorm:"type:varchar(100)"`
	PasswordHash string    `gorm:"type:varchar(150);not null" json:"-"`
	Role         string    `gorm:"type:varchar(20);not null;default:'client'"`
	Address      string    `gorm:"type:varchar(100)"`
	City         string    `gorm:"type:varchar(100)"`
	PostalCode   string    `gorm:"type:varchar(7)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Messages []Message `gorm:"foreignKey:SenderID"`
}

type Renovation struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClientID    uuid.UUID `gorm:"type:uuid;not null"`
	Client      User      `gorm:"foreignKey:ClientID"`
	Name        string    `gorm:"type:varchar(50);not null"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"type:varchar(20);default:'estimation'"` // estimation, in_progress, completed
	CreatedAt   time.Time
	UpdatedAt   time.Time

	LaborTasks      []LaborTask      `gorm:"foreignKey:RenovationID"`
	Transactions    []Transaction    `gorm:"foreignKey:RenovationID"`
	ProgressUpdates []ProgressUpdate `gorm:"foreignKey:RenovationID"`
	Messages        []Message        `gorm:"foreignKey:RenovationID"`
}

type LaborTask struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RenovationID uuid.UUID `gorm:"type:uuid;not null"`
	Label        string    `gorm:"type:varchar(100);not null"`
	Status       string    `gorm:"type:varchar(20);default:'pending'"` //pending, in_progress, completed

	UnitPrice float64 `gorm:"not null"`
	Unit      string  `gorm:"type:varchar(20);not nul"`
	Quantity  float64 `gorm:"not null"`

	Amount     float64 `gorm:"not null"`
	Date       time.Time
	Note       string     `gorm:"type:text"`
	Renovation Renovation `gorm:"foreignKey:RenovationID" json:"-"`
}

type Transaction struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RenovationID uuid.UUID  `gorm:"type:uuid;not null"`
	Type         string     `gorm:"type:varchar(30);not null"` //material_deposit, material_expense, labor_payment
	Amount       float64    `gorm:"not null"`
	Date         time.Time  `gorm:"not null"`
	Note         string     `gorm:"type:text"`
	DocumentURL  *string    `gorm:"type:varchar(255)"`
	Renovation   Renovation `gorm:"foreignKey:RenovationID" json:"-"`
}

type ProgressUpdate struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RenovationID uuid.UUID  `gorm:"type:uuid;not null" json:"renovation_id"`
	LaborTaskID  *uuid.UUID `gorm:"type:uuid" json:"labor_task_id,omitempty"`
	
	Title       string    `gorm:"type:varchar(150);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	
	Photos      []string  `gorm:"type:jsonb;serializer:json" json:"photos"` 
	Date        time.Time `json:"date"`
	Renovation  Renovation `gorm:"foreignKey:RenovationID" json:"-"`
	LaborTask   *LaborTask `gorm:"foreignKey:LaborTaskID" json:"-"`
}

type Message struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RenovationID uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderID     uuid.UUID `gorm:"type:uuid;not null"`
	ReceiverID   uuid.UUID `gorm:"type:uuid;not null"`

	Content string `gorm:"type:text;not null"`
	IsRead  bool   `gorm:"default:false"`

	CreatedAt  time.Time
	Renovation Renovation `gorm:"foreignKey:RenovationID" json:"-"`
	Sender     User       `gorm:"foreignKey:SenderID" json:"-"`
	Receiver   User       `gorm:"foreignKey:ReceiverID" json:"-"`
}
