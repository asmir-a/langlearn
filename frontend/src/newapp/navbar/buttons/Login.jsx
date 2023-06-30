import React from 'react';

import { authStateEnum } from '../../auth/auth';

const LoginButton = ({ setAuthState }) => {
    return (
        <button onClick={() => {
            setAuthState(authStateEnum.shouldLogin);
        }}>to login</button>
    )
}

export default LoginButton;