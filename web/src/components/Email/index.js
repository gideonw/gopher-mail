import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";

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
        setEmail(result.data.email);
        console.log(result.data.email);
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
        <ul className="email">
          <li>id: {messageID}</li>
          <li>subject: {email.Subject}</li>
          <li>From: {email.From[0].Address}</li>
          <li>To: {email.To[0].Address}</li>
        </ul>
        <div className="email-body">{email.TextBody}</div>
      </div>
    );
  }
};

export default Email;
