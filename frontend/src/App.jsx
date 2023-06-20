import React from 'react';
import { useState, useEffect } from 'react';
import WordsGame from './components/wordgame/WordsGame';
import LoadingDisplay from './components/general/Loading';
import AuthErrorDisplay from './components/general/AuthError';
import NavBar from './components/navbar/NavBar';
import LoginForm from './components/auth/LoginForm';
import SignupForm from './components/auth/SignupForm';
import * as common from './utilites';
import "./App.css";

const IS_AUTHED_ENDPOINT = "/api/is-authed";

const getCurrentUserInfo = async () => {
  const response = await fetch(IS_AUTHED_ENDPOINT);
  if (response.status === common.HTTP_STATUS_OK) {//todo:handle other status codes
    const user = await response.json();
    return user;
  } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
    return null;
  } else {
    //todo: handle later
    throw new Error("not implemented");
  }
}

const App = () => {
  const [authInfo, setAuthInfo] = useState(
    {
      authState: common.AUTH_STATE_ENUM.Loading,
      user: null,
    }
  );
  
  useEffect(() => {
    const wrapperAroundAuthLogic = async () => {
      const userJson = await getCurrentUserInfo();
      if (userJson) {
        setAuthInfo({
          authState: common.AUTH_STATE_ENUM.Authed,
          user: userData
        });
        return;
      } else {
        setAuthInfo({
          authState: common.AUTH_STATE_ENUM.ShouldLogin,
          user: null
        });
        return;
      }
    }
    wrapperAroundAuthLogic()
  }, []);

  const selectComponent = () => {
    switch (authInfo.authState) {
      case common.AUTH_STATE_ENUM.Loading:
        return <LoadingDisplay />//todo: can this be put into a closure somehow
      case common.AUTH_STATE_ENUM.Authed:
        return <WordsGame username = {authInfo.user.username} setAuthInfo = {setAuthInfo}/>
      case common.AUTH_STATE_ENUM.ShouldSignup:
        return <SignupForm setAuthInfo = {setAuthInfo}/>
      case common.AUTH_STATE_ENUM.ShouldLogin:
        return <LoginForm setAuthInfo = {setAuthInfo}/>
      default:
        throw new Error("not implemented")
    }
  }

  return (
    <div id="app">
      <NavBar authInfo={authInfo} setAuthInfo={setAuthInfo}/>
      {selectComponent()}
    </div>
  );
}

export default App;