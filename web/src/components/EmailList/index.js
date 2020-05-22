import React from "react";
import { Link, useRouteMatch } from "react-router-dom";
import "./style.css";

function EmailList() {
  let { url } = useRouteMatch();
  return (
    <div id="email-list">
      <ul className="email-list">
        <Link to={`${url}/rendering`}>
          <li>one</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>asd</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>oasdasdasne</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>oasdasdne</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>oaadadssdasdne</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>onasdasde</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>oaadadne</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>onadadadae</li>
        </Link>
        <Link to={`${url}/rendering`}>
          <li>ocacacne</li>
        </Link>
      </ul>
    </div>
  );
}

export default EmailList;
