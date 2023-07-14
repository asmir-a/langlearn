import React from 'react';

const WordRow = ({word}) => {
    return (
        <li>{word}</li>
    )
}

const LearningWordsColumn = ({learningWords}) => {
    return (
        <section>
            <h2>words you are still learning</h2>
            {
                learningWords.length > 0 ?
                    <ul>
                        {learningWords.map((word, index) => {
                            return <WordRow word={word} key={index} />
                        })}
                    </ul>
                    : <p>you have not started learning words yet</p>
            }
        </section>
    )
}

export default LearningWordsColumn;