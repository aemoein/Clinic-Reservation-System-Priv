// src/components/SignUp.js

import React, { useState } from 'react';
import axios from 'axios';
import './SignUp.css'

const SignUp = () => {
  const [userType, setUserType] = useState('patient');
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    // Additional fields for patients
    dob: '',
    gender: '',
    // Additional field for doctors
    specialization: '',
  });
  const [successMessage, setSuccessMessage] = useState('');

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    
    try {
      const response = await axios.post('your_backend_endpoint', {
        userType,
        ...formData,
      });

      // Assuming the backend returns a success message
      setSuccessMessage(response.data.message);
    } catch (error) {
      // Handle errors, e.g., display an error message
      console.error('Sign-up failed:', error);
    }
  };

  const handleUserTypeChange = (e) => {
    setUserType(e.target.value);
  };

  return (
    <div className="signup-page">
      <h1>Sign Up</h1>
      <div className="user-type-selection">
        <label>
          <input
            type="radio"
            value="patient"
            checked={userType === 'patient'}
            onChange={handleUserTypeChange}
          />
          Patient
        </label>
        <label>
          <input
            type="radio"
            value="doctor"
            checked={userType === 'doctor'}
            onChange={handleUserTypeChange}
          />
          Doctor
        </label>
      </div>

      {successMessage && <p className="success-message">{successMessage}</p>}

      <form onSubmit={handleFormSubmit}>
        <input
          type="text"
          placeholder="Name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        />
        <input
          type="email"
          placeholder="Email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        />
        <input
          type="password"
          placeholder="Password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
        />
        
        {userType === 'patient' && (
          <div>
            <input
              type="date"
              placeholder="Date of Birth"
              value={formData.dob}
              onChange={(e) => setFormData({ ...formData, dob: e.target.value })}
            />
            <select
              value={formData.gender}
              onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
            >
              <option value="">Select Gender</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
            </select>
          </div>
        )}
        
        {userType === 'doctor' && (
          <input
            type="text"
            placeholder="Specialization"
            value={formData.specialization}
            onChange={(e) => setFormData({ ...formData, specialization: e.target.value })}
          />
        )}
        
        <button type="submit">Sign Up</button>
      </form>
    </div>
  );
};

export default SignUp;