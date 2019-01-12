
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

window.fit = (function () {
    function proposeGeometry(term) {
        if (!term.element.parentElement) {
            return null;
        }
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
        var geometry = proposeGeometry(term);
        if (geometry) {
            if (term.rows !== geometry.rows || term.cols !== geometry.cols) {
                term._core.renderer.clear();
                term.resize(geometry.cols, geometry.rows);
            }
        }
    }
})();