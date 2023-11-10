package main

import (
	"fmt"
	"log"
)

type Doctor struct {
	DoctorID   int    `json:"doctor_id"`
	DoctorName string `json:"doctor_name"`
}

func GetDoctors() ([]Doctor, error) {
	rows, err := DB.Query(`
		SELECT userid AS doctor_id, name AS doctor_name
		FROM users
		WHERE usertype = 'doctor'
	`)

	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	var doctors []Doctor

	for rows.Next() {
		var doctor Doctor
		err := rows.Scan(&doctor.DoctorID, &doctor.DoctorName)
		if err != nil {
			log.Printf("Error scanning rows: %v", err)
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		doctors = append(doctors, doctor)
	}

	log.Printf("Retrieved doctors: %+v", doctors)

	return doctors, nil
}

func FetchEmptyDoctorSlots(doctorID int) ([]Appointment, error) {
	rows, err := DB.Query(`
		SELECT appointment_id, IFNULL(patient_id, 0) as patient_id, 
		       appointment_date, start_time, end_time 
		FROM appointments 
		WHERE doctor_id = ? AND patient_id IS NULL
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
