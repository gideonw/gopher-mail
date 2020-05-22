import React from "react";
import { Switch, Route, useRouteMatch } from "react-router-dom";

import "./style.css";

import MailboxList from "../../components/MailboxList";
import EmailList from "../../components/EmailList";
import Email from "../../components/Email";

function Mailbox() {
  let { path } = useRouteMatch();
  return (
    <div className="mailbox">
      <MailboxList />
      <Switch>
        <Route exact path={`${path}`}>
          <EmailList />
        </Route>
        <Route path={`${path}/:messageID`}>
          <Email />
        </Route>
      </Switch>
    </div>
  );
}

export default Mailbox;
