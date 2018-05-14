(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = AOISyncUserState;

    var CMD = NetMsgMsgId.AOISyncUserStateMsgID;
    function AOISyncUserState() {
        this.Num = 0;
        this.EIDS = [];
        this.EDS = [];
    }

    var proto = AOISyncUserState.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;
        this.Num = unpacker.uint32.unpackerl(buf.slice(pos));
        pos += unpacker.uint32.size();
        for (var i = 0; i < this.Num; i++) {
            var EID = unpacker.uint64.unpackerl(buf.slice(pos));
            pos += unpacker.uint64.size();
            var ED = unpacker.binary.unpackerl(buf.slice(pos));
            pos += unpacker.binary.sizel(buf.slice(pos));
            this.EIDS.push(EID);
            this.EDS.push(ED);
        }
    };

})();