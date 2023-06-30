export const httpCodes = {
    ok: 200,
    badRequest: 400,
    unauthorized: 401,
    forbidden: 403,
    notFound: 404,
}

export const endpoints = {
    user:                 '/api/auth/user',
    login:                '/api/auth/login',
    signup:               '/api/auth/signup',
    logout:               '/api/auth/logout',
    gameEntry:            '/api/wordgame/entry',
    getWords: username => `/api/users/${username}/wordgame/words`,//this might be a closure; we can close the getter under the username
    getStats: username => `/api/users/${username}/wordgame/stats`,
}