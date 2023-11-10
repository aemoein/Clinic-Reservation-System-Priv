package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/writer"
)

func main() {
	InitializeDB()
	defer CloseDB()

	CreateTables()

	r := mux.NewRouter()
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/signin", SignInHandler)
	http.HandleFunc("/AddSlot/{DId}/{APTime}/{STime}/{ETime}" , SetDoctorSchedulHandler).Methods("POST")
	http.HandleFunc("/CancelAppiontment/{APId}",  CancelAppiontmentHandler).Methods("DELETE")
	http.HandleFunc("/reviewReservations" , ViewPatientAppointmentsHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}


func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	usertype := r.FormValue("usertype")

	user := User{
		Name:     username,
		Email:    email,
		Password: password,
		UserType: usertype,
	}

	err := SignUp(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error signing up: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Signup successful")
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := SignIn(email, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error signing in: %v", err)
		return
	}

	func SetDoctorSchedulHandler(writer http.ResponseWriter, request *http.Request){
		newSlot := &Appointment{}
		var appointmentRequest AppointmentRequest
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&appointmentRequest)
		if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
		}
		r.ParseForm(request, SetDoctorSchedule())
		slot := newSlot.SetDoctorSchedul(appointmentRequest.DoctorID, appointmentRequest.AppointmentDate, appointmentRequest.start_time,appointmentRequest.end_time)
		res, _ := json.Marshal(slot)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(res)
	}

	func CancelAppointmentHandler(writer http.http.ResponseWriter, request *http.Request){
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

	func ViewPatientAppointmentsHandler(writer http.http.ResponseWriter, request *http.Request)  {
		reservations := mux.Vars(request)
		patientID := vars ["PId"]

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

	// Prepare the response
	response := map[string]interface{}{
		"message":   "Sign in successful",
		"user_type": user.UserType,
	}

	// Serialize the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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
