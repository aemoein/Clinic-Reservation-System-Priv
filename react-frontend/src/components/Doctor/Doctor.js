import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';
import styles from './Doctor.module.css';
const API_BASE_URL = process.env.REACT_APP_API_PORT;

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
      const response = await axios.post(`${API_BASE_URL}/slots/view`, {
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
      await axios.post(`${API_BASE_URL}/slots/add`, newSlot);
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
    <div className={styles.doctordiv}>
      <header className={styles.doctorheader}>
        <h1>Hello, {username}</h1>
        <h2>User Type: Doctor</h2>
        <div>
        <Link to="/" className="button">Log Out</Link>
        <Link to={`/Kafka/${userid}`} className="button">Messages</Link>
        </div>
      </header>
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
                <button className={styles.doctorbutton} onClick={() => handleEdit(slot.id)}>Edit</button>
                <button className={styles.doctorbutton} onClick={() => handleCancel(slot.id)}>Cancel</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className={styles.addslot}>
        <h2>Add New Slot</h2>
        <form className={styles.doctorform}>
          <label className={styles.doctorlabel}>
            Date:
            <input type="date" className={styles.doctorinput} name="appointment_date" value={newSlot.appointment_date} onChange={handleInputChange} />
          </label>
          <label className={styles.doctorlabel}>
            Start Time:
            <input type="time" className={styles.doctorinput} name="start_time" value={newSlot.start_time} onChange={handleInputChange} />
          </label>
          <label className={styles.doctorlabel}>
            End Time:
            <input type="time" className={styles.doctorinput} name="end_time" value={newSlot.end_time} onChange={handleInputChange} />
          </label>
          <button type="button" className={styles.addslotbutton} onClick={handleAddSlot}>
            Add Slot
          </button>
        </form>
      </div>
    </div>
  );
};

export default Doctor;