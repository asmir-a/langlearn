import React from 'react';
import * as common from '../../utilites';

const LoginNavButton = ({setAuthState}) => {
    const handleClick = (_) => {
        setAuthState(common.AUTH_STATE_ENUM.ShouldLogin);
    }
    return (
        <button onClick={handleClick}>to login</button>
    );
}

const SignupNavButton = ({setAuthState}) => {
    const handleClick = (_) => {
        setAuthState(common.AUTH_STATE_ENUM.ShouldSignup);
    }
    return (
        <button onClick={handleClick}>to signup</button>
    );
}

const LOGOUT_ENDPOINT = "/api/logout";
function LogoutNavButton({setAuthState}) {
    const handleClick = async (event) => {
        const response = await fetch(LOGOUT_ENDPOINT, {
            method: "POST",
        });
        if (response.status === common.HTTP_STATUS_OK) {
            setAuthState(common.AUTH_STATE_ENUM.ShouldLogin);
            return;
        } else {
            throw new Erorr("not implemented");
        }
    }

    return (
        <button onClick = {handleClick}>logout</button>
    );
}

export default function NavBar({authState, setAuthState}) {
    const selectAuthComponents = () => {
        if (authState !== common.AUTH_STATE_ENUM.Authed) {
            return (
                <React.Fragment>
                    <LoginNavButton setAuthState={setAuthState}/>
                    <SignupNavButton setAuthState={setAuthState}/>
                </React.Fragment>
            );
        } else {
            return <LogoutNavButton setAuthState={setAuthState}/>;
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