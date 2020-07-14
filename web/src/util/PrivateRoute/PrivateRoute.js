import React from "react";
import { Route, Redirect } from "react-router";
import { authHandler } from "../AuthHandler/AuthHandler";

const PrivateRoute = ({ component: Component, ...rest }) => {
    return (
        <Route {...rest} render={(props) => (
            authHandler.isLoggedIn() === true
              ? <Component {...props} />
              : <Redirect to={{
                  pathname: '/login',
                  from: props.location
                }} />
          )} /> 
    );
};

export default PrivateRoute;