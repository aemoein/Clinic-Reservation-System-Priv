import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import './Doctor.css';

const Doctor = () => {
    let { username, userid } = useParams();

    console.log (username, userid);

    const [slots, setSlots] = useState([]);

  useEffect(() => {
    const fetchDoctorSlots = async () => {
      try {
        const response = await axios.get(`http://localhost:8081/slots/view?username=${userid}`);
        if (Array.isArray(response.data)) {
            setSlots(response.data);
          } else {
            console.error('Invalid data format for doctor slots:', response.data);
          }
      } catch (error) {
        console.error('Error fetching doctor slots:', error);
      }
    };

    // Fetch doctor's slots when the component mounts
    fetchDoctorSlots();
  }, [username]);

  const handleEdit = (slotId) => {
    // Implement edit logic for the selected slot
    console.log(`Edit slot with ID ${slotId}`);
  };

  const handleCancel = (slotId) => {
    // Implement cancel logic for the selected slot
    console.log(`Cancel slot with ID ${slotId}`);
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
            <tr key={slot.id}>
              <td>{slot.date}</td>
              <td>{slot.startTime}</td>
              <td>{slot.endTime}</td>
              <td>
                <button onClick={() => handleEdit(slot.id)}>Edit</button>
                <button onClick={() => handleCancel(slot.id)}>Cancel</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Doctor;
