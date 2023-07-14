import StatsPage from "./stats/Stats"
import WordGamePage from "./wordgame/WordGame"

export const whichContentPageEnum = {
    stats: "stats",
    wordGame: "wordGame",
}

export const selectContentPage = (whichContentPage) => {
    switch (whichContentPage) {
        case whichContentPageEnum.stats:
            return StatsPage;
        case whichContentPageEnum.wordGame:
            return WordGamePage;
        default:
            throw new Error("not handled; whichPage: ", whichPage)
    }
}
