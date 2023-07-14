import React from 'react';

import { getNavbarButtonsBuilder, passStateSettersToNavbarButtons } from './logic';

import "./Navbar.css";

const Navbar = ({
    authState, //maybe, it would have been better if app was building the navbar by itself. or, there might be an approach that utilizes custom hooks.
    setAuthState,
    whichAuthPage,
    setWhichAuthPage,
    whichContentPage,
    setWhichContentPage
}) => {
    const buttonsBuilder = getNavbarButtonsBuilder(authState, setAuthState);
    const buttons = buttonsBuilder(whichAuthPage, whichContentPage);
    const buttonsWithStateSetters = passStateSettersToNavbarButtons(buttons, setWhichAuthPage, setWhichContentPage);

    return (
        <nav>
            {buttonsWithStateSetters}
        </nav>
    )
}

export default Navbar;