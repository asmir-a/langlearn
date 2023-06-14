import React from 'react';
import {useState, useEffect} from 'react';
import * as common from "../../utilites";

const loginEndpoint = "/api/login";

export default function LoginForm ({setAuthState}) {
    const [responseError, setResponseError] = useState(null);

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleUsernameChange = (event) => {
        event.preventDefault();
        setUsername(event.target.value);
    }

    const handlePasswordChange = (event) => {
        event.preventDefault();
        setPassword(event.target.value);
    }

    const handleSubmit = async (event) => {
        event.preventDefault();

        const formData = new FormData();
        formData.append("username", username);
        formData.append("password", password);

        const response = await fetch(loginEndpoint, {
            method: "POST",
            body: formData,
        });

        if (response.status === 200) {
            setAuthState(common.AUTH_STATE_ENUM.Authed);
            return;
        } else if (response.status == common.HTTP_STATUS_UNAUTHORIZED) {
            setResponseError("wrong login or password")
            return
        } else {
            throw new Error("not implemented");
        }
    }

    return (
        <form onSubmit = {handleSubmit}>
            {responseError && <p>{responseError}</p>}
            <input type="text" placeholder="username" onChange = {handleUsernameChange}/>
            <input type="password" placeholder="password" onChange = {handlePasswordChange}/>
            <input type = "submit" value = "login"/>
        </form>
    )
}