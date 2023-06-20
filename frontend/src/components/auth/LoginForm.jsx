import React from 'react';
import {useState, useEffect} from 'react';
import * as common from "../../utilites";

const loginEndpoint = "/api/login";

export default function LoginForm ({setAuthInfo}) {
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
            const response = await fetch("/api/is-authed");
            const userData = await response.json();
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.Authed,
                user: userData
            });
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