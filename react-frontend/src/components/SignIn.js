// src/components/SignIn.js
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import './SignIn.css';

const SignIn = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const history = useNavigate();

  const handleSignIn = async () => {
    try {
      const response = await axios.post('http://localhost:8081/signin', { email, password });

      // Assuming the backend returns user type in the response
      const userType = response.data.user_type;

      if (userType === 'doctor') {
        history.push('/doctor-dashboard');
      } else if (userType === 'patient') {
        history.push('/patient-dashboard');
      }
    } catch (error) {
      // Handle authentication errors, e.g., show an error message
      console.error('Authentication failed:', error);
    }
  };

  return (
    <div className="signin-page">
      <h1>Sign In</h1>
      <div className="form">
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button onClick={handleSignIn}>Sign In</button>
      </div>
      <Link to="/signup">Don't have an account? Sign Up</Link>
    </div>
  );
};

export default SignIn;