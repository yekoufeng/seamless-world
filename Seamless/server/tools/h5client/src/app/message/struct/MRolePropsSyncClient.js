(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = MRolePropsSyncClient;

    var CMD = NetMsgMsgId.MRolePropsSyncClientMsgID;
    function MRolePropsSyncClient() {
        this.EntityID = 0;
    }

    var proto = MRolePropsSyncClient.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var pos = NetMsgHead.len;

        this.EntityID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();
    };

})();