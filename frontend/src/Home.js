import React, { Component } from 'react';
import TransactionsTable from './TransactionsTable';
import Client from './Client';
import PlaidLink from 'react-plaid-link';



class Home extends Component {

  constructor(props) {
    super(props);

    this.state = {transactions: [], accounts: []};
    this.handleOnSuccess = this.handleOnSuccess.bind(this);
  }

  fetchTransactions() {
    const updateTs = (tdata) => {
        console.log(tdata);
        this.setState({transactions: tdata.transactions})
    }
    Client.transactions(updateTs);
  }

  handleOnSuccess(resp) {
    console.log("success!");
    console.log(resp);
    const cb = (data) => {
      console.log("resp!");
      console.log(data);
    }
    Client.addAccount(cb, 'new account1', resp);
  }




  render() {
    const transactions = this.state.transactions;

    return (
        <div>
        <PlaidLink
          publicKey="5cf2c831a6e43805a92b01fa703ee8"
          product="transactions"
          env="sandbox"
          clientName="Cash Coach"
          onSuccess={this.handleOnSuccess}
          apiVersion="v2"
        />

      <TransactionsTable transactions={transactions}/>
      </div>
    )
  }
}


export default Home;
