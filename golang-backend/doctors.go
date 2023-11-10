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
