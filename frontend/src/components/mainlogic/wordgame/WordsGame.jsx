import React from 'react';
import {useEffect, useState} from 'react';
import * as common from './../../../utilites';

import "./WordsGame.css"

const WORDS_GAME_ENDPOINT = "/api/wordgame/entries/random";
const HANDLE_ANSWER_ENDPOINT = "/api/wordgame/entries/submit"

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

const IncorrectWordButton = ({username, word, reloadData}) => {//answer should be validated on the server side for now. when the times for optimizations comes, this might make more sense
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


const StatsShower = ({stats}) => {
    return (
        <ul className = "stats-shower">
            <li>Words Learning: {stats.learningCount}</li>
            <li>Words Learned: {stats.learnedCount}</li>
        </ul>
    )
}

const GameEntry = ({username, gameEntry, reloadData, stats}) => {
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
            <StatsShower stats={stats}/>
            <img src={gameEntry.correctWordImageUrl}/>
            <div className="button-group">
                {allButtons}
            </div>
        </article>
    );
}

const makeRequestToStats = async (username) => {
    const wordCountsUrl = `/api/users/${username}/wordgame/stats`
    const response = await fetch(wordCountsUrl)

    if (response.status === common.HTTP_STATUS_OK) {
        return await response.json()
    } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
        return null;
    } else {
        throw new Error("not implemented: request to stats")
    }
}

const makeRequestToGameEntry = async () => {
    const response = await fetch(WORDS_GAME_ENDPOINT)
    if (response.status === common.HTTP_STATUS_OK) {
        return await response.json();
    } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
        return null;
    } else {
        throw new Error("not implemented: request to gameentry");
    }
}


export default function WordsGame({username, setAuthInfo}) {//todo: this component might be useless; maybe, just using the component above would suffice
    const [currentGameEntryData, setCurrentGameEntryData] = useState(null);
    const [stats, setStats] = useState(null);

    const fetchData = async () => {
        const gameEntry = await makeRequestToGameEntry();
        const stats = await makeRequestToStats(username);


        if (gameEntry !== null) {
            setCurrentGameEntryData(gameEntry)
        } else {
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.ShouldLogin,
                user: null,
            })
        }

        if (stats !== null) {
            setStats(stats)
        } else {
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.ShouldLogin,
                user: null,
            })
        }
    }

    useEffect(() => {
        fetchData()
    }, []);

    return (
        <>
          {currentGameEntryData && 
            <GameEntry 
                username = {username} 
                gameEntry = {currentGameEntryData} 
                reloadData = {fetchData}
                stats = {stats}
            />
          }
        </>
    );
}