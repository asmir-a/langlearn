import React from 'react';

import { whichPageEnum } from '../../pages/pages';

const WordGameButton = ({ setWhichPage }) => {
    return (
        <button onClick={(_) => { setWhichPage(whichPageEnum.wordGame) }}></button>
    );
}

export default WordGameButton;