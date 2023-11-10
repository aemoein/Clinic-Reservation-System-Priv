package main

import (
	"fmt"
	"log"
)

type Appointment struct {
	AppointmentID   int    `json:"appointment_id"`
	DoctorID        int    `json:"doctor_id"`
	PatientID       int    `json:"patient_id"`
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

// 3. Doctor sets his schedule
func SetDoctorSchedule(doctorID int, appointmentDate, startTime, endTime string) error {
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
		INSERT INTO appointments (doctor_id, appointment_date, start_time, end_time, patient_id) 
		VALUES (?, ?, ?, ?, NULL)
	`, doctorID, appointmentDate, startTime, endTime)

	if err != nil {
		return err
	}

	return nil
}

// IsSlotOccupied checks if the slot is already occupied
func IsSlotOccupied(doctorID int, appointmentDate, startTime, endTime string) (bool, error) {
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

func FetchDoctorSlots(doctorID int) ([]Appointment, error) {
	rows, err := DB.Query(`
		SELECT appointment_id, IFNULL(patient_id, 0) as patient_id, 
		       appointment_date, start_time, end_time 
		FROM appointments 
		WHERE doctor_id = ?
	`, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []Appointment

	for rows.Next() {
		var slot Appointment

		err := rows.Scan(&slot.AppointmentID, &slot.PatientID,
			&slot.AppointmentDate, &slot.StartTime, &slot.EndTime)
		if err != nil {
			return nil, err
		}

		slots = append(slots, slot)

		// Print the fetched data in the console
		log.Printf("Fetched Data: %+v\n", slot)
	}

	return slots, nil
}
