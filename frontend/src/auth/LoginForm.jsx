import React from 'react';
import { useState, useRef } from 'react';

import { loginHandler } from './auth';

import "./Form.css";

const LoginForm = ({ setAuthState }) => {
    const [errorMessage, setErrorMessage] = useState(null);
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const formContainer = useRef(null);

    const handleSubmit = async (event) => {
        event.preventDefault();

        const formData = new FormData();
        formData.append("username", username);
        formData.append("password", password);

        const responseErrorMessage = await loginHandler(setAuthState, formData);
        if (responseErrorMessage != null) {
            setErrorMessage(responseErrorMessage);
            formContainer.current.className = " with-error";
        }

        setUsername("");
        setPassword("");
    }

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        formContainer.current.className = "";
    }

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
        formContainer.current.className = "";
    }

    return (
        //form should be wrapped in main; the main prolly should be some sort of a container for auth pages handling the auth logic as well
        <form onSubmit={handleSubmit} ref={formContainer}>
            <h1>log in to your account</h1>
            {errorMessage && <h2 className="error-message">{errorMessage}</h2>} {/* this works only if the error message is null when there is no error: "" does not work with conditional rendering that uses && */}
            <label htmlFor="username">
                username:
                <input
                    value={username}
                    onChange={handleUsernameChange}
                    name="username"
                    required
                />
            </label>
            <label htmlFor="password">
                password:
                <input
                    type="password"
                    value={password}
                    onChange={handlePasswordChange}
                    name="password"
                    required
                />
            </label>
            <input type="submit" value="login" />
        </form>
    );
}

export default LoginForm;