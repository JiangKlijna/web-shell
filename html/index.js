// xterm.js
window.term = (function () {
    fit.apply(Terminal);
    webLinks.apply(Terminal);
    var dom = document.createElement("div");
    dom.id = "terminal";
    document.body.appendChild(dom);
    var term = new Terminal();

    var OnWebLinkClick = function (e, url) {
        open(url);
    };

    var buf = [];
    term.on("key", function (c, e) {
        console.log(e)
        if (e.which === 13) {
            var line = buf.join('') + "\n";
            buf = [];
            console.log(line);
            conn.send(line);
            // term.writeln('');
        }
    });

    term.on('data', function (data) {
        buf.push(data);
        term.write(data);
    });


    term.on('paste', function (data) {
        term.write(data);
        // this.copy = term.getSelection();
    });

    // term.on("selection", function() {
    //     if (term.hasSelection()) {
    //         this.copy = term.getSelection();
    //     }
    // });

    term.open(dom);
    term.fit();
    term.webLinksInit(OnWebLinkClick);
    return term;
})();

function blobToString(blob, encoding, fun) {
    var reader = new FileReader();
    reader.onloadend = function () {
        fun(reader.result);
    };
    reader.readAsText(blob, encoding);
}

// websocket client
window.conn = (function () {
    var conn = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + "/cmd");
    conn.onclose = function (e) {
        term.writeln("connection closed.");
    };
    var onWrite = function (message) {
        term.write(message);
    };
    conn.onmessage = function (e) {
        blobToString(e.data, "gbk", onWrite);
    };
    return conn;
})();
