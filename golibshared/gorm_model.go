package golibshared

import "time"

// BaseModel definition of base gorm model
type BaseModel struct {
	ID       string `gorm:"primary_key" json:"id"`
	IsActive bool   `gorm:"column:is_active" json:"is_active,omitempty"`
}

// UserVerification definition of base gorm model with verification created user and modified user
type UserVerification struct {
	CreatedUser  string `gorm:"column:created_user;" json:"created_user,omitempty"`
	ModifiedUser string `gorm:"column:modified_user;" json:"modified_user,omitempty"`
}

// DateTimeStruct definition of date time base gorm model
type DateTimeStruct struct {
	CreatedAt  time.Time  `gorm:"column:created_at" json:"-"`
	ModifiedAt time.Time  `gorm:"column:modified_at" json:"-"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"-,omitempty"`
}

// JSONB define JsonBe for go
type JSONB map[string]interface{}
