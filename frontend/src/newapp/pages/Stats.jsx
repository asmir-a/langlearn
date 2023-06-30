import React from 'react';

import { whichPageEnum } from './pages';

const StatsButton = ({setWhichPage}) => {
    return (
        <button onClick={(_) => setWhichPage(whichPageEnum.stats)}></button>
    );
}

export default StatsButton;