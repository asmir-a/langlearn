import React from 'react';
import { useState, useEffect } from 'react';
import axios from 'axios';

import Navbar from './navbar/Navbar';

import { authStateEnum, selectAuthComponent } from './auth/auth';
import { whichPageEnum, selectPageComponent } from './pages/pages';

import { endpoints, httpCodes } from './utilities';

const selectPageComponentWithAuth = (authState, pageState) => {
    if (authState.state !== authStateEnum.authed) {
        return selectAuthComponent(authState)
    } else {
        return selectPageComponent(pageState)
    }
}

const App = () => {
    const [authState, setAuthState] = useState({
        state: authStateEnum.shouldLogin,
        user: null,
    });
    const [whichPage, setWhichPage] = useState(whichPageEnum.stats);

    useEffect(() => {
        axios.interceptors.response.use((response) => {
            return response;
        }, (error) => {
            if (error.status === httpCodes.unauthorized) {
                setAuthState(prevState => {
                    return {...prevState, state: authStateEnum.shouldLogin}
                });
            } else {
                throw new Error("not handled; response: ", response)
            }
        });
    }, []);

    return (
        <React.Fragment>
            <Navbar
                authState={authState}
                setAuthState={setAuthState}
                whichPage={whichPage}
                setWhichPage={setWhichPage}
            />
        </React.Fragment>
    )
}

export default App;