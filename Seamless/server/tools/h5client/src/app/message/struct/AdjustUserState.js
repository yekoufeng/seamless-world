(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = AdjustUserState;

    var CMD = NetMsgMsgId.AdjustUserStateMsgID;
    function AdjustUserState() {
        this.Data = null;
    }

    var proto = AdjustUserState.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;
        this.Data = unpacker.binary.unpackerl(buf.slice(pos));
        pos += unpacker.binary.sizel(buf.slice(pos));
    };

})();