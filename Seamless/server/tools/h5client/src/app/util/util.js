(function () {
    'use strict';

    require('sprintf-js');

    module.exports = Util;

    function Util() { }

    Util.toUpper = function (s) {
        return s.substring(0, 1).toUpperCase() + s.substring(1);
    };

    Util.initWebSocket = function (tip, tport, onopen, onmessage, onclose) {
        if ("WebSocket" in window) {
            var config = require('../../../webpack.config.js');
            var ws = new WebSocket(sprintf('ws://%s:%d', config.devServer.addr, config.devServer.port + 1));
            ws.onopen = function () {
                ws.send(sprintf('{"tip":"%s", "tport":%d}', tip, tport));
                if (!!onopen) {
                    onopen();
                }
            };
            ws.onmessage = function (evt) {
                var data = evt.data;
                if (!!onmessage) {
                    onmessage(data);
                }
            };
            ws.onclose = function () {
                if (!!onclose) {
                    onclose();
                }
            };
            return ws;
        }
        else {
            alert("您的浏览器不支持WebSocket!");
            return null;
        }
    };

})();