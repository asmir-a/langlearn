import React from 'react';

import Quiz from './Quiz';
import WordCountsShower from './WordCountsShower';

import "./WordGame.css";

const WordGamePage = ({user}) => {
    return (
        <main className = "word-game">
            <WordCountsShower user={user} />
            <Quiz
                user={user}
            />
        </main>
    )
}

export default WordGamePage;