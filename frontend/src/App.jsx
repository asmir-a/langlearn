import React from 'react';
import { useState, useEffect } from 'react';
import WordsGame from './components/wordgame/ContentImagesGame';
import LoadingDisplay from './components/general/Loading';
import AuthErrorDisplay from './components/general/AuthError';
import NavBar from './components/navbar/NavBar';
import LoginForm from './components/auth/LoginForm';
import SignupForm from './components/auth/SignupForm';
import * as common from './utilites';
import "./App.css";

const IS_AUTHED_ENDPOINT = "/api/is-authed";
const requestToCheckIfAuthed = async () => {

  const response = await fetch(IS_AUTHED_ENDPOINT);
  if (response.status === common.HTTP_STATUS_OK) {//todo:handle other status codes
    return true;
  } else if (response.status === common.HTTP_STATUS_UNAUTHORIZED) {
    return false;
  } else {
    //todo: handle later
    throw new Error("not implemented");
  }
}

const App = () => {
  const [authState, setAuthState] = useState(common.AUTH_STATE_ENUM.Loading);
  
  useEffect(() => {
    const wrapperAroundAuthLogic = async () => {
      const isAuthedResponse = await requestToCheckIfAuthed();

      if (isAuthedResponse) {
        setAuthState(common.AUTH_STATE_ENUM.Authed);
        return;
      } else {
        setAuthState(common.AUTH_STATE_ENUM.ShouldLogin)//might be login, which depends on some further logic like sessionstorage, but for now it is okay
        return;
      }
    }
    wrapperAroundAuthLogic()
  }, []);

  const selectComponent = () => {
    switch (authState) {
      case common.AUTH_STATE_ENUM.Loading:
        return <LoadingDisplay />//todo: can this be put into a closure somehow

      case common.AUTH_STATE_ENUM.Authed:
        return <WordsGame setAuthState = {setAuthState}/>

      case common.AUTH_STATE_ENUM.ShouldSignup:
        return <SignupForm setAuthState = {setAuthState}/>

      case common.AUTH_STATE_ENUM.ShouldLogin:
        return <LoginForm setAuthState = {setAuthState}/>

      default:
        return <AuthErrorDisplay setAuthState = {setAuthState}/>
    }
  }

  return (
    <div id="app">
      <NavBar authState={authState} setAuthState={setAuthState}/>
      {selectComponent()}
    </div>
  );
}

export default App;