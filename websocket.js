console.log('Hello world!');

var connection = new WebSocket('ws://localhost:61489/tweets', ['soap', 'xmpp']);

connection.onopen = function() {
  console.log('Connection open to tweets');
};

connection.onerror = function(error) {
  console.log('WebSocket Error ' + error);
};

connection.onmessage = function (e) {
  console.log('Server: ' + e.data);
};


