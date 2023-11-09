package main

import (
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

// SetDoctorSchedule sets the schedule for a doctor
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
