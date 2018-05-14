(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = EntityAOIS;

    var CMD = NetMsgMsgId.EntityAOISMsgID;
    function EntityAOIS() {
        this.Num = 0;
        this.data = [];
    }

    var proto = EntityAOIS.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;
        this.Num = unpacker.uint32.unpackerl(buf.slice(pos));
        pos += unpacker.uint32.size();
        for (var i = 0; i < this.Num; i++) {
            var d = unpacker.binary.unpackerl(buf.slice(pos));
            pos += unpacker.binary.sizel(buf.slice(pos));
            this.data.push(d);
        }
    };

})();