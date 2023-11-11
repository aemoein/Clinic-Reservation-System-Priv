package main

import (
	"fmt"
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

func ReserveAppointment(appointmentID, patientID int) error {
	// Update the patient ID for the given appointment
	_, err := DB.Exec(`
        UPDATE appointments 
        SET patient_id = ? 
        WHERE appointment_id = ?
    `, patientID, appointmentID)

	if err != nil {
		return err
	}

	return nil
}

// 5- Patient can update his appointment by change the doctor or the slot.
func UpdateAppointment(appointmentID, newDoctorID int, appointmentDate, newStartTime, newEndTime string) error {
	isSlotOccupied, err := IsSlotOccupied(newDoctorID, appointmentDate, newStartTime, newEndTime)
	if err != nil {
		return err
	}

	if isSlotOccupied {
		return fmt.Errorf("the new slot is already occupied")
	}

	_, err = DB.Exec(`
		UPDATE appointments 
		SET doctor_id = ?, start_time = ?, end_time = ? 
		WHERE appointment_id = ?
	`, newDoctorID, newStartTime, newEndTime, appointmentID)

	if err != nil {
		return err
	}

	return nil
}

// 6- Patient can cancel his appointment.
func CancelAppointment(appointmentID int) error {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM appointments 
		WHERE appointment_id = ? AND patient_id IS NOT NULL
	`, appointmentID).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("appointment not found or already canceled")
	}

	_, err = DB.Exec(`
		UPDATE appointments 
		SET patient_id = NULL 
		WHERE appointment_id = ?
	`, appointmentID)

	if err != nil {
		return err
	}

	return nil
}

type AppointmentWithName struct {
	AppointmentID   int    `json:"appointment_id"`
	DoctorID        int    `json:"doctor_id"`
	PatientID       int    `json:"patient_id"`
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DoctorName      string `json:"doctor_name"`
}

// 7- Patients can view all his reservations.
func ViewPatientAppointments(patientID int) ([]AppointmentWithName, error) {
	rows, err := DB.Query(`
		SELECT a.appointment_id, a.doctor_id, a.patient_id, a.appointment_date, a.start_time, a.end_time, u.name as doctor_name
		FROM appointments a
		JOIN users u ON a.doctor_id = u.userid
		WHERE a.patient_id = ?
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentWithName

	for rows.Next() {
		var appointment AppointmentWithName
		err := rows.Scan(&appointment.AppointmentID, &appointment.DoctorID, &appointment.PatientID,
			&appointment.AppointmentDate, &appointment.StartTime, &appointment.EndTime, &appointment.DoctorName)
		if err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}
