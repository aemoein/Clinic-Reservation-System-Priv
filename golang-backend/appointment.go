package main

import (
	"database/sql"
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

	_, err = DB.Exec(`
		INSERT INTO appointments (doctor_id, appointment_date, start_time, end_time, patient_id) 
		VALUES (?, ?, ?, ?, NULL)
	`, doctorID, appointmentDate, startTime, endTime)

	if err != nil {
		return err
	}

	return nil
}

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

		log.Printf("Fetched Data: %+v\n", slot)
	}

	return slots, nil
}

func GetDoctorIDFromAppointment(appointmentID int) (int, error) {
	row := DB.QueryRow("SELECT doctor_id FROM appointments WHERE appointment_id = ?", appointmentID)

	var doctorID int

	if err := row.Scan(&doctorID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("appointment not found")
		}
		return 0, err
	}

	return doctorID, nil
}

func GetPatientIDFromAppointment(appointmentID int) (int, error) {
	row := DB.QueryRow("SELECT patient_id FROM appointments WHERE appointment_id = ?", appointmentID)

	var patientID int

	if err := row.Scan(&patientID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("appointment not found")
		}
		return 0, err
	}

	return patientID, nil
}
