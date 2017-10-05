
function post(cb, endpoint, jsonData) {
  return fetch(endpoint, {
    method: 'post',
    accept: 'application/json',
    headers: {
      "Content-type": "application/json"
    },
    credentials: "same-origin",
    body: JSON.stringify(jsonData)
  }).then(checkStatus)
    .then(parseJSON)
    .then(cb);
}

function get(cb, endpoint) {
  return fetch(endpoint, {
    accept: 'application/json',
    credentials: "same-origin"
  }).then(checkStatus)
    .then(parseJSON)
    .then(cb);

}


function transactions(cb) {
  return get(cb, 'api/transactions');
}

function accounts(cb) {
  return get(cb, 'api/accounts');
}


function me(cb) {
  return get(cb, 'api/me');
}

function verify(cb, idToken) {
  return post(cb, 'api/jwt', {'idtoken': idToken}); 
}

function addAccount(cb, account_nick, public_token) {
  return post(cb, 'api/accounts/add', {'name': account_nick,
                                       'public_token': public_token});
}



function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  const error = new Error(`HTTP Error ${response.statusText}`);
  error.status = response.statusText;
  error.code = response.status;
  error.response = response;
  throw error;
}

function parseJSON(response) {
  return response.json();
}

const Client = { transactions, verify, me, addAccount, accounts};
export default Client;
