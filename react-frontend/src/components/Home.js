// src/components/HomePage.js

import React from 'react';
import { Link } from 'react-router-dom';
import './Home.css';

const HomePage = () => {
  return (
    <div className="home-page">
      <h1>Welcome to Clinical Reservation</h1>
      <p>Your Trusted Medical Appointment System</p>
      <div className="button-container">
        <Link to="/signin" className="button">Sign In</Link>
        <Link to="/signup" className="button">Sign Up</Link>
      </div>
      <Link to="/kafka" className="button">All Messages</Link>
    </div>
  );
};

export default HomePage;