/**
 * @author jiangklijna
 */

// init
(function (W) {
    W.DataType = {
        Err: 0,
        Data: 1,
        Resize: 2,
    };
    W.Terminal = Terminal.Terminal;
    W.FitAddon = FitAddon.FitAddon;
    W.WebLinksAddon = WebLinksAddon.WebLinksAddon;

    W.NewWebSocket = function () {
        return new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + "/cmd");
    };
    W.NewTerminal = function () {
        return new Terminal({useStyle: true, screenKeys: true});
    };
})(window);

// web-shell component
window.WebShell = function (dom) {
    var term = NewTerminal();
    var conn = NewWebSocket();
    var fitAddon = new FitAddon();
    var webLinksAddon = new WebLinksAddon();

    // websocket connect
    conn.onclose = function (e) {
        term.writeln("connection closed.");
    };

    conn.onopen = function () {
        fitAddon.fit();
    };

    conn.onmessage = function (msg) {
        term.write(msg.data);
    };

    var send = function (dataType, data) {
        conn.send(JSON.stringify({'t': dataType, "d": data}));
    };

    // terminal term
    term.onTitleChange(function (title) {
        document.title = title;
    });

    // term.on('paste', function (data) {
    //     term.write(data);
    //     // this.copy = term.getSelection();
    // });

    term.onResize(function (data) {
        send(DataType.Resize, [data.cols, data.rows]);

    });
    term.onData(function (data) {
        send(DataType.Data, data);
    });
    term.open(dom);


    term.loadAddon(fitAddon);
    term.loadAddon(webLinksAddon);

    this.fit = function () {
        fitAddon.fit();
    };
    this.term = term;
    this.conn = conn;

    dom.oncontextmenu = function(){
        if (term.hasSelection()) {
            send(DataType.Data, term.getSelection());
            // term.write(term.getSelection());
            return false;
        }
        return true;
    }
};

// run
(function (W) {
    var dom = document.createElement("div");
    dom.style.width = "100%";
    dom.style.height = "100%";
    document.body.appendChild(dom);
    W.singleWebShell = new WebShell(dom);
    W.onresize = function () {
        W.singleWebShell.fit();
    }
})(window);
