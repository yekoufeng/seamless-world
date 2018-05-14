(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = ClientVertifySucceedRet;

    var CMD = NetMsgMsgId.ClientVertifySucceedRetMsgID;
    function ClientVertifySucceedRet() {
        this.Source = 0;
        this.UID = 0;
        this.SourceID = 0;
        this.Type = 0;
    }

    var proto = ClientVertifySucceedRet.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;

        this.Source = unpacker.uint8.unpackerl(buf.slice(pos));
        pos += unpacker.uint8.size();

        this.UID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.SourceID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.Type = unpacker.uint8.unpackerl(buf.slice(pos));
        pos += unpacker.uint8.size();
    };

})();