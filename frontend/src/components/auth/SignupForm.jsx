import React from 'react';
import {useState} from 'react';
import * as common from "./../../utilites";

const SIGNUP_ENDPOINT = "/api/signup";

export default function SignupForm ({setAuthState}) {
    const [responseError, setResponseError] = useState(null);

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
    }

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    }

    const handleSubmit = async (event) => {
        event.preventDefault();

        const formData = new FormData();
        formData.append("username", username);
        formData.append("password", password);

        const response = await fetch(SIGNUP_ENDPOINT, {
            method: "POST",
            body: formData,
        });

        if (response.status === common.HTTP_STATUS_OK) {
            setAuthState(common.AUTH_STATE_ENUM.Authed);
            return;
        } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {//todo: make the server return more information in www-unauthorized header or something like that
            setResponseError("provided username or password have an invalid format");
            return;
        } else if (response.status === common.HTTP_STATUS_CONFLICT) {
            setResponseError("username or password already exists.")
            return;
        } else {
            throw new Error("not implemented");
        }
    }

    return (
        <form onSubmit = {handleSubmit}>
            {responseError && <p>{responseError}</p>}
            <input 
                type = "text" 
                placeholder = "username"
                onChange = {handleUsernameChange}
            />
            <input 
                type = "password"
                placeholder = "password"
                onChange={handlePasswordChange}
            />
            <input type = "submit" value = "signup"/>
        </form>
    )
}
