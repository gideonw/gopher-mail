import React from "react";
import { Link, useRouteMatch } from "react-router-dom";
import _ from "lodash";

import "./style.css";

const EmailList = (props) => {
  let { url } = useRouteMatch();

  let emailList = [];
  _.mapKeys(props.emails, (email, messageID) => {
    emailList.push(
      <Link key={messageID} to={`${url}/${messageID}`}>
        <li>
          {email.Date} - {email.Subject}
        </li>
      </Link>
    );
  });

  return (
    <div id="email-list">
      <ul className="email-list">{emailList}</ul>
    </div>
  );
};

export default EmailList;
