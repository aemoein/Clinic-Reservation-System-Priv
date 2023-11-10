package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Appointment struct {
	AppointmentID   int 'jason: "AId" '
	DoctorID        int  'jason: "DId"'
	PatientID       int  'jason: "PId"'
	AppointmentDate time.Time 'jason: "APTime"'
	StartTime       time.Time 'jason: "STime"'
	EndTime         time.Time 'jason: "ETime"'
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

/*func getAvailableSlotsFromDB(doctorID int) ([]Slot, error) {
		query := DB.QueryRow(`SELECT * FROM appointments`)
		return query

}*/

func getAvailableSlotsFromDB(db *sql.DB, doctorID int) ([]int, error) {
	// Define the query to select available slots for a specific doctor
	query := "SELECT * FROM appointments WHERE doctor_id = ?"

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
