import React from 'react';
import axios from 'axios';
import {useEffect, useState} from 'react';

import LearnedWordsColumn from './LearnedWordsList';
import LearningWordsColumn from './LearningWordsList';

import { endpoints } from '../../utilities';

import "./Stats.css";

const fetchWords = async (wordsUrl) => {
    const response = await axios.get(wordsUrl);
    return response.data
}

const StatsPage = ({user}) => {//todo: maybe these should be wrapped in a Content Container for easier styling
    const [words, setWords] = useState(null);

    useEffect(() => {
        const fetchWordsWrapper = async () => {
            const fetchedWords = await fetchWords(endpoints.getWordsEndpoint(user.username));//maybe, it would be better if the component did not know anything about the endpoints
            setWords(fetchedWords);
        }
        fetchWordsWrapper();
    }, []);

    return (
        <main className="stats">
            {
                words ?
                    <article>
                        <h1>your lists of words</h1>
                        <LearningWordsColumn learningWords={words.learning} />
                        <LearnedWordsColumn learnedWords={words.learned} />
                    </article> :
                    <p>loading...</p>
            }
        </main>
    )
}

export default StatsPage;