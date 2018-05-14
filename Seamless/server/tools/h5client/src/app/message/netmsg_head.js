(function () {
    'use strict';

    var unpacker = require('./struct/unpacker/unpacker.js');

    module.exports = NetMsgHead;

    function NetMsgHead(size, cmd) {
        this.size = size + 2;
        this.flag = 0;
        this.cmd = cmd;
    }

    NetMsgHead.len = 6;

    var proto = NetMsgHead.prototype;

    proto.encode = function () {
        var buf1 = unpacker.uint24.packerl(this.size);
        var buf2 = unpacker.uint8.packerl(this.flag);
        var buf3 = unpacker.uint16.packerl(this.cmd);
        var len = buf1.length + buf2.length + buf3.length;
        return Buffer.concat([buf1, buf2, buf3]);

    };

    proto.decode = function (buf) {
        var pos = 0;

        this.size = unpacker.uint24.unpackerl(buf.slice(pos));
        pos += unpacker.uint24.size();

        this.flag = unpacker.uint8.unpackerl(buf.slice(pos));
        pos += unpacker.uint8.size();

        this.cmd = unpacker.uint16.unpackerl(buf.slice(pos));
        pos += unpacker.uint16.size();
    };

    proto.msgSize = function () {
        return this.size + 4;
    };

})();