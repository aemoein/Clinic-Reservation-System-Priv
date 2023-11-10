import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import './Doctor.css';

const Doctor = () => {
  const { username, userid } = useParams();

  const [slots, setSlots] = useState([]);
  const [newSlot, setNewSlot] = useState({
    doctor_id: Number(userid),
    patient_id: 0,
    appointment_date: '',
    start_time: '',
    end_time: '',
  });

  const fetchDoctorSlots = async () => {
    try {
      const response = await axios.get('http://localhost:8081/slots/view', {
        params: { doctorid: userid },
      });

      if (Array.isArray(response.data)) {
        console.log("Received data:", response.data);
        setSlots(response.data);
      } else {
        console.error('Invalid data format for doctor slots:', response.data);
      }
    } catch (error) {
      console.error('Error fetching doctor slots:', error);
    }
  };

  useEffect(() => {
    fetchDoctorSlots();
  }, [username]);

  const handleEdit = (slotId) => {
    console.log(`Edit slot with ID ${slotId}`);
    // Implement edit logic for the selected slot
  };

  const handleCancel = (slotId) => {
    console.log(`Cancel slot with ID ${slotId}`);
    // Implement cancel logic for the selected slot
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewSlot((prevSlot) => ({
      ...prevSlot,
      [name]: value,
    }));
  };  

  const handleAddSlot = async () => {
    try {
      // Send a request to the server to add the new slot
      console.log('New slot sent:', newSlot);
      await axios.post('http://localhost:8081/slots/add', newSlot);
      // Refresh the list of slots after adding a new one
      fetchDoctorSlots();
      // Clear the input fields
      setNewSlot({
        doctor_id: Number(userid),
        patient_id: 0,
        appointment_date: '',
        start_time: '',
        end_time: '',
      });
    } catch (error) {
      console.error('Error adding slot:', error);
    }
  };

  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: Doctor</p>
      <h2>My Slots</h2>
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {slots.map((slot) => (
            <tr key={slot.appointment_id}>
              <td>{slot.appointment_date}</td>
              <td>{slot.start_time}</td>
              <td>{slot.end_time}</td>
              <td>
                <button onClick={() => handleEdit(slot.id)}>Edit</button>
                <button onClick={() => handleCancel(slot.id)}>Cancel</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Form to add a new slot */}
      <h2>Add New Slot</h2>
      <form>
        <label>
          Date:
        </label>
        <input type="date" name="appointment_date" value={newSlot.appointment_date} onChange={handleInputChange} />
        <label>
          Start Time:
        </label>
        <input type="time" name="start_time" value={newSlot.start_time} onChange={handleInputChange} />
        <label>
          End Time:
        </label>
        <input type="time" name="end_time" value={newSlot.end_time} onChange={handleInputChange} />
        <button type="button" onClick={handleAddSlot}>
          Add Slot
        </button>
      </form>
    </div>
  );
};

export default Doctor;