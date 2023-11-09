package main

import (
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
