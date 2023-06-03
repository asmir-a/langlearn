import React from 'react';
import { useState } from 'react';

import "./App.css";

const LoginForm = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");


  const handleSubmit = async (event) => {
    event.preventDefault();
    let formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);

    let response = await fetch("/api/login", {
      body: formData,
      method: "post"
    });
    let responseJson = await response.text();
    console.log(responseJson);
  }

  return (
    <form action = "/login" method = "post" onSubmit = {handleSubmit}>
      <input 
        type = "text" 
        name = "username" 
        placeholder = "username" 
        onChange = {(event) => setUsername(event.target.value)}
      />
      <input 
        type = "password" 
        name = "password" 
        placeholder = "password" 
        onChange = {event => setPassword(event.target.value)}
      />
      <input type="submit" value = "login"/>
    </form>
  )
}


const App = () => {
  
  return (
    <>
      <LoginForm />
    </>
  );
}

export default App;
