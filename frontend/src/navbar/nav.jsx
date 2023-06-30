import React from 'react';

import LoginButton from './buttons/Login';
import SignupButton from './buttons/Signup';
import LogoutButton from './buttons/logout';
import StatsButton from './buttons/Stats';
import WordGameButton from './buttons/WordGame';

import { authStateEnum } from './../auth/auth';
import { whichPageEnum } from './../pages/pages';

const getNavbarButtonsWhenNotAuthed = (authState) => {
    const buttons = {
        authButtons: [],
        pageButtons: [],
    };
    switch (authState.state) {
        case authStateEnum.shouldLogin:
            buttons.authButtons.push(SignupButton);
            break;
        case authStateEnum.shouldSignup:
            buttons.authButtons.push(LoginButton);
            break;
        case authStateEnum.authed:
            throw new Error("the state is not supposed to be authed")
        default:
            throw new Error("not handled; authState: ", authState.state)
    }
    return buttons;
}

const getNavbarButtonsWhenAuthed = (whichPage) => {
    const buttons = {
        authButtons: [LogoutButton],
        pageButtons: [],
    };
    switch (whichPage) {
        case whichPageEnum.stats:
            buttons.pageButtons.push(WordGameButton);
            break;
        case whichPageEnum.wordGame:
            buttons.pageButtons.push(StatsButton);
            break;
        default:
            throw new Error("not handled; whichPage: ", whichPage)
    }
    return buttons;
}

export const getNavbarButtons = (authState, whichPage) => {
    if (authState.state !== authStateEnum.authed) {
        return getNavbarButtonsWhenNotAuthed(authState);
    } else {
        return getNavbarButtonsWhenAuthed(whichPage);
    }
}