package models

import (
	"errors"

	"gorm.io/gorm"
)

// the logic for checking valid student/prof may change later if it can be handled via Gin

// all slots of a professor, should show which ones are available (or not) later on
func GetProfessorSlots(db *gorm.DB, profID uint8) ([]Appointment, error) {
	var slots []Appointment
	result := db.Where("user_id = ?", profID).Find(&slots)
	if result.Error != nil {
		return nil, result.Error
	}
	return slots, nil
}

// stand alone for available slots only, may not need
func GetAvailableSlots(db *gorm.DB, profID uint8) ([]Appointment, error) {
	var slots []Appointment
	result := db.Where("user_id = ? AND availability = ?", profID, true).Find(&slots)
	if result.Error != nil {
		return nil, result.Error
	}
	return slots, nil
}

// change the availability of the slot, allows professors to mark slots as busy/free
func UpdateSlotAvailability(db *gorm.DB, profID uint8, slotID uint8, available bool) error {
	var slot Appointment
	result := db.Where("id = ? AND user_id = ?", slotID, profID).First(&slot)
	if result.Error != nil {
		return errors.New("slot not found or does not belong to professor")
	}

	slot.Availability = available
	if available == false {
		slot.StudentID = 0
	}
	return db.Save(&slot).Error
}

func BookAppointment(db *gorm.DB, studentID, slotTime uint8, profMail string) error {
	var student User
	result := db.First(&student, studentID)
	if result.Error != nil {
		return errors.New("student not found")
	}
	// only done by students
	if student.Type != STUDENT_TYPE {
		return errors.New("user is not a student")
	}

	var prof User
	result = db.Where("email = ? AND type = ?", profMail, PROFESSOR_TYPE).First(&prof)
	if result.Error != nil {
		return errors.New("professor not found")
	}

	// check for slot under this prof
	var slot Appointment
	result = db.Where("user_id = ? AND start_time = ? AND availability = ?", prof.ID, slotTime, true).First(&slot)
	if result.Error != nil {
		return errors.New("professor not found")
	}
	slot.Availability = false
	slot.StudentID = studentID
	return db.Save(&slot).Error
}

func GetStudentAppointments(db *gorm.DB, studentID uint8) ([]struct {
	StartTime uint8
	ProfName  string
}, error) {
	var appointments []struct {
		StartTime uint8
		ProfName  string
	}
	result := db.Table("appointments").
		Select("appointments.start_time, users.name as prof_name").
		Joins("JOIN users ON appointments.user_id = users.id").
		Where("appointments.student_id = ? AND appointments.availability = ?", studentID, false).
		Scan(&appointments)

	if result.Error != nil {
		return nil, result.Error
	}
	return appointments, nil
}
