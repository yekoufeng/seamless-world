(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = RPCMsg;
    var CMD = NetMsgMsgId.RPCMsgID;
    function RPCMsg() {
        this.ServerType = 0;
        this.SrcEntityID = 0;
        this.MethodName = 0;
        this.Data = null;
    }

    var proto = RPCMsg.prototype;

    proto.encode = function () {
        var buf1 = unpacker.uint8.packerl(this.ServerType);
        var buf2 = unpacker.uint64.packerl(this.SrcEntityID);
        var buf3 = unpacker.string.packerl(this.MethodName);
        var buf4 = this.Data;
        return Buffer.concat([buf1, buf2, buf3, buf4]);
    };

    proto.encodeWithHead = function () {
        var buf1 = unpacker.uint8.packerl(this.ServerType);
        var buf2 = unpacker.uint64.packerl(this.SrcEntityID);
        var buf3 = unpacker.string.packerl(this.MethodName);
        var buf4 = unpacker.binary.packerl(this.Data);
        var len = buf1.length + buf2.length + buf3.length + buf4.length;
        return Buffer.concat([new NetMsgHead(len, CMD).encode(), buf1, buf2, buf3, buf4]);
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;

        this.ServerType = unpacker.uint8.unpackerl(buf.slice(pos));
        pos += unpacker.uint8.size();

        this.SrcEntityID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.MethodName = unpacker.string.unpackerl(buf.slice(pos));
        pos += unpacker.string.sizel(buf.slice(pos));

        this.Data = unpacker.binary.unpackerl(buf.slice(pos));
        pos += unpacker.binary.sizel(buf.slice(pos));
    };

})();