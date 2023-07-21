import React from 'react';

import {useState, useEffect} from 'react';
import {shuffle} from 'lodash';

import IncorrectWordButton from './IncorrectWordButton';
import CorrectWordButton from './CorrectWordButtons';

import { fetchGameEntryFor } from './logic';

import "./Quiz.css";

const Quiz = ({user}) => {
    const [gameEntry, setGameEntry] = useState(null);
    const [atLeastOneImageLoaded, setAtLeastOneImageLoaded] = useState(false);//cannot count images because corb warning cannot be catched; this should work fine bacause the images are stored in s3 so if one loads; others are gonna be loaded quite soon after that

    const updateGame = () => {
        const clearGame = () => {
            setGameEntry(null);
            setAtLeastOneImageLoaded(false);
        }
        const fetchAndUpdateGameEntry = async () => {
            const newGameEntry = await fetchGameEntryFor(user.username);
            setGameEntry(newGameEntry);
        }

        clearGame();
        fetchAndUpdateGameEntry();
    }

    useEffect(() => {
        updateGame();
    }, []);

    const handleImageLoad = async (_) => {
        setAtLeastOneImageLoaded(true);
    }

    return (
        gameEntry ?
            <article className = "quiz">
                <h1>select the right word matching the image</h1>
                <figure>
                    {gameEntry.correctWordImageUrls.map((imageUrl, index) => {
                        return <img
                            src={imageUrl}
                            key={index}
                            onLoad={handleImageLoad}
                            style={
                                atLeastOneImageLoaded ? {} : { display: "none" }
                            }
                        />
                    })}
                </figure>
                <section>
                    {
                        shuffle([
                            ...gameEntry.incorrectWords,
                            gameEntry.correctWord
                        ].map((word, index) => {
                            return word === gameEntry.correctWord ?
                                <CorrectWordButton
                                    user={user}
                                    word={word}
                                    updateGame={updateGame}
                                    key={index}
                                />
                                : <IncorrectWordButton
                                    user={user}
                                    word={word}
                                    updateGame={updateGame}
                                    key={index}
                                />
                        }))
                    }
                </section>
            </article>
            : <p>loading...</p>
    )
}

export default Quiz;