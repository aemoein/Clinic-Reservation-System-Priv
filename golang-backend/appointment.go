package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Appointment struct {
	AppointmentID   int
	DoctorID        int
	PatientID       int
	AppointmentDate time.Time
	StartTime       time.Time
	EndTime         time.Time
	IsBooked        bool
}

// 3. Doctor sets his schedule
func SetDoctorSchedule(doctorID int, appointmentDate time.Time, startTime time.Time, endTime time.Time) error {
	// Check if the slot is already occupied
	isSlotOccupied, err := IsSlotOccupied(doctorID, appointmentDate, startTime, endTime)
	if err != nil {
		return err
	}

	if isSlotOccupied {
		return fmt.Errorf("the slot is already occupied")
	}

	// If the slot is not occupied, insert the new slot
	_, err = DB.Exec(`
		INSERT INTO appointments (doctor_id, appointment_date, start_time, end_time) 
		VALUES (?, ?, ?, ?)
	`, doctorID, appointmentDate, startTime, endTime)

	if err != nil {
		return err
	}

	return nil
}

// for this function I implemented another one depending on the new added field is_booked called IsSlotBooked
// IsSlotOccupied checks if the slot is already occupied
func IsSlotOccupied(doctorID int, appointmentDate time.Time, startTime time.Time, endTime time.Time) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM appointments 
		WHERE doctor_id = ? 
		AND appointment_date = ? 
		AND ((start_time <= ? AND end_time >= ?) OR (start_time <= ? AND end_time >= ?))
	`, doctorID, appointmentDate, startTime, startTime, endTime, endTime).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsSlotBooked(AppointmentID int) (bool, error) {
	// Query the database to check the is_booked status for the given AppointmentID
	var isBooked bool
	err := DB.QueryRow("SELECT is_booked FROM appointments WHERE appointment_id = ?", AppointmentID).Scan(&isBooked)
	if err != nil {
		return false, err
	}

	return isBooked, nil
}

/*func getAvailableSlotsFromDB(doctorID int) ([]Slot, error) {
		query := DB.QueryRow(`SELECT * FROM appointments`)
		return query

}*/

func getAvailableSlotsFromDB(db *sql.DB, doctorID int) ([]int, error) {
	// Define the query to select available slots for a specific doctor
	query := "SELECT * FROM appointments WHERE doctor_id = ? AND is_booked = TRUE"

	// Execute the query and retrieve the rows
	rows, err := db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []int

	// Iterate over the rows and scan the results into the slots slice
	for rows.Next() {
		var slotID int
		err := rows.Scan(&slotID)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slotID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
