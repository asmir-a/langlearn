import React from 'react';

import { whichAuthPageEnum } from '../../auth/auth';

const SignupButton = ({ setWhichAuthPage }) => {
    return (
        <button onClick={() => {
            setWhichAuthPage(whichAuthPageEnum.signup);
        }}>to signup</button>
    )
}

export default SignupButton;