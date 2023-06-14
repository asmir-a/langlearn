import React from 'react';
import {useEffect, useState} from 'react';
import * as common from '../../utilites';

const WORDS_GAME_ENDPOINT = "/api/word-game/game-entries/random";

const GameEntryExample = {
    imageUrl: "",
    correctWord: "",
    incorrectWords: ["", "", ""],
}

const CorrectWordButton = ({word, reloadData}) => {
    const handleClick = () => {
        console.log("WIN: correct word");
    }
    return <button onClick = {handleClick}>{word}</button>
}

const IncorrectWordButton = ({word, reloadData}) => {
    const handleClick = () => {
        console.log("LOSE: incorrect word")
    }
    return <button onClick = {handleClick}>{word}</button>
}

const GameEntry = ({gameEntry, reloadData}) => {
    const correctWordButton = <CorrectWordButton //if we want the user not be able to cheat (i do not want for what purpose, it might be better not to indicate the correct word and let the server decide); but it feels like it would load the server for needing to keep the word session/context or game in the database
        word={gameEntry.correctWord} reloadData={reloadData} key = {gameEntry.correctWord}
    />
    const incorrectWordButtons = gameEntry.incorrectWords.map(
        incorrectWord => <IncorrectWordButton word = {incorrectWord} reloadData={reloadData} key = {incorrectWord}/> 
    )

    const randomIndex = Math.floor(Math.random() * 3);
    const allButtons = [...incorrectWordButtons.slice(0, randomIndex), correctWordButton, ...incorrectWordButtons.slice(randomIndex)];

    return (
        <article className="game-entry">
            <img src={gameEntry.correctWordImageUrl}/>
            <div className="button-group">
                {allButtons}
            </div>
        </article>
    );
}

export default function WordsGame({setAuthState}) {//todo: this component might be useless; maybe, just using the component above would suffice
    const [currentGameEntryData, setCurrentGameEntryData] = useState(null);

    const fetchGameEntry = async () => {
        const response = await fetch(WORDS_GAME_ENDPOINT);
        if (response.status === common.HTTP_STATUS_OK) {
            const gameEntryFromServer = await response.json();
            setCurrentGameEntryData(gameEntryFromServer);
            return;
        } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
            setAuthState(common.AUTH_STATE_ENUM.ShouldLogin)
            return;
        } else {
            throw new Error("not implemented");
        }
    }

    useEffect(() => {
        fetchGameEntry();
    }, []);

    return (
        <>
          {currentGameEntryData && <GameEntry gameEntry={currentGameEntryData} reloadData = {fetchGameEntry}/>}
        </>
    );
}