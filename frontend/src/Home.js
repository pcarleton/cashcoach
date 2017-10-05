import React, { Component } from 'react';
import TransactionsTable from './TransactionsTable';
import Client from './Client';
import PlaidLink from 'react-plaid-link';


class AddAccount extends Component {
  render () {
        return <PlaidLink
          publicKey="5cf2c831a6e43805a92b01fa703ee8"
          product="transactions"
          env="sandbox"
          clientName="Cash Coach"
          onSuccess={this.handleOnSuccess}
          apiVersion="v2"
          buttonText="Add Account"
        />
  }

}

class Account extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    const props = this.props;
    return (<div>
        <p>{props.model.name}</p>
        <p>{props.model.bank}</p>
        <p>{props.model.balance}</p>
      </div>)
    
  }
}


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

  componentDidMount() {
    const cb = (data) => {
      console.log(data);
    }
    
    const accounts = Client.accounts(cb);
    Client.transactions(cb);
    //this.setState({accounts: accounts});
  }

  render() {
    const transactions = this.state.transactions;
    const model = {name: "Test Account", bank: "Bank of Test", balance: 200.0};

    return (
        <div>
        <Account model={model} />

      </div>
    )
  }
}


export default Home;
