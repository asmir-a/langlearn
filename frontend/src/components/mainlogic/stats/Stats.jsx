import React from 'react';
import {useEffect, useState} from 'react';
import * as common from './../../../utilites';


const makeRequestToWords = async (username) => {
    const wordsEndpoint = `/api/users/${username}/wordgame/words`
    const response = await fetch(wordsEndpoint)
    if (response.status === common.HTTP_STATUS_OK) {
        return response.json()
    } else if (response.status === common.HTTP_STATUS_FORBIDDEN) {
        return nil
    } else {
        throw new Error("not handled: makeRequestToWords")
    }
}

const StatsPage = ({username, setAuthInfo}) => {
    const [wordLists, setWordLists] = useState(nil);
    useEffect(() => {
        const wordListsFromResponse = makeRequestToWords(username);
        if (wordListsFromResponse === null) {
            setAuthInfo({
                authState: common.AUTH_STATE_ENUM.ShouldLogin,
                user: null,
            });
        } else {
            setWordLists(wordListsFromResponse);
        }
    }, [])

    return (
        <div>{JSON.stringify(wordLists)}</div>
    )
}

export default StatsPage;

