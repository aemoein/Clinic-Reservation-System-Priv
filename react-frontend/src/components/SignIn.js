import React, { useState } from 'react';
import { Link, useNavigate, NavLink } from 'react-router-dom';
import axios from 'axios';
import './SignIn.css';

const SignIn = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loginStatus, setLoginStatus] = useState(null);
  const navigate = useNavigate();

  const handleSignIn = async () => {
    try {
      console.log('Email and password being sent:', { email, password });
  
      const response = await axios.post('http://localhost:8081/signin', { email, password });
  
      console.log('Response data:', response.data);

      const userid = response.data.userid;
      const userType = response.data.user_type;
      const username = response.data.username;
  
      console.log('Name, id and type recieved:', { userid, userType, username })

      setLoginStatus(true);

      if (userType === "doctor") {
        navigate( "/doctor/" + username + "/" + userid );
      } else if (userType === "patient") {
        navigate( "/patient/" + username + "/" + userid );
      }
    } catch (error) {
      setLoginStatus(false);
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