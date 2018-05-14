(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = EnterSpace;

    var CMD = NetMsgMsgId.EnterSpaceMsgID;
    function EnterSpace() {
        this.SpaceID = 0;
        this.MapName = "";
        this.EntityID = 0;
        this.Addr = "";
        this.TimeStamp = 0;
    }

    var proto = EnterSpace.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;

        this.SpaceID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.MapName = unpacker.string.unpackerl(buf.slice(pos));
        pos += unpacker.string.sizel(buf.slice(pos));

        this.EntityID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.Addr = unpacker.string.unpackerl(buf.slice(pos));
        pos += unpacker.string.sizel(buf.slice(pos));

        this.TimeStamp = unpacker.uint32.unpackerl(buf.slice(pos));
        pos += unpacker.uint32.size();
    };

})();