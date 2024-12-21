package models

// this file contains the schemas for Users and Appointments

const (
	STUDENT_TYPE uint8 = iota
	PROFESSOR_TYPE
)

// Type is student or prof, not kept 0 or 1 as we may have more users later
type User struct {
	ID           uint8         `gorm:"primaryKey"`
	Type         uint8         `gorm:"not null"`
	Email        string        `gorm:"uniqueIndex;size:255;not null"` // should be unique
	Password     string        `gorm:"not null"`                      // no size limit as this will be encrypted
	Name         string        `gorm:"size:255;not null"`
	Appointments []Appointment `gorm:"foreignKey:UserID"`
}

// assume each appointment would be an hour long, we will not be storing the "EndTime" as such.
type Appointment struct {
	ID           uint8 `gorm:"primaryKey"`
	UserID       uint8 `gorm:"not null"`              // foreign, for the prof
	StudentID    uint8 `gorm:"default:null"`          // student assigned this slot
	StartTime    uint8 `gorm:"not null"`              // 24-hour format (0-23)
	Availability bool  `gorm:"not null;default:true"` // 1 for free slot, 0 for occupied/busy (for profs)
}
