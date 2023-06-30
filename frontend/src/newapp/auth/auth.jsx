import React from 'react';

export const authStateEnum = {
    authed: "authed",
    shouldLogin: "shouldLogin",
    shouldSignup: "shouldSignup",
}

export const selectAuthComponent = (authState) => {
    switch (authState) {
        case authStateEnum.shouldLogin:
        case authStateEnum.shouldSignup:
        case authStateEnum.authed:
            throw new Error("authState cannot be authed")
        default:
            throw new Error("not handled; authState: ", authState)
    }
}
