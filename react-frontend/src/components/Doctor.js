import React from 'react';
import { useParams } from 'react-router-dom';

const Doctor = () => {
    let { username } = useParams();

    console.log (username);

  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: Doctor</p>
      {/* Add other content for the doctor page */}
    </div>
  );
};

export default Doctor;
