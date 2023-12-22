// src/components/SignUp.js

import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import styles from './Signup.module.css'
const API_BASE_URL = process.env.REACT_APP_API_PORT;

const SignUp = () => {
  const navigate = useNavigate();
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
      const response = await axios.post(`${API_BASE_URL}/signup`, {
        username,
        email,
        password,
        usertype,
      });

      console.log('Response data:', response.data);

      navigate('/signin');
    } catch (error) {
      console.error('Error during sign-up:', error);
      console.log('Error response data:', error.response.data);
      setError('Error signing up. Please try again.');
    }
  };

  return (
    <div className={styles.signuppage}>
      <h2>Sign Up</h2>
      {error && <div style={{ color: 'red', border: '1px solid red', padding: '10px', marginBottom: '10px' }}>{error}</div>}
      <form className={styles.form}>
        <label className={styles.signuplabel}>
          Full Name:
          <input type="text" className={styles.signupinput} value={username} onChange={(e) => setUsername(e.target.value)} />
        </label>
        <br />
        <label className={styles.signuplabel}>
          Email:
          <input type="email" className={styles.signupinput} value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <br />
        <label className={styles.signuplabel}>
          Password:
          <input type="password" className={styles.signupinput} value={password} onChange={(e) => setPassword(e.target.value)} />
        </label>
        <br />

        <label className={styles.radiomain}>
          User Type:
          <div className={styles.radio}>
            <label className={styles.radiolabel}>
              Doctor
              <input
                type="radio"
                value="doctor"
                checked={usertype === 'doctor'}
                onChange={() => setUsertype('doctor')}
                className={styles.radiobtn}
              />
            </label>
            <label className={styles.radiolabel}>
              Patient
              <input
                type="radio"
                value="patient"
                checked={usertype === 'patient'}
                onChange={() => setUsertype('patient')}
                className={styles.radiobtn}
              />
            </label>
          </div>
        </label>
        <br />
        <div className={styles.signupLinks}>
          <button className={styles.signupbutton} type="button" onClick={handleSignUp}>
            Sign Up
          </button>
        </div>

        <div className={styles.signupLinks}>
          <Link to="/signin">Already Signed Up? Login</Link>
        </div>
      </form>
    </div>
  );
};

export default SignUp;
