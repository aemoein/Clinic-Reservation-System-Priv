import React from 'react';
import { useParams } from 'react-router-dom';

const Patient = () => {
    let { username } = useParams();

    console.log (username);

  return (
    <div>
      <h1>Hello, {username}</h1>
      <h2>User Type: Patient</h2>
      {/* Add other content for the patient page */}
    </div>
  );
};

export default Patient;