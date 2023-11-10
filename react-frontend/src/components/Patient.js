import React from 'react';
import { useLocation } from 'react-router-dom';

const Patient = ({ username, userType }) => {
  const location = useLocation();

  // If the user data is not passed from the login page, redirect to the login page
  if (!location.state) {
    return window.location.replace('/login');
  }

  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: {userType}</p>
      {/* Add other content for the patient page */}
    </div>
  );
};

export default Patient;