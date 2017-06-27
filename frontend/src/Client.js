function transactions(cb) {
  return fetch(`api/transactions`, {
    accept: 'application/json',
    credentials: "same-origin"
  }).then(checkStatus)
    .then(parseJSON)
    .then(cb);
}

function me(cb) {
  return fetch(`api/me`, {
    accept: 'application/json',
    credentials: "same-origin"
  }).then(checkStatus)
    .then(parseJSON)
    .then(cb);
}

function verify(cb, idToken) {
  return fetch(`api/jwt`,
    {
        method: 'post',
        headers: {
          "Content-type": "application/json"
        },
        body: JSON.stringify({"idtoken": idToken}),
        accept: 'application/json',
        credentials: "same-origin"
      }
).then(checkStatus)
 .then(parseJSON)
 .then(cb);
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

const Client = { transactions, verify, me };
export default Client;
