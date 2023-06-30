import React from 'react';

import { authStateEnum } from '../../auth/auth';

import { endpoints, httpCodes } from '../../utilities';

const LogoutButton = ({ setAuthState }) => {
    const handleClick = async (_) => {
        const responseToLogout = await axios.post(endpoints.logout);
        if (responseToLogout.status === httpCodes.ok) {
            setAuthState(prev => {
                return { ...prev, state: authStateEnum.shouldLogin }
            })
        } else {
            throw new Error(
                "not implemented; status code from response: ",
                responseToLogout.status,
            )
        }
    }

    return (
        <button onClick={handleClick}>logout</button>
    );
};

export default LogoutButton;