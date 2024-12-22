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

// allows students and professors to cancel appointments
// returns error if user does not have permission, or if appointment does not exist
func CancelAppointment(db *gorm.DB, userID uint8, startTime uint8, userType uint8, profEmail string) error {
	var slot Appointment
	var query *gorm.DB

	// if student is cancelling, need prof ID to locate the slot
	if userType == STUDENT_TYPE {
		var prof User
		result := db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
		if result.Error != nil {
			return errors.New("professor not found")
		}
		// locate slot under this prof, at this time, booked by this student
		query = db.Where("user_id = ? AND start_time = ? AND student_id = ?",
			prof.ID, startTime, userID)
	} else {
		// profs can cancel by their ID and time only
		query = db.Where("user_id = ? AND start_time = ?",
			userID, startTime)
	}

	result := query.First(&slot)
	if result.Error != nil {
		return errors.New("appointment not found")
	}

	// reset slot to available state
	slot.Availability = true
	slot.StudentID = 0

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
