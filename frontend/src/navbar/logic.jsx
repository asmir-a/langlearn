import React from 'react';

import LoginButton from './buttons/Login';
import SignupButton from './buttons/Signup';
import LogoutButton from './buttons/logout';
import StatsButton from './buttons/Stats';
import WordGameButton from './buttons/WordGame';

import { authStateEnum, whichAuthPageEnum } from '../auth/auth';
import { whichContentPageEnum } from '../pages/pages';

export const getNavbarButtonsBuilder = (authState, setAuthState) => {
    switch (authState.state) {
        case authStateEnum.notAuthed: case authStateEnum.loading: {//todo: loading state should be handled somehow by its own
            return (whichAuthPage, _) => {//unnamed to have a common interface
                return getNavbarButtonsWhenNotAuthed(whichAuthPage);//maybe, this needs refactoring; logout might be better of to be handled separately
            }
        }
        case authStateEnum.authed: {
            return (whichAuthPage, whichContentPage) => {
                return getNavbarButtonsWhenAuthed(setAuthState, whichAuthPage, whichContentPage);
            }
        }
        default:
            throw new Error("not handled; authState: ", authState)
    }
}

export const getNavbarButtonsWhenNotAuthed = (whichAuthPage) => {
    const allButtons = {//this way of structuring buttons might seems unnecarry as there are only two buttons to show when the user is authed. However, there can be more buttons in the future so I just took this approach.
        authButtons: [],
        contentButtons: [],
    };

    switch (whichAuthPage) {
        case whichAuthPageEnum.login:
            allButtons.authButtons = [...allButtons.authButtons, SignupButton];
            break;
        case whichAuthPageEnum.signup:
            allButtons.authButtons = [...allButtons.authButtons, LoginButton];
            break;
        default:
            throw new Error("not handled; whichAuthPage: ", whichAuthPage)
    }

    return allButtons;
}

export const getNavbarButtonsWhenAuthed = (setAuthState, whichAuthPage, whichContentPage) => {
    const allButtons = {
        authButtons: [],
        contentButtons: [],
    }

    switch (whichAuthPage) {
        case whichAuthPageEnum.login: case whichAuthPageEnum.signup: 
            allButtons.authButtons = [...allButtons.authButtons, () => <LogoutButton setAuthState = {setAuthState}/>];
            break;
        default:
            throw new Error(`not handled; whichAuthPage: ${whichAuthPage}` )
    }

    switch (whichContentPage) {
        case whichContentPageEnum.stats:
            allButtons.contentButtons = [...allButtons.contentButtons, WordGameButton];
            break;
        case whichContentPageEnum.wordGame:
            allButtons.contentButtons = [...allButtons.contentButtons, StatsButton];
            break;
        default:
            throw new Error("not handled; whichContentPage: ", whichContentPage)
    }

    return allButtons;
}

export const passStateSettersToNavbarButtons = (buttons, setWhichAuthPage, setWhichContentPage) => {
    const authButtonsPartial = buttons.authButtons.map((AuthButton) => {
        return key => <AuthButton 
            setWhichAuthPage={setWhichAuthPage} 
            key={key} 
            className="auth-nav-button"//todo this might be needed to seperate the buttons and to create borders in-between them
        />
    })
    const contentButtonsPartial = buttons.contentButtons.map((ContentButton) => {
        return key => <ContentButton 
            setWhichContentPage={setWhichContentPage} 
            key={key}
            className="content-nav-button"
        />
    })

    const allButtonsPartial = [...contentButtonsPartial, ...authButtonsPartial];
    const allButtons = allButtonsPartial.map((buttonPartial, index) => buttonPartial(index));

    return allButtons
}