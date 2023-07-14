import React from 'react';

import { whichAuthPageEnum } from '../../auth/auth';


const LoginButton = ({ setWhichAuthPage }) => {
    return (
        <button onClick={() => {
            setWhichAuthPage(whichAuthPageEnum.login)
        }}>to login</button>
    )
}

export default LoginButton;