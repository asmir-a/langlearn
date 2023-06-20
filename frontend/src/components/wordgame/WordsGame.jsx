import React from 'react';
import {useEffect, useState} from 'react';
import * as common from '../../utilites';

const WORDS_GAME_ENDPOINT = "/api/wordgame/entries/random";
const HANDLE_ANSWER_ENDPOINT = "/api/wordgame/entries/submit"

const GameEntryExample = {
    imageUrl: "",
    correctWord: "",
    incorrectWords: ["", "", ""],
}

const AnswerExample = {
    isAnswerCorrect: false,
    word: "",
    username: ""
}


const CorrectWordButton = ({username, word, reloadData}) => {
    const handleClick = async () => {
        const response = await fetch(HANDLE_ANSWER_ENDPOINT, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({username: username, isAnswerCorrect: true, word: word})
        });
        if (response.status === common.HTTP_STATUS_OK) {
            reloadData();
            return;
        } else {
            throw new Error("not implemented");//we need access to the setAuthState prolly
        }
    }
    return <button onClick = {handleClick}>{word}</button>
}

const IncorrectWordButton = ({username, word, reloadData}) => {
    const handleClick = async () => {
        const response = await fetch(HANDLE_ANSWER_ENDPOINT, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                isAnswerCorrect: false,
                word: word
            })
        })
        if (response.status === common.HTTP_STATUS_OK) {
            reloadData();
            return;
        } else {
            throw new Error("not implemented");
        }
    }
    return <button onClick = {handleClick}>{word}</button>
}

const GameEntry = ({username, gameEntry, reloadData}) => {
    const correctWordButton = <CorrectWordButton //if we want the user not be able to cheat (i do not want for what purpose, it might be better not to indicate the correct word and let the server decide); but it feels like it would load the server for needing to keep the word session/context or game in the database
        username={username} word={gameEntry.correctWord} reloadData={reloadData} key = {gameEntry.correctWord}
    />
    const incorrectWordButtons = gameEntry.incorrectWords.map(
        incorrectWord => <IncorrectWordButton username={username} word = {incorrectWord} reloadData={reloadData} key = {incorrectWord}/> 
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

export default function WordsGame({username, setAuthInfo}) {//todo: this component might be useless; maybe, just using the component above would suffice
    const [currentGameEntryData, setCurrentGameEntryData] = useState(null);

    const fetchGameEntry = async () => {
        const response = await fetch(WORDS_GAME_ENDPOINT);
        if (response.status === common.HTTP_STATUS_OK) {
            const gameEntryFromServer = await response.json();
            setCurrentGameEntryData(gameEntryFromServer);
            return;
        } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.ShouldLogin,
                user: null
            });
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
          {currentGameEntryData && <GameEntry username = {username} gameEntry={currentGameEntryData} reloadData = {fetchGameEntry}/>}
        </>
    );
}