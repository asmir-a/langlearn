export const httpCodes = {
    ok: 200,
    badRequest: 400,
    unauthorized: 401,
    forbidden: 403,
    notFound: 404,
    conflict: 409,
}

export const endpoints = {
    user: '/api/auth/user',
    login: '/api/auth/login',
    signup: '/api/auth/signup',
    logout: '/api/auth/logout',
    getRandomGameEntry: username => `/api/wordgame/entries/users/${username}/random`,
    getSubmitAnswer: username => `/api/wordgame/entries/users/${username}/submit`,
    getWordsEndpoint: username => `/api/wordgame/users/${username}/words`,//this might be a closure; we can close the getter under the username
    getWordCountsEndpoint: username => `/api/wordgame/users/${username}/word-counts`,
}