// src/App.js

import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './components/Home';
import SignIn from './components/SignIn';
import SignUp from './components/SignUp';
import Doctor from './components/Doctor';
import Patient from './components/Patient';
import WebSocketComponent from './components/WebSocketComponent';

function App() {
  return (
    <Router>
        <Routes>
          <Route path="/" exact element={<Home />} />
          <Route path="/signin" element={<SignIn />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/doctor/:username/:userid" element={<Doctor />} />
          <Route path="/patient/:username/:userid" element={<Patient />} />
          <Route path="/Kafka/:doctorid" element={<WebSocketComponent />} />
          <Route path="/Kafka" element={<WebSocketComponent />} />
        </Routes>
    </Router>
  );
}

export default App;