import React from "react";
import { Link } from "react-router-dom";

import "./style.css";

const MailboxList = () => (
  <div id="mailbox-list">
    <ul className="mailboxes">
      <Link to="/m">
        <li>Inbox</li>
      </Link>
      <li>Sent</li>
      <li>Spam</li>
    </ul>
  </div>
);

export default MailboxList;
