import React from 'react';
import Client from './Client';

const MATCHING_ITEM_LIMIT = 25;

class TransactionsTable extends React.Component {
  constructor(props) {
    super(props);
    this.setState({transactions: props.transactions});
  }

  state = {
    transactions: [],
    startDate: '',
    endDate: '',
  };

  handleSearchChange = (e) => {
    const value = e.target.value;

    this.setState({
      startDate: value,
    });

    if (value === '') {
      this.setState({
        transactions: [],
      });
    } else {
      Client.transactions((transactions) => {
        this.setState({
          transactions: transactions.slice(0, MATCHING_ITEM_LIMIT),
        });
      });
    }
  };

  render() {
    const { transactions } = this.state;
    const transRows = transactions.map((trans, idx) => (
      <tr
        key={idx}
      >
        <td className='right aligned'>{trans.date}</td>
        <td className='right aligned'>{trans.name}</td>
        <td className='right aligned'>{trans.amount}</td>

        <td className='right aligned'>{trans.category}</td>
      </tr>
    ));


    return (
      <div id='food-search'>
        <table className='ui selectable structured large table'>
          <thead>
          <tr>
              <th colSpan='4'>
              <div className='ui icon input'>
                    <input
                      className='prompt'
                      type='text'
                      placeholder='Start date'
                      value={this.state.startDate}
                      onChange={this.handleSearchChange}
                    />
                    <i className='search icon' />
                  </div>
              </th>
              </tr>
            <tr>
              <th className='eight wide'>Date</th>
              <th>Description</th>
              <th>Amount</th>
              <th>Category</th>
            </tr>
          </thead>
          <tbody>
            {transRows}
          </tbody>
        </table>
      </div>
    );
  }
}


export default TransactionsTable;
