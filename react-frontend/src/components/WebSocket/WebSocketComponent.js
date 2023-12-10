import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import './Ws.css';

const WebSocketComponent = () => {
  const { doctorid } = useParams();
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8081/ws');

    socket.onopen = () => {
      console.log('WebSocket connection opened');
    };

    socket.onmessage = (event) => {
      console.log('Raw message received:', event.data);

      const rawMessage = event.data;
      console.log('Raw message received:', rawMessage);

      const trimmedMessage = rawMessage.replace(/\u0000/g, '');
      try {
        var newMessage = JSON.parse(trimmedMessage);
        setMessages((prevMessages) => [...prevMessages, newMessage]);
      } catch (error) {
        console.error('Error parsing JSON:', error);
        console.error('Raw JSON string:', event.data);
      }
    };

    socket.onclose = (event) => {
      console.log('WebSocket connection closed:', event);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);

      if (error && error.message) {
        console.error('Error message:', error.message);
      }
    };

    return () => {
      if (socket.readyState === 0) {
        socket.close();
      }
    };
  }, []);

  console.log("message recieved: ", messages)
  console.log("id used: ", doctorid)

  return (
    <div className="container">
      <h1>WebSocket Messages</h1>
      <Link to="/" className="button">
        Log Out
      </Link>
      <ul>
        {messages
          .filter(message => (doctorid ? parseInt(message.doctorId, 10) === parseInt(doctorid, 10) : true))
          .map((message, index) => (
            <li key={index}>{JSON.stringify(message)}</li>
          ))}
      </ul>
    </div>
  );
};

export default WebSocketComponent;