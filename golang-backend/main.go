package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	InitializeDB()
	defer CloseDB()
	CreateTables()

	router := mux.NewRouter()

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to My Go Backend!")
	}).Methods("GET")

	router.HandleFunc("/signup", SignUpHandler).Methods("POST")
	router.HandleFunc("/signin", SignInHandler).Methods("POST")
	router.HandleFunc("/slots/view", ViewDoctorSlotsHandler).Methods("GET").Queries("doctorid", "{doctorid}")
	router.HandleFunc("/slots/view/empty", ViewEmptyDoctorSlotsHandler).Methods("GET").Queries("doctorid", "{doctorid}")
	router.HandleFunc("/slots/add", SetDoctorScheduleHandler).Methods("POST")
	router.HandleFunc("/appointments/reserve", ReserveAppointmentHandler).Methods("POST")
	router.HandleFunc("/appointments/update", UpdateAppointmentHandler).Methods("POST")
	router.HandleFunc("/appointments/cancel", CancelAppointmentHandler).Methods("PUT").Queries("appointmentid", "{appointmentid}")
	router.HandleFunc("/appointments/view", ViewPatientAppointmentsHandler).Methods("GET").Queries("patientid", "{patientid}")
	router.HandleFunc("/doctors", GetDoctorsHandler).Methods("GET")
	router.HandleFunc("/ws", handleWebSocket)

	http.ListenAndServe(":8081",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router))
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	log.Printf("Received data: %+v", user)

	err = SignUp(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error signing up: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Signup successful")
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	user, err := SignIn(credentials.Email, credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error signing in: %v", err)
		return
	}

	log.Printf("User data: %+v", user)

	response := map[string]interface{}{
		"message":   "Sign in successful",
		"userid":    user.UserID,
		"user_type": user.UserType,
		"username":  user.UserName,
	}

	log.Printf("User data id: %+v", user.UserID)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func ViewDoctorSlotsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorIDStr, ok := vars["doctorid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Doctor ID not provided in the URL")
		return
	}

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid doctor ID: %v", err)
		return
	}

	log.Printf("ID received: %d", doctorID)

	slots, err := FetchDoctorSlots(doctorID)
	if err != nil {
		log.Printf("Error fetching doctor slots: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching doctor slots: %v", err)
		return
	}

	jsonData, err := json.Marshal(slots)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating JSON response: %v", err)
		return
	}
	log.Printf("JSON being sent: %s", jsonData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func ViewEmptyDoctorSlotsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorIDStr, ok := vars["doctorid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Doctor ID not provided in the URL")
		return
	}

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid doctor ID: %v", err)
		return
	}

	log.Printf("ID received: %d", doctorID)

	slots, err := FetchEmptyDoctorSlots(doctorID)
	if err != nil {
		log.Printf("Error fetching doctor slots: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching doctor slots: %v", err)
		return
	}

	jsonData, err := json.Marshal(slots)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating JSON response: %v", err)
		return
	}
	log.Printf("JSON being sent: %s", jsonData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func SetDoctorScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var newSlot Appointment
	err := json.NewDecoder(r.Body).Decode(&newSlot)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: %+v", newSlot)

	if newSlot.DoctorID <= 0 || newSlot.AppointmentDate == "" || newSlot.StartTime == "" || newSlot.EndTime == "" {
		http.Error(w, "Invalid slot data", http.StatusBadRequest)
		return
	}

	err = SetDoctorSchedule(newSlot.DoctorID, newSlot.AppointmentDate, newSlot.StartTime, newSlot.EndTime)
	if err != nil {
		log.Printf("Error adding slot: %v", err)
		http.Error(w, "Failed to add slot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Slot added successfully")
}

func ReserveAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var reservationRequest struct {
		AppointmentID int `json:"appointment_id"`
		PatientID     int `json:"patient_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reservationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the function to reserve the appointment
	if err := ReserveAppointment(reservationRequest.AppointmentID, reservationRequest.PatientID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Appointment reserved successfully")
}

func UpdateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var updateRequest struct {
		AppointmentID    int `json:"appointment_id"`
		OldAppointmentID int `json:"old_appointment_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding JSON request: %v", err)
		return
	}
	log.Printf("JSON request decoded successfully: %v", updateRequest)

	// Call the function to reserve the appointment
	if err := UpdateAppointment(updateRequest.AppointmentID, updateRequest.OldAppointmentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error updating appointment: %v", err)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Appointment reserved successfully")
}

func CancelAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	AppointmentIdStr, ok := vars["appointmentid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "appointment id not provided in the URL")
		return
	}

	AppointmentId, err := strconv.Atoi(AppointmentIdStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	log.Printf("ID received: %d", AppointmentId)

	appointments := CancelAppointment(AppointmentId)
	if err != nil {
		log.Printf("Error fetching patient appointments: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching patient appointments: %v", err)
		return
	}

	jsonData, err := json.Marshal(appointments)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating JSON response: %v", err)
		return
	}
	log.Printf("JSON being sent: %s", jsonData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func ViewPatientAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientIDStr, ok := vars["patientid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "patient id not provided in the URL")
		return
	}

	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, "Invalid Patient ID", http.StatusBadRequest)
		return
	}

	log.Printf("patient ID received: %d", patientID)

	appointments, err := ViewPatientAppointments(patientID)
	if err != nil {
		log.Printf("Error fetching patient appointments: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching patient appointments: %v", err)
		return
	}

	jsonData, err := json.Marshal(appointments)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating JSON response: %v", err)
		return
	}
	log.Printf("JSON being sent: %s", jsonData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	doctors, err := GetDoctors()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching doctors: %v", err)
		return
	}

	// Serialize the doctors to JSON
	jsonResponse, err := json.Marshal(doctors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating JSON response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	messages := consumeOldKafkaMessages()

	for _, message := range messages {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			//fmt.Println("Error broadcasting:", err)
			return
		}
	}
	message := consumeKafkaMessages()
	var oldMessage = message

	for {
		message := consumeKafkaMessages()
		if !reflect.DeepEqual(message, oldMessage) {
			err := conn.WriteMessage(websocket.TextMessage, message.Value)
			if err != nil {
				fmt.Println("Error broadcasting:", err)
				return
			}
			oldMessage = message
		}
	}
}

/*
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:7R26@llg4grb$&@tcp(127.0.0.1:3306)/Clinical_Reservation")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", "Ahmed Elsayed", "ahmed33elsayed22@gmail.com", "12345", "patient")
	if err != nil {
		panic(err.Error())
	}

	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, "Ahmed Elsayed", "2002-06-13", "male")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("User added to both users and patients tables successfully!")
}*/
