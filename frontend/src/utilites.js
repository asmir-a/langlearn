export const HTTP_STATUS_OK = 200;
export const HTTP_STATUS_UNAUTHORIZED = 401;
export const HTTP_STATUS_FORBIDDEN = 403;
export const HTTP_STATUS_CONFLICT = 409;
export const IS_AUTHED_ENDPOINT = "/api/is-authed";

export const AUTH_STATE_ENUM = {
  ShouldLogin: "shouldLogin",
  ShouldSignup: "shouldSignup",
  Authed: "authed",
  Loading: "loading",
  SomeError: "someError",//should also be under consideration, but for now it is okay
}

