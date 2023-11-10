package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	//http.HandleFunc("/view",)
	router.HandleFunc("/slots/view", DoctorSlotsHandler).Methods("GET").Queries("doctorid", "{doctorid}")

	http.ListenAndServe(":8081",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router))
	r := mux.NewRouter()

	http.HandleFunc("/AddSlot/{DId}/{APTime}/{STime}/{ETime}", SetDoctorSchedulHandler).Methods("POST")
	http.HandleFunc("/CancelAppiontment/{APId}", CancelAppiontmentHandler).Methods("DELETE")
	http.HandleFunc("/reviewReservations", ViewPatientAppointmentsHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
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

	//log.Printf("Received data: %+v", credentials)

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

func DoctorSlotsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorIDStr, ok := vars["doctorid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Doctor ID not provided in the URL")
		return
	}
  
	// Convert doctorIDStr to an integer
	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid doctor ID: %v", err)
		return
	}

	log.Printf("ID received: %d", doctorID)

	slots, err := FetchDoctorSlots(doctorID)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error fetching doctor slots: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching doctor slots: %v", err)
		return
	}

	jsonData, err := json.Marshal(slots)
	if err != nil {
		// Log the error for debugging
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

func SetDoctorScheduleHandler(writer http.ResponseWriter, request *http.Request) {
	newSlot := &Appointment{}
	var appointmentRequest AppointmentRequest
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&appointmentRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	r.ParseForm(request, SetDoctorSchedule())
	slot := newSlot.SetDoctorSchedul(appointmentRequest.DoctorID, appointmentRequest.AppointmentDate, appointmentRequest.start_time, appointmentRequest.end_time)
	res, _ := json.Marshal(slot)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(res)
}

func CancelAppointmentHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	AppointmentId := vars["APId"]

	id, err := strconv.Atoi(AppointmentId)
	if err != nil {
		http.Error(writer, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var appointmentRequest AppointmentRequest
	err = json.NewDecoder(request.Body).Decode(&appointmentRequest)
	if err != nil {
		http.Error(writer, "Error decoding request body", http.StatusBadRequest)
		return
	}

	res, _ := json.Marshal(CancelAppointment(int(id)))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
	writer.Write(res)
}

func ViewPatientAppointmentsHandler(writer http.ResponseWriter, request *http.Request) {
	reservations := mux.Vars(request)
	patientID := vars["PId"]

	id, err := strconv.Atoi(patientID)
	if err != nil {
		http.Error(writer, "Invalid Patient ID", http.StatusBadRequest)
		return
	}

	var appointmentRequest AppointmentRequest
	err = json.NewDecoder(request.Body).Decode(&appointmentRequest)
	if err != nil {
		http.Error(writer, "Error decoding request body", http.StatusBadRequest)
		return
	}

	res, _ := json.Marshal(ViewPatientAppointments(patientID))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
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
