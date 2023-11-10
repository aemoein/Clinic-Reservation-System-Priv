import React, { useState } from 'react';
import { Link, useNavigate, NavLink } from 'react-router-dom';
import axios from 'axios';
import './SignIn.css';

const SignIn = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loginStatus, setLoginStatus] = useState(null); // null: initial state, true: success, false: failure
  const navigate = useNavigate();

  const handleSignIn = async () => {
    try {
      // Log email and password before making the API call
      console.log('Email and password being sent:', { email, password });
  
      const response = await axios.post('http://localhost:8081/signin', { email, password });
  
      console.log('Response data:', response.data);
      // Assuming the backend returns user type in the response
      const userType = response.data.user_type;
      const username = response.data.username;
  
      console.log('Name and type recieved:', { userType, username })
      // Set login success status
      setLoginStatus(true);

      if (userType === "doctor") {
        navigate("/doctor/"+username);
      } else if (userType === "patient") {
        navigate("/patient/"+username);
      }
    } catch (error) {
      // Set login failure status
      setLoginStatus(false);
  
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
      {loginStatus === true && <div className="success-message">Login successful!</div>}
      {loginStatus === false && <div className="error-message">Login failed. Please try again.</div>}
      <Link to="/signup">Don't have an account? Sign Up</Link>
    </div>
  );
};

export default SignIn;