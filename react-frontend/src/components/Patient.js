import React from 'react';
import { useParams } from 'react-router-dom';

const Patient = () => {
    let { username } = useParams();

    console.log (username);

  return (
    <div>
      <h1>Hello, {username}</h1>
      <p>User Type: Patient</p>
      {/* Add other content for the patient page */}
    </div>
  );
};

export default Patient;