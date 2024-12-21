package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
)

func HashPassword(password string) string {
	// Go inbuilt SHA256
	hasher := sha256.New()
	hasher.Write([]byte(password))             // convert to bytes before Writing it
	return hex.EncodeToString(hasher.Sum(nil)) // convert to string
}

// default appointment slots for professors, 10 am to 1pm, continued to 2pm to 4pm
// this will be 6 slots, T1 to T6.
func createProfessorSlots(db *gorm.DB, professorID uint8) error {
	var slots []Appointment

	for start := uint8(10); start <= 16; start++ {
		if start == 13 {
			// break slot
			continue
		}
		slots = append(slots, Appointment{
			UserID:       professorID,
			StartTime:    start,
			Availability: true,
		})
	}
	result := db.Create(&slots)
	return result.Error
}

func RegisterUser(db *gorm.DB, email, name, password string, userType uint8) error {
	// check for existing user with this email
	var existingUser User
	result := db.Where("email = ?", email).First(&existingUser)
	if result.RowsAffected > 0 {
		return errors.New("user already exists\n")
	}
	user := User{
		Email:    email,
		Name:     name,
		Password: HashPassword(password),
		Type:     userType,
	}

	result = db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	if userType == PROFESSOR_TYPE {
		// for proffs
		err := createProfessorSlots(db, user.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoginUser(db *gorm.DB, email, password string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	if HashPassword(password) != user.Password {
		// hashed password did not match with input
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
