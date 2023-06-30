import React from 'react';

const authStateEnum = {
    shouldLogin: "shouldLogin",
    shouldSignup: "shouldSignup",
    authed: "authed",
}

const LoginButton = () => {
}

const SignupButton = () => {
}

const LogoutButton = () => {
}

const LoginForm = () => {
}

const SignupForm = () => {
}

const getAuthNavButtons = (authState) => {
    switch (authState) {
        case authStateEnum.authed:
            return [LogoutButton]
        case authStateEnum.shouldLogin: case authStateEnum.shouldSignup:
            return [LoginButton, SignupButton]
        default:
            throw new Error("not handled; auth state is not valid: ", authState)
    }
}

const getAuthContent = (authState) => {
    switch (authState) {
        case authStateEnum.authed:
            throw new Error("getContent from auth cannot be given authState = authed")
        case authStateEnum.shouldLogin:
            return LoginForm
        case authStateEnum.shouldSignup:
            return SignupForm
        default:
            throw new Error("not handled; auth state is not valid: ", authState)
    }
}

const whichPageEnum = {
    stats: "stats",
    wordGame: "wordGame"
}

const StatsButton = () => {
}

const WordGameButton = () => {
}

const StatsPage = () => {
}

const WordGamePage = () => {
}

const getPageNavButtons = () => {
    return [StatsButton, WordGameButton]
}

const getPageContent = (whichPage) => {
    switch (whichPage) {
        case whichPageEnum.stats:
            return StatsPage
        case whichPageEnum.wordGame:
            return WordGamePage
        default:
            throw new Error("invalid whichPage state: ", whichPage)
    }
}

const getNavButtons = (authState = authStateEnum.shouldLogin, pageState = whichPageEnum.stats) => {
    const allNavButtons = []
    const authNavButtons = getAuthNavButtons(authState)
    allNavButtons = [...authNavButtons, ...allNavButtons]
    if (authState.authed) {
        const pageNavButtons = getPageNavButtons(pageState)
        allNavButtons = [...pageNavButtons, ...allNavButtons]
    }
    //todo: need to add default nav buttons like a button for "about" page
    return allNavButtons
}

const getContent = (authState = "shouldLogin", pageState = "stats") => {
    if (!authState.authed) {
        return getAuthContent(authState)
    }
    return getPageContent(pageState)
}


function NavButtonsWrapper (authState, whichPage) {
    const NavButtons = (setAuthState, setWhichPage) => {
        return (
            <button onClick={setAuthState}></button>
        )
    }
    return NavButtons
}

const App = () => {
    const [authState, setAuthState] = useState(authStateEnum.shouldLogin);
    const [whichPage, setWhichPage] = useState(whichPageEnum.stats);
    const handleNotAuthorized = (response) => {
        if (response.status === 401) {
            setAuthState(authStateEnum.shouldLogin)
        }
    }
    const navButtons = getNavButtons(authState, setAuthState, whichPage, setWhichPage)
    const content = getContent(authState, setAuthState, whichPage, setWhichPage)

    return (
        <React.Fragment>
            {navButtons}
            {content}
        </React.Fragment>
    )
}















