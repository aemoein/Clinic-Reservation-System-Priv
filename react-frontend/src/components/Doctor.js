import React from 'react';

const Doctor = ({ username, userType }) => {
  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: {userType}</p>
      {/* Add other content for the doctor page */}
    </div>
  );
};

export default Doctor;
