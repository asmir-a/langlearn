import axios from "axios";

import { endpoints } from "../../utilities";

const colors = {
    red: "var(--main-red-color)",
    green: "var(--main-green-color)",
}

const resultIndicationDuration = 800;

const fetchGameEntryFor = async (username) => {
    const response = await axios.get(endpoints.getRandomGameEntry(username));
    return response.data
}

const submitAnswerFor = async (username, word, isAnswerCorrect) => {
    await axios.post(endpoints.getSubmitAnswer(username), {
        word: word,
        isAnswerCorrect: isAnswerCorrect,
    });
}

const promiseToWait = (someFunc, duration) => new Promise((res) => {
    const timerId = setTimeout(() => {
        someFunc();
        res(timerId);
    }, duration);
})

const submitAnswerAndToggleResultIndicator = async (
    submitAnswer,
    toggleResultIndicator,
) => {
    toggleResultIndicator();
    const submitAnswerPromise = submitAnswer();
    const timerId = setTimeout(() => {//promise status cannot be answered so i cannot use promiseToWait; maybe there is a way though to make it uniform
        toggleResultIndicator();
        timerId = null;
    }, resultIndicationDuration)

    const startTime = Date.now();
    await submitAnswerPromise;
    const endTime = Date.now();

    if (timerId !== null) {
        const promiseDuration = endTime - startTime;
        const remainingDuration = Math.max(0, resultIndicationDuration - promiseDuration);
        clearTimeout(timerId);
        await promiseToWait(toggleResultIndicator, remainingDuration);
    }
}

const getResultIndicatorToggler = (color) => {//this is a backdoor to dom; this can be done without refs or document
    let resultIndicator = null;
    const toggle = () => {
        if (resultIndicator === null) {
            resultIndicator = document.createElement("div");
            const indicatorStyle = {
                zIndex: "100",
                opacity: "70%",
                backgroundColor: color,
                position: "fixed",
                minHeight: "100%",
                minWidth: "100%",
                inset: 0,
            }
            for (const styleProperty in indicatorStyle) {
                resultIndicator.style[styleProperty] = indicatorStyle[styleProperty];
            }
            document.body.append(resultIndicator);
        } else {
            resultIndicator.remove();
            resultIndicator = null;
        }
    }
    return toggle;
}

export {colors, fetchGameEntryFor, submitAnswerFor, getResultIndicatorToggler, submitAnswerAndToggleResultIndicator};