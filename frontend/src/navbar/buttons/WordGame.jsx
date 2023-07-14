import React from 'react';

import { whichContentPageEnum } from '../../pages/pages';

const WordGameButton = ({ setWhichContentPage }) => {
    return (
        <button onClick={(_) => { setWhichContentPage(whichContentPageEnum.wordGame) }}>WordGame</button>
    );
}

export default WordGameButton;