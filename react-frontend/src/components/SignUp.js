// src/components/SignUp.js

import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import './SignUp.css'

const SignUp = () => {
  const history = useNavigate();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [usertype, setUsertype] = useState('');
  const [error, setError] = useState('');

  const handleSignUp = async () => {
    console.log('Data sent to the backend:', {
      username,
      email,
      password,
      usertype,
    });

    try {
      const response = await axios.post('http://localhost:8081/signup', {
        username,
        email,
        password,
        usertype,
      });

      console.log('Response data:', response.data);

      history.push('/login');
    } catch (error) {
      console.error('Error during sign-up:', error);
      console.log('Error response data:', error.response.data);
      setError('Error signing up. Please try again.');
    }
  };

  return (
    <div className="signup-page">
      <h2>Sign Up</h2>
      {error && <div style={{ color: 'red', border: '1px solid red', padding: '10px', marginBottom: '10px' }}>{error}</div>}
      <form>
        <label>
          User Name:
          <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} />
        </label>
        <br />
        <label>
          Email:
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <br />
        <label>
          Password:
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </label>
        <br />
        <label className="radiomain">
          User Type:
          <div className="radio">
            <label>
              <input
                type="radio"
                value="doctor"
                checked={usertype === 'doctor'}
                onChange={() => setUsertype('doctor')}
                className="radiobtn"
              />
              doctor
            </label>
            <label>
              <input
                type="radio"
                value="patient"
                checked={usertype === 'patient'}
                onChange={() => setUsertype('patient')}
                className="radiobtn"
              />
              patient
            </label>
          </div>
        </label>
        <br />
        <button type="button" onClick={handleSignUp}>
          Sign Up
        </button>
      </form>
      <Link to="/signin">Already Signed Up? Login</Link>
    </div>
  );
};

export default SignUp;
