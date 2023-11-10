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

func getAvailableSlotsFromDB(db *sql.DB, doctorID int) ([]Appointment, error) {
	// Define the query to select available slots for a specific doctor
	query := "SELECT * FROM appointments WHERE doctor_id = ?"

	// Execute the query and retrieve the rows
	rows, err := db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []Appointment

	// Iterate over the rows and scan the results into the slots slice
	for rows.Next() {
		var slot Appointment
		err := rows.Scan(&slot.AppointmentID, &slot.DoctorID, &slot.PatientID, &slot.AppointmentDate, &slot.StartTime, &slot.EndTime)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
