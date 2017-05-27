import React, { Component } from 'react';
import TransactionsTable from './TransactionsTable';
import './App.css';
import GoogleLogin from 'react-google-login';
import Client from './Client';

class App extends Component {

  constructor(props) {
    super(props);

    this.state = {transactions: []};
  }

  responseGoogle = (response) => {

    const updateTs = (tdata) => {
        console.log(tdata);
        this.setState({transactions: tdata.transactions})
    }
    console.log(response);
    Client.verify(function(data) {
      console.log("Authed!")
      console.log(data);
      Client.transactions(updateTs);
    }, response.tokenId)
  }

  render() {
    const transactions = this.state.transactions;
    return (
      <div className='App'>
        <div className='ui text container'>
        <GoogleLogin
           clientId="217209923893-kv53i3hqgk1plrapk0eub1p4jr7sipet.apps.googleusercontent.com"
           buttonText="Sign in with Google"
           onSuccess={this.responseGoogle}
           onFailure={this.responseGoogle}
         />,
          <TransactionsTable transactions={transactions}/>
        </div>
      </div>
    );
  }
}

export default App;
