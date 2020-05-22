import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import Header from "./components/Header";
import Landing from "./routes/Landing";
import Mailbox from "./routes/Mailbox";

import "./App.css";

function App() {
  return (
    <Router>
      <div className="App">
        <Header />
        <Switch>
          <Route exact path="/">
            <Landing />
          </Route>
          <Route path="/m">
            <Mailbox />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
