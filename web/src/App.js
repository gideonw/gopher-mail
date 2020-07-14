import React from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";

import PrivateRoute from "./util/PrivateRoute/PrivateRoute";
import Header from "./components/Header";
import Login from "./routes/Login";
import Mailbox from "./routes/Mailbox";

import "./App.css";

function App() {
  return (
    <Router>
      <div className="App flex flex-col flex-1">
        <Header />
        {/* <PrivateRoute exact path="/" component={Login}>
        </PrivateRoute> */}
        <PrivateRoute path="/" component={Mailbox} />
        <Route exact path="/login" component={Login} />
      </div>
    </Router>
  );
}

export default App;
