import React from 'react';
import { useState, useEffect } from 'react';

import Navbar from './navbar/Navbar';

import "./App.css";

import {
    authStateEnum,
    whichAuthPageEnum,
    makeRequestToUser,
    interceptUnauthorizedResponses,
    selectAuthPage,
    stopInterceptingUnauthorizedResponses,
} from './auth/auth';
import { selectContentPage, whichContentPageEnum } from './pages/pages';

const selectPage = (authState, setAuthState, whichAuthPage, whichContentPage) => {
    switch (authState.state) {
        case authStateEnum.authed:
            return selectContentPage(whichContentPage);
        case authStateEnum.notAuthed:
            return selectAuthPage(setAuthState, whichAuthPage);
        case authStateEnum.loading:
            return () => <p>loading</p>//todo: need to handle this properly and to consider how does it affect navbar
        default:
            throw new Error("not handled; authState: ", authState)
    }
}

const useAuthState = (authState) => {
    const [whichAuthPage, setWhichAuthPage] = useState(whichAuthPageEnum.login);
    const [whichContentPage, setWhichContentPage] = useState(whichContentPageEnum.stats);

    useEffect(() => {
        switch (authState.state) {
            case authStateEnum.authed: case authStateEnum.notAuthed:
                setWhichAuthPage(whichAuthPageEnum.login);//it might be better to call this from the loginfrom component; the branch above also logically can call the setwhichpage with first pages to show (login and stats). however, maybe some higher level abstraction is kinda needed. we migth have a auth page container component. that way we will be able to remove the axios interceptor when it unmounts, but for now it is okay
                setWhichContentPage(whichContentPageEnum.stats);
                break;
            case authStateEnum.loading: default:
                break;
        }
    }, [authState.state]);

    return [whichAuthPage, setWhichAuthPage, whichContentPage, setWhichContentPage];
}

const useInterceptorForAuthed = (authState, setAuthState) => {
    //may we could have an interceptor saved in here.
    useEffect(() => {
        switch (authState.state) {
            case authStateEnum.authed:
                const interceptor = interceptUnauthorizedResponses(setAuthState);
                setAuthState(prev => {
                    return {
                        ...prev,
                        interceptorForAuthed: interceptor,
                    }
                })
            case authState.notAuthed:
                stopInterceptingUnauthorizedResponses(authState.interceptorForAuthed);
                setAuthState(prev => {
                    return {
                        ...prev,
                        interceptorForAuthed: interceptor,
                    }
                })
            default:
                break;
        }
    }, [authState.state]);
}

const App = () => {
    const [authState, setAuthState] = useState({
        state: authStateEnum.loading,
        user: null,
        interceptorForAuthed: null
    });
    const [
        whichAuthPage, 
        setWhichAuthPage, 
        whichContentPage, 
        setWhichContentPage
    ] = useAuthState(authState);

    useEffect(() => {
        makeRequestToUser(setAuthState);
    }, []);

    useInterceptorForAuthed(authState, setAuthState);

    const Page = selectPage(
        authState,
        setAuthState,
        whichAuthPage,
        whichContentPage
    );//this should not be done; passing props become more tedious. There should be a Content functional components defined.

    return (
        <React.Fragment>
            <Navbar //todo: 5 parameters to pass seems too much. Other alternatives are described in the file that implements navbar
                authState={authState}
                setAuthState={setAuthState}
                whichAuthPage={whichAuthPage}
                setWhichAuthPage={setWhichAuthPage}
                whichContentPage={whichContentPage}
                setWhichContentPage={setWhichContentPage}
            />
            <Page user = {authState.user}/>
        </React.Fragment>
    );
}

export default App;