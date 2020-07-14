import axios from 'axios';

class AuthHandler {
    getTokenCookie(){

    }
    
    setTokenCookie(token){

    }

    getTokenLS(){
        return true;
    }

    setTokenLS(token){

    }

    getLoggedInLS(){
        return true;
    }

    setLoggedInLS(){

    }
    
    authenticate(username, password, callback) {
        axios.post(
            "https://gps.gideonw.xyz/api/auth/login", 
            {username: username, password: password})
        .then((d) => {
            console.log(d);
            
            callback();
        });
    }

    isLoggedIn() {
        return false;
    }
};

export let authHandler = new AuthHandler();