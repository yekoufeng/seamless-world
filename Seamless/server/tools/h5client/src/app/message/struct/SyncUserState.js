(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = SyncUserState;

    var CMD = NetMsgMsgId.SyncUserStateMsgID;
    function SyncUserState() {
        this.EntityID = 0;
        this.Data = null;
    }

    var proto = SyncUserState.prototype;

    proto.encode = function () {
        var buf1 = unpacker.uint64.packerl(this.EntityID);
        var buf2 = unpacker.binary.packerl(this.Data);
        var len = buf1.length + buf2.length;
        return Buffer.concat([new NetMsgHead(len, CMD).encode(), buf1, buf2]);
    };

    proto.decode = function (buf) {
        throw "don't need to do ";
    };

})();