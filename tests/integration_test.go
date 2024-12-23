package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	model "github.com/bassamadnan/ucb-001xhttp/models"
	router "github.com/bassamadnan/ucb-001xhttp/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var duration int

func TestMain(m *testing.M) {
	// delay between each T.C
	duration = 2 // default 2 seconds
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &duration)
	}

	// run tests
	os.Exit(m.Run())
}

func TestAppointmentWorkflow(t *testing.T) {
	fmt.Printf("\nStarting test with %d second delay between steps...\n", duration)

	os.Remove("test.db") // entirely new db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Appointment{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	r := router.InitRouter(db)

	// store cookies for different users
	var A1Cookie, A2Cookie, P1Cookie string

	testCases := []struct {
		name     string
		method   string
		url      string
		body     interface{}
		cookie   *string // pointer to store cookie
		expected int
	}{
		{
			name:   "Register Student A1",
			method: "POST",
			url:    "/api/auth/register",
			body: map[string]interface{}{
				"email":    "A1@test.com",
				"password": "test123",
				"name":     "A1",
				"type":     0,
			},
			expected: 201,
		},
		{
			name:   "Register Student A2",
			method: "POST",
			url:    "/api/auth/register",
			body: map[string]interface{}{
				"email":    "A2@test.com",
				"password": "test123",
				"name":     "A2",
				"type":     0,
			},
			expected: 201,
		},
		{
			name:   "Register Professor P1",
			method: "POST",
			url:    "/api/auth/register",
			body: map[string]interface{}{
				"email":    "P1@test.com",
				"password": "test123",
				"name":     "P1",
				"type":     1,
			},
			expected: 201,
		},
		{
			name:   "Login Student A1",
			method: "POST",
			url:    "/api/auth/login",
			body: map[string]interface{}{
				"email":    "A1@test.com",
				"password": "test123",
			},
			cookie:   &A1Cookie,
			expected: 200,
		},
		{
			name:   "Login Professor P1",
			method: "POST",
			url:    "/api/auth/login",
			body: map[string]interface{}{
				"email":    "P1@test.com",
				"password": "test123",
			},
			cookie:   &P1Cookie,
			expected: 200,
		},
		{
			name:   "P1 Updates Slot T1 Availability",
			method: "PUT",
			url:    "/api/professor/slot",
			body: map[string]interface{}{
				"start_time": 10, // T1
				"available":  true,
			},
			cookie:   &P1Cookie,
			expected: 200,
		},
		{
			name:   "P1 Updates Slot T2 Availability",
			method: "PUT",
			url:    "/api/professor/slot",
			body: map[string]interface{}{
				"start_time": 11, // T2
				"available":  true,
			},
			cookie:   &P1Cookie,
			expected: 200,
		},
		{
			name:     "A1 Views P1's Available Slots",
			method:   "GET",
			url:      "/api/student/professor/P1@test.com/slots",
			cookie:   &A1Cookie,
			expected: 200,
		},
		{
			name:   "A1 Books Slot T1",
			method: "POST",
			url:    "/api/student/book",
			body: map[string]interface{}{
				"professor_email": "P1@test.com",
				"start_time":      10,
			},
			cookie:   &A1Cookie,
			expected: 200,
		},
		{
			name:   "Login Student A2",
			method: "POST",
			url:    "/api/auth/login",
			body: map[string]interface{}{
				"email":    "A2@test.com",
				"password": "test123",
			},
			cookie:   &A2Cookie,
			expected: 200,
		},
		{
			name:   "A2 Books Slot T2",
			method: "POST",
			url:    "/api/student/book",
			body: map[string]interface{}{
				"professor_email": "P1@test.com",
				"start_time":      11,
			},
			cookie:   &A2Cookie,
			expected: 200,
		},
		{
			name:   "P1 Cancels A1's Appointment",
			method: "POST",
			url:    "/api/professor/cancel",
			body: map[string]interface{}{
				"start_time": 10,
			},
			cookie:   &P1Cookie,
			expected: 200,
		},
		{
			name:     "A1 Checks Appointments",
			method:   "GET",
			url:      "/api/student/appointments",
			cookie:   &A1Cookie,
			expected: 200,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\nStep %d: %s\n", i+1, tc.name)

			var req *http.Request
			if tc.body != nil {
				jsonBody, _ := json.Marshal(tc.body)
				req, _ = http.NewRequest(tc.method, tc.url, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tc.method, tc.url, nil)
			}

			w := httptest.NewRecorder()

			// add session cookie if exists
			if tc.cookie != nil && *tc.cookie != "" {
				req.Header.Set("Cookie", *tc.cookie)
			}

			r.ServeHTTP(w, req)

			// store cookie from response if this test case expects to handle cookies
			if tc.cookie != nil {
				if cookies := w.Result().Cookies(); len(cookies) > 0 {
					*tc.cookie = cookies[0].String()
				}
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			fmt.Printf("Response: %+v\n", response)

			if w.Code != tc.expected {
				t.Errorf("Expected status %d, got %d", tc.expected, w.Code)
			}

			// wait till timeout
			time.Sleep(time.Duration(duration) * time.Second)
		})
	}

	os.Remove("test.db")
}
