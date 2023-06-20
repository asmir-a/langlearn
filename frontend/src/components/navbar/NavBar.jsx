import React from 'react';
import * as common from '../../utilites';

const LoginNavButton = ({setAuthInfo}) => {
    const handleClick = (_) => {
        setAuthInfo({
            authState: common.AUTH_STATE_ENUM.ShouldLogin,
            user: null
        });
    }
    return (
        <button onClick={handleClick}>to login</button>
    );
}

const SignupNavButton = ({setAuthInfo}) => {
    const handleClick = (_) => {
        setAuthInfo({
            authState: common.AUTH_STATE_ENUM.ShouldSignup,
            user: null
        });
    }
    return (
        <button onClick={handleClick}>to signup</button>
    );
}

const LOGOUT_ENDPOINT = "/api/logout";
function LogoutNavButton({setAuthInfo}) {
    const handleClick = async (event) => {
        const response = await fetch(LOGOUT_ENDPOINT, {
            method: "POST",
        });
        if (response.status === common.HTTP_STATUS_OK) {
            const response = await fetch("/api/is-authed");
            const userData = await response.json();
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.ShouldLogin,
                user: userData
            });
            return;
        } else {
            throw new Erorr("not implemented");
        }
    }

    return (
        <button onClick = {handleClick}>logout</button>
    );
}

export default function NavBar({authInfo, setAuthInfo}) {
    const selectAuthComponents = () => {
        if (authInfo.authState !== common.AUTH_STATE_ENUM.Authed) {
            return (
                <React.Fragment>
                    <LoginNavButton setAuthState={setAuthInfo}/>
                    <SignupNavButton setAuthState={setAuthInfo}/>
                </React.Fragment>
            );
        } else {
            return <LogoutNavButton setAuthState={setAuthInfo}/>;
        }
    }
    return (
        <nav>
            <ul>
                {selectAuthComponents()}
            </ul>
        </nav>
    );
}