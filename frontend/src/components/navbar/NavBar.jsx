import React from 'react';
import * as common from '../../utilites';
import * as componentSelector from './../mainlogic/ComponentSelector';

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

function StatsNavButton({setPage}) {
    const handleClick = async (event) => {
        setPage(componentSelector.pagesEnum.stats)
    }
    return (
        <button onClick={handleClick}>stats</button>
    )
}

export default function NavBar({authInfo, setAuthInfo, setPage}) {
    const selectAuthComponents = () => {
        if (authInfo.authState !== common.AUTH_STATE_ENUM.Authed) {
            return (
                <React.Fragment>
                    <StatsNavButton setPage={setPage} />
                    <LoginNavButton setAuthInfo={setAuthInfo} />
                    <SignupNavButton setAuthInfo={setAuthInfo} />
                </React.Fragment>
            );
        } else {
            return <LogoutNavButton setAuthInfo={setAuthInfo}/>;
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