
export const whichPageEnum = {
    stats: "stats",
    wordGame: "wordGame",
}

export const selectPageComponent = (whichPage) => {
    switch (whichPage) {
        case whichPageEnum.stats:
        case whichPageEnum.wordGame:
        default:
            throw new Error("not handled; whichPage: ", whichPage)
    }
}