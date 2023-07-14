import React from 'react';
import axios from 'axios';

import { logoutHandler } from './../../auth/auth';

const LogoutButton = ({ setAuthState }) => {
    return (
        <button
            onClick={(_) => logoutHandler(setAuthState)}
        >logout</button>
    );
};

export default LogoutButton;