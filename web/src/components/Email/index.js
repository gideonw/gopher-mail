import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";

import { IconContext } from "react-icons";
import { BsReply, BsReplyAll, BsForward } from "react-icons/bs";

import "./style.css";

const Email = () => {
  let { messageID } = useParams();

  const [email, setEmail] = useState();
  const [isLoaded, setIsLoaded] = useState(false);
  const [loading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isLoaded || loading) return;
    setIsLoading(true);
    axios
      .get(`https://gps.gideonw.xyz/api/gideon/email/${messageID}`)
      .then((result) => {
        setEmail(result.data.Email);
        setIsLoaded(true);
        setIsLoading(false);
      });
  }, [messageID, email, isLoaded, loading]);

  if (!isLoaded) {
    return (
      <div id="email">
        <ul className="email">
          <li>id: {messageID}</li>
          <li>Loading</li>
        </ul>
        <div className="email-body">loading...</div>
      </div>
    );
  } else {
    return (
      <div id="email">
        <div className="email-actions">
          <IconContext.Provider
            value={{
              size: "1.25em",
              style: { marginRight: "4px" },
            }}
          >
            <div>
              <BsReply />
              Reply
            </div>
            <div>
              <BsReplyAll />
              Reply All
            </div>
            <div>
              <BsForward />
              Forward
            </div>
          </IconContext.Provider>
        </div>
        <ul className="email">
          <li>id: {messageID}</li>
          <li>subject: {email.Subject}</li>
          <li>From: {email.From[0].Address}</li>
          <li>To: {email.To[0].Address}</li>
        </ul>
        <div className="email-body">
          <pre>{email.TextBody}</pre>
        </div>
      </div>
    );
  }
};

export default Email;
