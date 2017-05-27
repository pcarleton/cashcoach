function transactions(cb) {
  return fetch(`transactions`, {
    accept: 'application/json',
    credentials: "same-origin"
  }).then(checkStatus)
    .then(parseJSON)
    .then(cb);
}


function verify(cb, idToken) {
  return fetch(`jwt`,
    {
        method: 'post',
        headers: {
          "Content-type": "application/json"
        },
        body: JSON.stringify({"idtoken": idToken}),
        accept: 'application/json'
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
  error.response = response;
  console.log(error); // eslint-disable-line no-console
  throw error;
}

function parseJSON(response) {
  return response.json();
}

const Client = { transactions, verify };
export default Client;
