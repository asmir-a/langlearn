import React from 'react';

import { getNavbarButtons } from './nav';

const Navbar = ({ authState, setAuthState, whichPage, setWhichPage }) => {
    const buttons = getNavbarButtons(authState, whichPage);

    const authButtons = buttons.authButtons;
    const authButtonsComponents = authButtons.map((AuthButton, index) => {
        return <AuthButton
            setAuthState={setAuthState}
            key={`authButton#${index}`}
        />;
    });

    const pageButtons = buttons.pageButtons;
    const pageButtonsComponents = pageButtons.map((PageButton, index) => {
        return <PageButton
            setWhichPage={setWhichPage}
            key={`pageButton#${index}`}
        />;
    });

    const allButtonsComponents = [...authButtonsComponents, ...pageButtonsComponents];

    return (
        <nav>
            {allButtonsComponents}
        </nav>
    )
}

export default Navbar;