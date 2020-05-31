import React from "react";
import { Redirect } from "react-router-dom";
import _ from "lodash";
import axios from "axios";

import "./style.css";

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      auth: {
        username: "",
        password: "",
      },
      isLoaded: false,
      isLoading: false,
      loggedIn: false,
    };

    this.handleLogin = this.handleLogin.bind(this);
    this.handleUsername = this.handleUsername.bind(this);
    this.handlePassword = this.handlePassword.bind(this);
  }

  handleUsername(e) {
    this.setState(_.merge(this.state, { auth: { username: e.target.value } }));
  }
  handlePassword(e) {
    this.setState(_.merge(this.state, { auth: { password: e.target.value } }));
  }

  componentDidMount() {
    let token = localStorage.getItem("gopher-mail-auth-token");
    console.log("Logged in" + token);
    if (token) {
      this.setState(_.merge(this.state, { loggedIn: true }));
    }
  }

  handleLogin(e) {
    console.log(e);
    e.preventDefault();
    this.setState(_.merge(this.state, { isLoading: true }));
    axios
      .post("https://gps.gideonw.xyz/api/auth/login", this.state.auth)
      .then((result) => {
        console.log(result);
        if (result.status === 200) {
          localStorage.setItem("gopher-mail-auth-token", result.data.token);
          this.setState(_.merge(this.state, { loggedIn: true }));
        } else {
          // TODO: error states
          this.setState(
            _.merge(this.state, { isLoading: false, isLoaded: false })
          );
        }
      })
      .catch((error) => {
        // handle error
        console.log(error);
        this.setState(
          _.merge(this.state, { isLoading: false, isLoaded: false })
        );
      });
  }

  render() {
    if (this.state.loggedIn === true) {
      return <Redirect to="/m" />;
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
              disabled={this.state.isLoading}
              onChange={this.handleUsername}
              value={this.state.auth.username}
            />
            <label>Username</label>
          </div>
          <div className="user-box">
            <input
              type="password"
              name=""
              required={true}
              disabled={this.state.isLoading}
              onChange={this.handlePassword}
              value={this.state.auth.password}
            />
            <label>Password</label>
          </div>
          <button href="#" onClick={this.handleLogin}>
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

export default Home;
