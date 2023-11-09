package main

import (
	"time"
)

// 4- Patients select doctor, view his available slots, then patient chooses a slot.
func ViewDoctorSlots(doctorID int, appointmentDate time.Time) ([]Appointment, error) {
	rows, err := DB.Query(`
		SELECT appointment_id, start_time, end_time 
		FROM appointments 
		WHERE doctor_id = ? 
		AND appointment_date = ? 
		AND patient_id IS NULL
	`, doctorID, appointmentDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availableSlots []Appointment

	for rows.Next() {
		var slot Appointment
		err := rows.Scan(&slot.AppointmentID, &slot.StartTime, &slot.EndTime)
		if err != nil {
			return nil, err
		}

		availableSlots = append(availableSlots, slot)
	}

	return availableSlots, nil
}
