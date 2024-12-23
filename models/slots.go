package models

import (
	"errors"

	"gorm.io/gorm"
)

// the logic for checking valid student/prof may change later if it can be handled via Gin

// all slots of a professor, should show which ones are available (or not) later on
func GetProfessorSlots(db *gorm.DB, profEmail string) ([]Appointment, error) {
	var prof User
	result := db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
	if result.Error != nil {
		return nil, errors.New("professor not found")
	}

	var slots []Appointment
	result = db.Where("user_id = ?", prof.ID).Find(&slots)
	if result.Error != nil {
		return nil, result.Error
	}
	return slots, nil
}

// stand alone for available slots only, may not need
func GetAvailableSlots(db *gorm.DB, profEmail string) ([]Appointment, error) {
	var prof User
	result := db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
	if result.Error != nil {
		return nil, errors.New("professor not found")
	}

	var slots []Appointment
	result = db.Where("user_id = ? AND availability = ?", prof.ID, true).Find(&slots)
	if result.Error != nil {
		return nil, result.Error
	}
	return slots, nil
}

// allows students and professors to cancel appointments
// returns error if user does not have permission, or if appointment does not exist
func CancelAppointment(db *gorm.DB, userEmail string, startTime uint8, userType uint8, profEmail string) error {
	var slot Appointment
	var query *gorm.DB

	// if student is cancelling, need prof ID to locate the slot
	if userType == STUDENT_TYPE {
		var prof User
		result := db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
		if result.Error != nil {
			return errors.New("professor not found")
		}

		var student User
		result = db.Where("email = ? AND type = ?", userEmail, STUDENT_TYPE).First(&student)
		if result.Error != nil {
			return errors.New("student not found")
		}

		// locate slot under this prof, at this time, booked by this student, and must be booked (not available)
		query = db.Where("user_id = ? AND start_time = ? AND student_id = ? AND availability = ?",
			prof.ID, startTime, student.ID, false)
	} else {
		var prof User
		result := db.Where("email = ? AND type = ?", userEmail, PROFESSOR_TYPE).First(&prof)
		if result.Error != nil {
			return errors.New("professor not found")
		}
		// profs can cancel by their ID and time only, but slot must be booked
		query = db.Where("user_id = ? AND start_time = ? AND availability = ?",
			prof.ID, startTime, false)
	}

	result := query.First(&slot)
	if result.Error != nil {
		return errors.New("no booked appointment found")
	}

	// reset slot to available state
	slot.Availability = true
	slot.StudentID = 0

	return db.Save(&slot).Error
}

// allow profs to update their availability (cancel appointments) may be a redundant endpoint
// based the assumpiton that prof can declare busy before the slots are even booked
func UpdateSlotAvailability(db *gorm.DB, profEmail string, startTime uint8, available bool) error {
	var prof User
	result := db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
	if result.Error != nil {
		return errors.New("professor not found")
	}

	var slot Appointment
	result = db.Where("user_id = ? AND start_time = ?", prof.ID, startTime).First(&slot)
	if result.Error != nil {
		return errors.New("slot not found")
	}

	slot.Availability = available
	return db.Save(&slot).Error
}

func BookAppointment(db *gorm.DB, studentEmail string, slotTime uint8, profEmail string) error {
	var student User
	result := db.Where("email = ? AND type = ?", studentEmail, STUDENT_TYPE).First(&student)
	if result.Error != nil {
		return errors.New("student not found")
	}
	// only done by students
	if student.Type != STUDENT_TYPE {
		return errors.New("user is not a student")
	}

	var prof User
	result = db.Where("email = ? AND type = ?", profEmail, PROFESSOR_TYPE).First(&prof)
	if result.Error != nil {
		return errors.New("professor not found")
	}

	// check for slot under this prof
	var slot Appointment
	result = db.Where("user_id = ? AND start_time = ? AND availability = ?", prof.ID, slotTime, true).First(&slot)
	if result.Error != nil {
		// the prof exists but is unavailable
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("slot not available")
		}
		return errors.New("professor not found")
	}
	slot.Availability = false
	slot.StudentID = student.ID
	return db.Save(&slot).Error
}

func GetStudentAppointments(db *gorm.DB, studentEmail string) ([]struct {
	StartTime uint8
	ProfName  string
}, error) {
	var student User
	result := db.Where("email = ? AND type = ?", studentEmail, STUDENT_TYPE).First(&student)
	if result.Error != nil {
		return nil, errors.New("student not found")
	}

	var appointments []struct {
		StartTime uint8
		ProfName  string
	}
	result = db.Table("appointments").
		Select("appointments.start_time, users.name as prof_name").
		Joins("JOIN users ON appointments.user_id = users.id").
		Where("appointments.student_id = ? AND appointments.availability = ?", student.ID, false).
		Scan(&appointments)

	if result.Error != nil {
		return nil, result.Error
	}
	return appointments, nil
}

// get all professors in the system
func GetAllProfessors(db *gorm.DB) ([]struct {
	Name  string
	Email string
}, error) {
	var professors []struct {
		Name  string
		Email string
	}
	result := db.Table("users").
		Select("name, email").
		Where("type = ?", PROFESSOR_TYPE).
		Scan(&professors)

	if result.Error != nil {
		return nil, result.Error
	}
	return professors, nil
}
