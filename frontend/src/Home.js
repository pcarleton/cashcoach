import React, { Component } from 'react';
import TransactionsTable from './TransactionsTable';
import Client from './Client';



class Home extends Component {

  constructor(props) {
    super(props);

    this.state = {transactions: [], accounts: []};
  }

  fetchTransactions() {
    const updateTs = (tdata) => {
        console.log(tdata);
        this.setState({transactions: tdata.transactions})
    }
    Client.transactions(updateTs);
  }

  render() {
    const transactions = this.state.transactions;

    return (
      <TransactionsTable transactions={transactions}/>
    )
  }
}


export default Home;
