import React from 'react';
import axios from 'axios';

import LoginForm from './LoginForm';
import SignupForm from './SignupForm';

import { endpoints, httpCodes } from '../utilities';

export const authStateEnum = {
    authed: 'authed',
    notAuthed: 'notAuthed',
    loading: 'loading',
}

export const whichAuthPageEnum = {
    login: 'login',
    signup: 'signup',
}

export const makeRequestToUser = async (setAuthState) => {
    try {
        const response = await axios.get(endpoints.user)
        if (response.status === httpCodes.ok) {
            const user = response.data;
            setAuthState({
                state: authStateEnum.authed,
                user: user,
            })
        } else {
            throw new Error("not handled; response: ", response);
        }
    } catch (error) {
        if (error.response.status === httpCodes.unauthorized) {
            setAuthState({
                state: authStateEnum.notAuthed,
                user: null,
            })
        } else {
            throw new Error(`not handled; the error: ${error}`)
        }
    }
}

export const interceptUnauthorizedResponses = (setAuthState) => {
    return axios.interceptors.response.use(response => {
        return response;
    }, error => {
        if (error.status === httpCodes.unauthorized) {
            setAuthState({
                state: authStateEnum.notAuthed,
                user: null,
            });
            return error;
        } else {
            throw new Error(`not handled; error: ${error}`)
        }
    });
}

export const stopInterceptingUnauthorizedResponses = (interceptor) => {
    axios.interceptors.response.eject(interceptor);
}

export const loginHandler = async (setAuthState, formData) => {
    try {
        await axios.postForm(endpoints.login, formData);
        const response = await axios.get(endpoints.user);
        const user = response.data;
        setAuthState(prev => {
            return {//todo: think about the state not changing immediately. so the user still stays on the login page for some small time. what if the user presses the login button multiple times? maybe i should disable the buttons for some time after it is pressed
                ...prev,
                state: authStateEnum.authed,
                user: user,
            }
        })
    } catch (error) {
        switch (error.response.status) {
            case httpCodes.unauthorized:
                return "username or password is incorrect"
            default:
                throw new Error(`not handled; error: ${error}`)
        }
    }
}

export const signupHandler = async (setAuthState, formData) => {
    try {
        await axios.postForm(endpoints.signup, formData);
        const response = await axios.get(endpoints.user);
        const user = response.data;
        setAuthState(prev => {
            return {
                ...prev,
                state: authStateEnum.authed,
                user: user,
            }
        })
    } catch (error) {//todo: for signup and login, need to handle when error is received
        switch (error.response.status) {
            case httpCodes.unauthorized:
                return "username or password is invalid"
            case httpCodes.conflict:
                return "username already exists\nplease choose another one"
            default:
                throw new Error(`not handled; error: ${error}`)
        }
    }
}

export const logoutHandler = async (setAuthState) => {
    try {
        await axios.post(endpoints.logout);
        setAuthState(prev => {
            return {
                ...prev,
                state: authStateEnum.notAuthed,
                user: null
            }
        });
    } catch (error) {//todo: need to handle the error condition somehow
        throw new Error(`not handled; error: ${error}`);
    }
}

export const selectAuthPage = (setAuthState, whichAuthPage) => {
    switch (whichAuthPage) {
        case whichAuthPageEnum.login:
            return () => <LoginForm setAuthState = {setAuthState}/>
        case whichAuthPageEnum.signup:
            return () => <SignupForm setAuthState={setAuthState}/>
        default:
            throw new Error("not handled; authState: ", whichAuthPage)
    }
}

export const usernameFormatGuidelines = "username needs to only have letters and numbers, and be at least 4 characters long; not special symbols are allowed";
export const validateUsername = (username) => {
    return /^[A-Za-z][A-za-z0-9]/.test(username) &&
            username.length >= 4
}

export const passwordFormatGuidelines = "password needs to include at least one letter and at least one number and at least one special symbol and be at least 7 characters long";
export const validatePassword = (password) => {
    return /[A-Z]/       .test(password) &&
           /[a-z]/       .test(password) &&
           /[0-9]/       .test(password) &&
           /[^A-Za-z0-9]/.test(password) &&
           password.length >= 7 
}