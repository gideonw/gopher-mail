import React, { useState, useEffect } from "react";
import { Switch, Route, useRouteMatch } from "react-router-dom";
import _ from "lodash";

import "./style.css";

import MailboxList from "../../components/MailboxList";
import EmailList from "../../components/EmailList";
import Email from "../../components/Email";
import axios from "axios";

function Mailbox() {
  const [emails, setEmails] = useState({});
  const [isLoaded, setIsLoaded] = useState(false);
  const [loading, setIsLoading] = useState(false);

  let { path } = useRouteMatch();

  useEffect(() => {
    if (isLoaded || loading) return;
    setIsLoading(true);
    axios.get("https://gps.gideonw.xyz/api/gideon/emails").then((result) => {
      console.log(result.data);

      _.forEach(result.data.emails, (value) => {
        emails[value] = {};
        setEmails(emails);
      });
      setIsLoaded(true);
      setIsLoading(false);
    });
  }, [emails, isLoaded, loading]);

  let emailList = <EmailList emails={emails} />;
  if (!isLoaded) {
    emailList = <div>loading</div>;
  }

  return (
    <div className="mailbox">
      <MailboxList />
      <Switch>
        <Route exact path={`${path}`}>
          {emailList}
        </Route>
        <Route path={`${path}/:messageID`}>
          <Email />
        </Route>
      </Switch>
    </div>
  );
}

export default Mailbox;
