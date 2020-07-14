import React from "react";
import { Redirect } from "react-router-dom";

import "./style.css";
import { authHandler } from "../../util/AuthHandler/AuthHandler";

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.login = this.login.bind(this);
    this.usernameInput = React.createRef();
    this.passwordInput = React.createRef();
  }

  state = {
    isLoaded: false,
    isLoading: false,
    loggedIn: false,
    redirectToReferrer: false,
  }

  login = (e) => {
    e.preventDefault();
    // const { username, password } = this.state;
    const username = this.usernameInput.current.value;
    const password = this.passwordInput.current.value;
    authHandler.authenticate(username, password, (token) => {
      console.log(token);
      this.setState(() => ({
        redirectToReferrer: true,
      }));
    });
  }

  render() {
    const { from } = this.props.location.state || { from: { pathname: "/" } };
    const { 
      redirectToReferrer, 
      isLoading,
      auth } = this.state;

    if (redirectToReferrer === true) {
      console.log("how");
      return <Redirect to={from} />;
    }

    return (
      <div className="login-box">
        <h2>Login</h2>
        <form>
          <div className="user-box">
            <input
              type="text"
              name=""
              required={true}
              disabled={isLoading}
              // onChange={this.handleUsername}
              // value={auth.username}
              ref={this.usernameInput}
            />
            <label>Username</label>
          </div>
          <div className="user-box">
            <input
              type="password"
              name=""
              required={true}
              disabled={isLoading}
              // onChange={this.handlePassword}
              // value={auth.password}
              ref={this.passwordInput}
            />
            <label>Password</label>
          </div>
          <button href="#" onClick={this.login}>
            <span></span>
            <span></span>
            <span></span>
            <span></span>
            Submit
          </button>
        </form>
      </div>
    );
  }
}

export default Login;
