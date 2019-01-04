


var conn;

if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {
        term.writeln("Connection closed.");
    };
    conn.onmessage = function (evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            term.writeln(messages[i]);
        }
        //term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ');
    };
} else {
}

var buf = [];
var term = new Terminal();
term.open(document.getElementById('terminal-container'));
term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ');

term.on("key", function (c, e) {
    if (e.key == "Enter") {
        var line = buf.join('');
        buf = [];
        conn.send(line);
        term.writeln('');
    } else {
        buf.push(c);
        term.write(c)
    }

    //console.log(c, e)
});