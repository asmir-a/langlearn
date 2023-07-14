import React from 'react';

const WordRow = ({word}) => {
    return (
        <li>{word}</li>
    )
}

const LearnedWordsColumn = ({learnedWords}) => {
    return (
        <section>
            <h2>words you already learned</h2>
            {
                learnedWords.length > 0 ?
                    <ul>
                        {learnedWords.map((word, index) => {
                            return <WordRow word={word} key={index} />
                        })}
                    </ul>
                    : <p>you have not learn any words yet</p>
            }
        </section>
    )
}

export default LearnedWordsColumn;