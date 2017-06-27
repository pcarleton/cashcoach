import React, { Component } from 'react';
import Home from './Home';
import GoogleLogin from 'react-google-login';
import Client from './Client';

class App extends Component {

  constructor(props) {
    super(props);
    this.state = {loggedIn: false, error: ""};

    this.responseGoogle = this.responseGoogle.bind(this);
    this.login = this.login.bind(this);
  }

  login(data) {

    console.log("Authed!")
    console.log(data);
    // TODO: Display some profile info
    this.setState({loggedIn: true});
  }

  responseGoogle(response) {
    console.log(response);
    Client.verify(this.login, response.tokenId)
  }

  loginFromCookie() {
      try {
        Client.me(this.login);
        console.log("Logged in from cookie.")
      }
      catch (e) {
        if (e.code != 401) {
          throw e;
        }
      }
  }

  componentDidMount() {
    console.log("mount method");
    this.loginFromCookie();
  }

  render() {
    const loggedIn = this.state.loggedIn;

    if (loggedIn) {
      return (
      <div className="App">
          <Home />
          </div>
        )
    }
    return (
      <div className='App'>
        <div className='ui text container'>
        <GoogleLogin
           clientId="217209923893-kv53i3hqgk1plrapk0eub1p4jr7sipet.apps.googleusercontent.com"
           buttonText="Sign in with Google"
           onSuccess={this.responseGoogle}
           onFailure={this.responseGoogle}
         />
        </div>
      </div>
    );
  }
}

export default App;
