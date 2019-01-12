// xterm.js
window.term = (function () {
    var term_dom = document.getElementById('terminal');
    var term = new Terminal();
    term.open(term_dom);

    term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ');

    var buf = [];
    term.on("key", function (c, e) {
        buf.push(c);
        term.write(c);
        if (e.key === "Enter") {
            var line = buf.join('');
            buf = [];
            console.log(line);
            conn.send(line);
            term.writeln('');
        }
    });
    return term;
})();

// websocket client
window.conn = (function () {
    var conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (e) {
        term.writeln("connection closed.");
    };
    conn.onmessage = function (e) {
        term.writeln(e.data);
    };
    return conn;
})();

// xterm.js.addons.fit
window.fit = (function () {
    function proposeGeometry(term) {
        var parentElementStyle = window.getComputedStyle(term.element.parentElement);
        var parentElementHeight = parseInt(parentElementStyle.getPropertyValue('height'));
        var parentElementWidth = Math.max(0, parseInt(parentElementStyle.getPropertyValue('width')));
        var elementStyle = window.getComputedStyle(term.element);
        var elementPadding = {
            top: parseInt(elementStyle.getPropertyValue('padding-top')),
            bottom: parseInt(elementStyle.getPropertyValue('padding-bottom')),
            right: parseInt(elementStyle.getPropertyValue('padding-right')),
            left: parseInt(elementStyle.getPropertyValue('padding-left'))
        };
        var elementPaddingVer = elementPadding.top + elementPadding.bottom;
        var elementPaddingHor = elementPadding.right + elementPadding.left;
        var availableHeight = parentElementHeight - elementPaddingVer;
        var availableWidth = parentElementWidth - elementPaddingHor - term._core.viewport.scrollBarWidth;
        return {
            cols: Math.floor(availableWidth / term._core.renderer.dimensions.actualCellWidth),
            rows: Math.floor(availableHeight / term._core.renderer.dimensions.actualCellHeight)
        };
    }

    window.onresize = function (e) {
        if (!term.element.parentElement) {
            return null;
        }
        var geometry = proposeGeometry(term);
        if (term.rows !== geometry.rows || term.cols !== geometry.cols) {
            term._core.renderer.clear();
            term.resize(geometry.cols, geometry.rows);
        }
    }
})();