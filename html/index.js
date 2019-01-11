
window.term = (function () {
    var term = new Terminal();
    var term_dom = document.getElementById('terminal');
    term.open(term_dom);

    term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ');

    var buf = [];
    term.on("key", function (c, e) {
        if (e.key == "Enter") {
            var line = buf.join('') + "\n";
            buf = [];
            conn.send(line);
            term.writeln('');
        } else {
            buf.push(c);
            term.write(c)
        }
    });
    return term;
})();

window.conn = (function () {
    var conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {
        term.writeln("Connection closed.");
    };
    conn.onmessage = function (evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            term.writeln(messages[i]);
        }
    };
    return conn;
})();
