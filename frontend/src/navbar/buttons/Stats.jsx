import React from 'react';

import { whichContentPageEnum } from '../../pages/pages';

const StatsButton = ({ setWhichContentPage }) => {
    return (
        <button onClick={(_) => {setWhichContentPage(whichContentPageEnum.stats)}}>Stats</button>
    )
}

export default StatsButton;