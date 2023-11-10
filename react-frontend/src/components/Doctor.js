import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import './Doctor.css';

const Doctor = () => {
  let { username, userid } = useParams();

  console.log(username, userid);

  const [slots, setSlots] = useState([]);
  const [newSlot, setNewSlot] = useState({
    date: '',
    startTime: '',
    endTime: '',
  });

  const fetchDoctorSlots = async () => {
    try {
      const response = await axios.get('http://localhost:8081/slots/view', {
        params: { doctorid: userid },
      });
  
      if (Array.isArray(response.data)) {
        setSlots(response.data);
      } else {
        console.error('Invalid data format for doctor slots:', response.data);
      }
    } catch (error) {
      console.error('Error fetching doctor slots:', error);
    }
  };  

  useEffect(() => {
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
      await axios.post(`http://localhost:8081/slots/add?doctorID=${userid}`, newSlot);
      // Refresh the list of slots after adding a new one
      fetchDoctorSlots();
      // Clear the input fields
      setNewSlot({ date: '', startTime: '', endTime: '' });
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

      {/* Form to add a new slot */}
      <h2>Add New Slot</h2>
      <form>
        <label>
          Date:
          <input type="date" name="date" value={newSlot.date} onChange={handleInputChange} />
        </label>
        <label>
          Start Time:
          <input type="time" name="startTime" value={newSlot.startTime} onChange={handleInputChange} />
        </label>
        <label>
          End Time:
          <input type="time" name="endTime" value={newSlot.endTime} onChange={handleInputChange} />
        </label>
        <button type="button" onClick={handleAddSlot}>
          Add Slot
        </button>
      </form>
    </div>
  );
};

export default Doctor;