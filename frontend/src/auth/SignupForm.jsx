import React from 'react';
import { useState, useRef } from 'react';

import { 
    signupHandler, 
    usernameFormatGuidelines, 
    validateUsername, 
    passwordFormatGuidelines, 
    validatePassword 
} from './auth';

import "./Form.css";

const SignupForm = ({ setAuthState }) => {
    const [errorMessage, setErrorMessage] = useState(null);
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const formContainer = useRef(null);

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!validateUsername(username)) {
            setErrorMessage(usernameFormatGuidelines);
            formContainer.current.className = "with-error";
            return;
        } else if (!validatePassword(password)) {
            setErrorMessage(passwordFormatGuidelines);
            formContainer.current.className = "with-error";
            return;
        }

        const formData = new FormData();
        formData.append("username", username);
        formData.append("password", password);

        const responseErrorMessage = await signupHandler(setAuthState, formData);//todo: await or try/catch might be needed
        setErrorMessage(responseErrorMessage);

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
        <form onSubmit={handleSubmit} ref={formContainer}>
            <h1>create a new account</h1>
            {errorMessage && <h2 className="error-message">{errorMessage}</h2>} {/* this works only if the errorMessage is null when initially set and when returned from the handler. conditional rendering that uses && does not work if the error can be ""*/}
            <label htmlFor="username">username:
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
            <input type="submit" value="signup" />
        </form>
    )
}

export default SignupForm;