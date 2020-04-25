import React from 'react';
import { LoggedInMenu } from './menu/LoggedInMenu';
import { LoggedOutMenu } from './menu/LoggedOutMenu';
import { Round } from './round/Round';
import { Score } from './round/Score';
import { selectAuthenticated} from './auth/authSlice'
import { useSelector } from 'react-redux';
import {
  HashRouter,
  Switch,
  Route
} from "react-router-dom";

function App() {

  const authenticated = useSelector(selectAuthenticated);

  return (
      <HashRouter>
        <Routes authenticated={authenticated}/>
      </HashRouter>
  );
}

function Routes(props){
  if (props.authenticated){
    return LoggedInRoutes()
  } else {
    return LoggedOutRoutes()
  }
}

function LoggedOutRoutes() {
  return (
    <Switch>
      <Route path='/'children={<LoggedOutMenu/>}/>
    </Switch>
  )
}

function LoggedInRoutes() {
  return (
    <Switch>
      <Route path='/:id/s' children={<Score/>}/>
      <Route path='/:id' children={<Round/>}/>
      <Route path='/' children={<LoggedInMenu/>}/>
    </Switch>
  )
}

export default App;
