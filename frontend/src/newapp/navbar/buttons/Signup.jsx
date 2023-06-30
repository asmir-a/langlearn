import React from 'react';

import { authStateEnum } from '../../auth/auth';

const SignupButton = ({ setAuthState }) => {
    return (
        <button onClick={() => {
            setAuthState(authStateEnum.shouldSignup);
        }}>to signup</button>
    )
}

export default SignupButton;