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
  }




  render() {
    const transactions = this.state.transactions;

    return (
        <PlaidLink
          publicKey="5cf2c831a6e43805a92b01fa703ee8"
          product="transactions"
          env="sandbox"
          clientName="react-client"
          onSuccess={this.handleOnSuccess}
        />

      <TransactionsTable transactions={transactions}/>
    )
  }
}


export default Home;
