(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = SpaceEntityMsg;
    var CMD = NetMsgMsgId.SpaceEntityMsgID;
    function SpaceEntityMsg() {
        this.UID = 0;
        this.Data = null;
    }

    var proto = SpaceEntityMsg.prototype;

    proto.encode = function () {
        var buf1 = unpacker.uint64.packerl(this.UID);
        var buf2 = this.Data;
        var len = buf1.length + buf2.length;
        return Buffer.concat([new NetMsgHead(len, CMD).encode(), buf1, buf2]);
    };

    proto.decode = function (buf) {
        throw "don't need to do ";
    };

})();