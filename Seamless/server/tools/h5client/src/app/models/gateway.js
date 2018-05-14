(function () {
    'use strict';

    var Util = require('../util/util.js');
    var NetMsgHead = require('../message/netmsg_head.js');
    var NetMsgMsgId = require('../message/netmsg_msgid.js');
    var ClientVertifyReq = require('../message/struct/ClientVertifyReq.js');
    var ClientVertifySucceedRet = require('../message/struct/ClientVertifySucceedRet.js');
    var ProtoSync = require('../message/struct/ProtoSync.js');
    var MRolePropsSyncClient = require('../message/struct/MRolePropsSyncClient.js');
    var RPCMsg = require('../message/struct/RPCMsg.js');
    var def = require('../message/def.js');
    var EnterSpace = require('../message/struct/EnterSpace.js');

    module.exports = Gateway;

    function Gateway(user) {
        this.user = user;
        this.ws = null;                     // WebSocket对象
        this.recvbuf = null;                // 粘包处理用
        this.ProtoSync = null;              // TimeFire，server层协议信息(协议ID - 协议名)，登录后会发送。
        this.cmds = null;
        this.rpcdonothings = {};
    }

    var proto = Gateway.prototype;

    proto.Login = function () {
        var self = this;
        self.ws = Util.initWebSocket(self.user.gatewayIP, self.user.gatewayPort,
            self.onopen.bind(self),
            self.onmessage.bind(self),
            self.onclose.bind(self)
        );
    };

    proto.onopen = function () {
        this.ws.binaryType = 'arraybuffer';
        var req = new ClientVertifyReq();
        req.Source = 0;
        req.UID = this.user.uid;
        req.Token = this.user.token;
        var buf = req.encode();
        this.ws.send(buf);
    };

    proto.onmessage = function (data) {
        var buf = Buffer.from(data);
        if (!!this.recvbuf) {
            buf = Buffer.concat([this.recvbuf, buf]);
            this.recvbuf = null;
        }
        var pos = 0;
        while (pos < buf.length) {
            if (buf.length - pos < NetMsgHead.len) {
                break;
            }
            var curBuf = buf.slice(pos);
            var head = new NetMsgHead(0, 0);
            head.decode(curBuf);
            if (buf.length - pos < head.msgSize()) {
                break;
            }
            this.onmessage_switch(head, curBuf);
            pos += head.msgSize();
        }
        if (pos < buf.length) {
            this.recvbuf = buf.slice(pos);
        }
    };

    proto.onclose = function () {
        console.log('[gateway] onclose');
    };

    proto.onmessage_switch = function (head, buf) {
        if (head.cmd == NetMsgMsgId.HeartBeatMsgID) {
            // heart beat
            return;
        }

        if (head.flag == 2) {
            alert('不支持加密，请更改服务器配置！');
            throw 'msg unsupport error!';
        }

        if (this.cmds == null) {
            this.rpcdonothings.JumpAir = 0;
            this.rpcdonothings.SyncFriendList = 0;
            this.rpcdonothings.SyncApplyList = 0;
            this.rpcdonothings.InitOwnGoodsInfo = 0;
            this.rpcdonothings.InitNotifyMysqlDbAddr = 0;
            this.rpcdonothings.OnlineCheckMatchOpen = 0;

            this.cmds = {};
            this.cmds[NetMsgMsgId.ClientVertifySucceedRetMsgID] = this.onmessage_ClientVertifySucceedRet.bind(this);
            this.cmds[NetMsgMsgId.ClientVertifyFailedRetMsgID] = this.onmessage_ClientVertifyFailedRet.bind(this);
            this.cmds[NetMsgMsgId.ProtoSyncMsgID] = this.onmessage_ProtoSync.bind(this);
            this.cmds[NetMsgMsgId.MRolePropsSyncClientMsgID] = this.onmessage_MRolePropsSyncClient.bind(this);
            this.cmds[NetMsgMsgId.RPCMsgID] = this.onmessage_RPCMsg.bind(this);
            this.cmds[NetMsgMsgId.EnterSpaceMsgID] = this.onmessage_EnterSpace.bind(this);
            this.cmds[NetMsgMsgId.AOISyncUserStateMsgID] = this.user.room.onmessage_AOISyncUserState.bind(this.user.room);
            this.cmds[NetMsgMsgId.AdjustUserStateMsgID] = this.user.room.onmessage_AdjustUserState.bind(this.user.room);
            this.cmds[NetMsgMsgId.EntityAOISMsgID] = this.user.room.onmessage_EntityAOIS.bind(this.user.room);
        }

        if (head.cmd in this.cmds) {
            this.cmds[head.cmd](buf);
        } else {
            console.log('[gateway] recv message, cmd:', head.cmd, ', size:', head.size, ', flag:', head.flag);
        }
    };

    proto.onmessage_ClientVertifySucceedRet = function (buf) {
        var msg = new ClientVertifySucceedRet();
        msg.decode(buf);
        console.log('[gateway] ClientVertifySucceedRet.Source =', msg.Source);
        console.log('[gateway] ClientVertifySucceedRet.UID =', msg.UID);
        console.log('[gateway] ClientVertifySucceedRet.SourceID =', msg.SourceID);
        console.log('[gateway] ClientVertifySucceedRet.Type =', msg.Type);
        console.log('[gateway] login to gateway success!');
    };

    proto.onmessage_ClientVertifyFailedRet = function (buf) {
        alert('[gateway] 登录Gateway失败！');
    };

    proto.onmessage_MRolePropsSyncClient = function (buf) {
        var msg = new MRolePropsSyncClient();
        msg.decode(buf);
        console.log('[gateway] MRolePropsSyncClient.EntityID =', msg.EntityID);
        this.user.entityID = msg.EntityID;

        // 开始登录Lobby
        this.user.lobby.Login();
    };

    proto.onmessage_ProtoSync = function (buf) {
        this.ProtoSync = new ProtoSync();
        this.ProtoSync.decode(buf);
        console.log('[gateway] ProtoSync.Data:', this.ProtoSync.Data);
    };


    proto.onmessage_RPCMsg = function (buf) {
        var msg = new RPCMsg();
        msg.decode(buf);

        if (msg.MethodName in this.rpcdonothings) {
            return;
        }

        if (this.user.lobby.RecvRPCMsg(msg)) {
            return;
        }

        console.log('[gateway] RPCMsg.MethodName =', msg.MethodName);
        console.log('[gateway] RPCMsg.Data =', msg.Data);
    };

    proto.onmessage_EnterSpace = function (buf) {
        var msg = new EnterSpace();
        msg.decode(buf);

        console.log('[gateway] EnterSpace.SpaceID =', msg.SpaceID);
        console.log('[gateway] EnterSpace.MapName =', msg.MapName);
        console.log('[gateway] EnterSpace.EntityID =', msg.EntityID);
        console.log('[gateway] EnterSpace.Addr =', msg.Addr);
        console.log('[gateway] EnterSpace.TimeStamp =', msg.TimeStamp);

        var addrinfo = msg.Addr.split(":");
        this.user.spaceID = msg.SpaceID;
        this.user.roomIP = addrinfo[0];
        this.user.roomPort = parseInt(addrinfo[1]);
        this.user.room.Login();
    };

})();