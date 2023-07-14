import React, {useRef} from 'react';

import { 
    colors, 
    submitAnswerFor, 
    submitAnswerAndToggleResultIndicator, 
    getResultIndicatorToggler 
} from './logic';

import "./Button.css";

const IncorrectWordButton = ({
    user, 
    word, 
    updateGame,
}) => {
    const toggleResultIndicator = useRef(getResultIndicatorToggler(colors.red));

    const submitAnswer= async () => {
        await submitAnswerFor(user.username, word, true);
    }

    const handleClick = async (_) => {
        await submitAnswerAndToggleResultIndicator(
            submitAnswer,
            toggleResultIndicator.current,
        )
        updateGame()
    }

    return (
        <button
            onClick={handleClick}
            className = "word-button"
        >{word}</button>
    )
}

export default IncorrectWordButton;