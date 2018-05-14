(function () {
    'use strict';

    var unpacker = require('./unpacker/unpacker.js');
    var NetMsgHead = require('../netmsg_head.js');
    var NetMsgMsgId = require('../netmsg_msgid.js');

    module.exports = ClientVertifyReq;

    var CMD = NetMsgMsgId.ClientVertifyReqMsgID;
    function ClientVertifyReq() {
        this.Source = 0;
        this.UID = 0;
        this.Token = "";
    }

    var proto = ClientVertifyReq.prototype;

    proto.encode = function () {
        var buf1 = unpacker.uint8.packerl(this.Source);
        var buf2 = unpacker.uint64.packerl(this.UID);
        var buf3 = unpacker.string.packerl(this.Token);
        var len = buf1.length + buf2.length + buf3.length;
        return Buffer.concat([new NetMsgHead(len, CMD).encode(), buf1, buf2, buf3]);
    };

    proto.decode = function (buf) {
        throw "don't need to do ";
    };

})();