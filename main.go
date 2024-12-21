package main

import (
	"fmt"
	"log"

	model "github.com/bassamadnan/ucb-001xhttp/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("models/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// drop existing tables if any and recreate
	db.Migrator().DropTable(&model.User{}, &model.Appointment{})

	// migrate the schemas
	err = db.AutoMigrate(&model.User{}, &model.Appointment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// verify tables were created
	if db.Migrator().HasTable(&model.User{}) {
		fmt.Println("Users table created successfully")
	}
	if db.Migrator().HasTable(&model.Appointment{}) {
		fmt.Println("Appointments table created successfully")
	}

	// register a professor
	err = model.RegisterUser(db, "prof@example.com", "Professor Smith", "prof123", model.PROFESSOR_TYPE)
	if err != nil {
		log.Printf("Error registering professor: %v\n", err)
	}
	// register a student
	err = model.RegisterUser(db, "student@example.com", "Student Jones", "student123", model.STUDENT_TYPE)
	if err != nil {
		log.Printf("Error registering student: %v\n", err)
	}

	// test login for both users
	prof, err := model.LoginUser(db, "prof@example.com", "prof123")
	if err != nil {
		log.Printf("Professor login error: %v\n", err)
	} else {
		fmt.Printf("Professor logged in successfully: %s\n", prof.Name)
	}

	student, err := model.LoginUser(db, "student@example.com", "student123")
	if err != nil {
		log.Printf("Student login error: %v\n", err)
	} else {
		fmt.Printf("Student logged in successfully: %s\n", student.Name)
	}

	// get slots for prof
	slots, err := model.GetProfessorSlots(db, prof.ID)
	if err != nil {
		log.Printf("Error getting slots: %v\n", err)
	} else {
		fmt.Printf("Professor has %d slots\n", len(slots))
	}

	// try booking an appointment
	err = model.BookAppointment(db, student.ID, 10, "prof@example.com") // booking 10 AM slot
	if err != nil {
		log.Printf("Error booking appointment: %v\n", err)
	} else {
		fmt.Println("Appointment booked successfully")
	}

	// get student appointments
	appointments, err := model.GetStudentAppointments(db, student.ID)
	if err != nil {
		log.Printf("Error getting student appointments: %v\n", err)
	} else {
		for _, apt := range appointments {
			fmt.Printf("Appointment with %s at %d:00\n", apt.ProfName, apt.StartTime)
		}
	}

	// update slot availability (prof marking 2 PM slot as unavailable)
	err = model.UpdateSlotAvailability(db, prof.ID, slots[3].ID, false)
	if err != nil {
		log.Printf("Error updating slot: %v\n", err)
	}

	// print tables at the end
	fmt.Println("\nUsers table:")
	var users []model.User
	db.Find(&users)
	for _, u := range users {
		fmt.Printf("ID: %d, Type: %d, Name: %s, Email: %s\n", u.ID, u.Type, u.Name, u.Email)
	}

	fmt.Println("\nAppointments table:")
	var allSlots []model.Appointment
	db.Find(&allSlots)
	for _, s := range allSlots {
		fmt.Printf("ID: %d, UserID: %d, StudentID: %d, StartTime: %d, Available: %v\n",
			s.ID, s.UserID, s.StudentID, s.StartTime, s.Availability)
	}
}
