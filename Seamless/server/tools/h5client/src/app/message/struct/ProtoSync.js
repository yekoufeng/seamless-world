(function () {
    'use strict';

    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = ProtoSync;

    var CMD = NetMsgMsgId.ProtoSyncMsgID;
    function ProtoSync() {
        this.Data = null;
    }

    var proto = ProtoSync.prototype;

    proto.encode = function () {
        throw "don't need to do ";
    };

    proto.decode = function (buf) {
        var head = new NetMsgHead(0, 0);
        head.decode(buf);
        var jsonstr = buf.slice(NetMsgHead.len, head.msgSize()).toString('ascii');
        console.log(jsonstr);
        this.Data = JSON.parse(jsonstr);
    };

})();