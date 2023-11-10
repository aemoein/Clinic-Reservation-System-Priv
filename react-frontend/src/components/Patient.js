import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import './Patient.css';

const Patient = () => {
  const { username, userid } = useParams();

  const [doctors, setDoctors] = useState([]);
  const [selectedDoctor, setSelectedDoctor] = useState('');
  const [doctorSlots, setDoctorSlots] = useState([]);
  const [selectedSlot, setSelectedSlot] = useState('');
  const [patientAppointments, setPatientAppointments] = useState([]);

  const fetchDoctors = async () => {
    try {
      const response = await axios.get('http://localhost:8081/doctors');
      if (Array.isArray(response.data)) {
        setDoctors(response.data);
      } else {
        console.error('Invalid data format for doctors:', response.data);
      }
    } catch (error) {
      console.error('Error fetching doctors:', error);
    }
  };

  const fetchPatientAppointments = async () => {
    try {
      const response = await axios.get('http://localhost:8081/appointments/view', {
        params: { patientid: userid },
      });

      if (Array.isArray(response.data)) {
        setPatientAppointments(response.data);
      } else {
        console.error('Invalid data format for patient appointments:', response.data);
      }
    } catch (error) {
      console.error('Error fetching patient appointments:', error);
    }
  };

  useEffect(() => {
    fetchPatientAppointments();
  }, [username]);

  const handleEdit = async (appointmentId) => {
    try {
        // Send a request to the server to cancel the appointment
        console.log(`Cancel appointment with ID ${appointmentId}`);
        await axios.post('http://localhost:8081/appointments/cancel', {
          appointment_id: appointmentId,
          patient_id: Number(userid),
        });
  
        // Refresh the list of patient appointments after cancellation
        fetchPatientAppointments();
      } catch (error) {
        console.error('Error canceling appointment:', error);
      }
  };

  const handleCancel = async (appointmentId) => {
    try {
      // Send a request to the server to cancel the appointment
      console.log(`Cancel appointment with ID ${appointmentId}`);
      await axios.post('http://localhost:8081/appointments/cancel', {
        appointment_id: appointmentId,
        patient_id: Number(userid),
      });

      // Refresh the list of patient appointments after cancellation
      fetchPatientAppointments();
    } catch (error) {
      console.error('Error canceling appointment:', error);
    }
  };

  const fetchDoctorSlots = async () => {
    try {
      const response = await axios.get('http://localhost:8081/slots/view', {
        params: { doctorid: selectedDoctor },
      });

      if (Array.isArray(response.data)) {
        setDoctorSlots(response.data);
      } else {
        console.error('Invalid data format for doctor slots:', response.data);
      }
    } catch (error) {
      console.error('Error fetching doctor slots:', error);
    }
  };

  const handleDoctorChange = (e) => {
    const selectedDoctorId = e.target.value;
    setSelectedDoctor(selectedDoctorId);
    // Fetch available slots when a doctor is selected
    fetchDoctorSlots(selectedDoctorId);
  };

  const handleSlotChange = (e) => {
    const selectedSlotId = e.target.value;
    setSelectedSlot(selectedSlotId);
  };

  const handleReserve = async () => {
    try {
      // Send a request to the server to reserve the selected slot
      await axios.post('http://localhost:8081/appointments/reserve', {
        appointment_id: selectedSlot,
        patient_id: Number(userid),
      });

      // Refresh the list of patient appointments after reservation
      fetchPatientAppointments();
    } catch (error) {
      console.error('Error reserving appointment:', error);
    }
  };

  useEffect(() => {
    fetchDoctors();
    fetchPatientAppointments();
  }, [username]);

  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: Patient</p>
      <h2>Your Appointments</h2>
      <table>
      <thead>
          <tr>
            <th>Date</th>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Doctor Name</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {patientAppointments.map((appointment) => (
            <tr key={appointment.appointment_id}>
              <td>{appointment.appointment_date}</td>
              <td>{appointment.start_time}</td>
              <td>{appointment.end_time}</td>
              <td>{appointment.doctor_name}</td>
              <td>
                <button onClick={() => handleEdit(appointment.appointment_id)}>Edit</button>
                <button onClick={() => handleCancel(appointment.appointment_id)}>Cancel</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <h2>Make a Reservation</h2>
      <form>
        <label>
          Select a Doctor:
          <select value={selectedDoctor} onChange={handleDoctorChange}>
            <option value="" disabled>Select a doctor</option>
            {doctors.map((doctor) => (
              <option key={doctor.doctor_id} value={doctor.doctor_id}>
                {doctor.doctor_name}
              </option>
            ))}
          </select>
        </label>

        <label>
          Select a Slot:
          <select value={selectedSlot} onChange={handleSlotChange}>
            <option value="" disabled>Select a slot</option>
            {doctorSlots.map((slot) => (
              <option key={slot.appointment_id} value={slot.appointment_id}>
                {slot.appointment_date} - {slot.start_time} to {slot.end_time}
              </option>
            ))}
          </select>
        </label>

        <button type="button" onClick={handleReserve}>
          Reserve Appointment
        </button>
      </form>
    </div>
  );
};

export default Patient;