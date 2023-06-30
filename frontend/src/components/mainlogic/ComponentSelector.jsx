import StatsPage from "./stats/Stats"
import WordsGame from "./wordgame/WordsGame";

export const pagesEnum = {
    wordgame: "wordgame",
    stats: "stats",
    none: "none"
}

export const componentSelector = (
    whichPageToShow,
    username,
    setAuthInfo,
) => {
    switch (whichPageToShow.page) {
        case pagesEnum.wordgame:
            return <WordsGame
                username={username}
                setAuthInfo={setAuthInfo}
            />
        case pagesEnum.stats:
            return <StatsPage
                username={username}
                setAuthInfo={setAuthInfo}
            />
        default:
            console.log("whichPageToShow: ", whichPageToShow)
            // throw new Error("not handled; component selector")

    }
}
