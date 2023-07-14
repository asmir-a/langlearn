import React from 'react';
import axios from 'axios';

import {useState, useEffect} from 'react';
import { endpoints } from '../../utilities';

import "./WordCountsShower.css";

const fetchWordCounts = async (username) => {
    const response = await axios.get(endpoints.getWordCountsEndpoint(username));
    return response.data;
}

const WordCountsShower = ({user}) => {
    const [wordCounts, setWordCounts] = useState(null);

    useEffect(() => {
        const fetchWordCountsWrapper = async () => {
            const fetchedWordCounts = await fetchWordCounts(user.username);
            setWordCounts(fetchedWordCounts);
        }
        fetchWordCountsWrapper();
    }, []);

    return (
        wordCounts ?
            <article className="word-counts-shower">
                <h1>current progress</h1>
                <output>words learned: {wordCounts.learned}</output>
                <output>words learning: {wordCounts.learning}</output>
            </article>
        : <p>loading...</p>
    );
}

export default WordCountsShower;