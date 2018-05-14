(function () {
    'use strict';

    var Util = require('../util/util.js');
    var NetMsgHead = require('../message/netmsg_head.js');
    var NetMsgMsgId = require('../message/netmsg_msgid.js');
    var ClientVertifyReq = require('../message/struct/ClientVertifyReq.js');
    var ClientVertifySucceedRet = require('../message/struct/ClientVertifySucceedRet.js');
    var SpaceUserConnect = require('../message/struct/SpaceUserConnect.js');
    var Page = require('../pages/page.js');
    var PageRoom = require('../pages/room.controller.js');
    var def = require('../message/def.js');
    var RPCMsg = require('../message/struct/RPCMsg.js');
    var unserializer = require('../message/rpc/unserialize.js');
    var AOISyncUserState = require('../message/struct/AOISyncUserState.js');
    var AdjustUserState = require('../message/struct/AdjustUserState.js');
    var EntityAOIS = require('../message/struct/EntityAOIS.js');
    var Entity = require('./entity.js');
    var unpacker = require('../message/struct/unpacker/unpacker.js');
    var EventInfo_EnterAOI = require('./eventinfo_enteraoi.js');
    var EventInfo_LeaveAOI = require('./eventinfo_leaveaoi.js');
    var V2D = require('./v2d.js');
    var SyncUserState = require('../message/struct/SyncUserState.js');

    module.exports = Room;

    function Room(user) {
        this.user = user;
        this.ws = null;                     // WebSocket对象
        this.recvbuf = null;                // 粘包处理用
        this.cmds = null;
        this.rpcdonothings = {};
        this.rpcs = null;
        this.timer = null;
        this.timer2 = null;

        // logic
        this.totalnum = 0;
        this.alivenum = 0;
        this.killnum = 0;
        this.airlift = -1;
        this.entitys = {};

        this.v2d = null;
        this.downFlagX = 0;
        this.downFlagZ = 0;
        this.timeStamp = 0;
        this.epoch = 0;
    }

    var proto = Room.prototype;

    proto.Login = function () {
        var self = this;
        self.ws = Util.initWebSocket(self.user.roomIP, self.user.roomPort,
            self.onopen.bind(self),
            self.onmessage.bind(self),
            self.onclose.bind(self)
        );
        if (this.timer == null) {
            this.timer = setInterval(this.Draw.bind(this), 25);
        }
        if (this.timer2 == null) {
            this.timer2 = setInterval(this.Move.bind(this), 34);
        }
    };

    proto.SyncState = function (x, z, t, e) {
        var buf1 = unpacker.uint32.packerl(t);
        var buf2 = unpacker.uint32.packerl(Entity.Mask_Pos_X | Entity.Mask_Pos_Z);
        var buf3 = unpacker.float32.packerl(x);
        var buf4 = unpacker.float32.packerl(z);
        var buf5 = unpacker.uint8.packerl(e);
        var len = buf1.length + buf2.length + buf3.length + buf4.length + buf5.length;
        var data = Buffer.concat([buf1, buf2, buf3, buf4, buf5]);

        var req = new SyncUserState();
        req.EntityID = this.user.entityID;
        req.Data = data;
        var reqbuf = req.encode();
        this.ws.send(reqbuf);
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
        console.log('[room] onclose');
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
            this.cmds = {};
            this.cmds[NetMsgMsgId.ClientVertifySucceedRetMsgID] = this.onmessage_ClientVertifySucceedRet.bind(this);
            this.cmds[NetMsgMsgId.ClientVertifyFailedRetMsgID] = this.onmessage_ClientVertifyFailedRet.bind(this);
            this.cmds[NetMsgMsgId.SpaceUserConnectSucceedRetMsgID] = this.onmessage_SpaceUserConnectSucceedRet.bind(this);
            this.cmds[NetMsgMsgId.AOISyncUserStateMsgID] = this.onmessage_AOISyncUserState.bind(this);
            this.cmds[NetMsgMsgId.AdjustUserStateMsgID] = this.onmessage_AdjustUserState.bind(this);
            this.cmds[NetMsgMsgId.EntityAOISMsgID] = this.onmessage_EntityAOIS.bind(this);
            this.cmds[NetMsgMsgId.RPCMsgID] = this.onmessage_RPCMsg.bind(this);
        }

        if (head.cmd in this.cmds) {
            this.cmds[head.cmd](buf);
        } else {
            console.log('[room] recv message, cmd:', head.cmd, ', size:', head.size, ', flag:', head.flag);
        }
    };


    proto.onmessage_ClientVertifySucceedRet = function (buf) {
        var msg = new ClientVertifySucceedRet();
        msg.decode(buf);
        console.log('[room] ClientVertifySucceedRet.Source =', msg.Source);
        console.log('[room] ClientVertifySucceedRet.UID =', msg.UID);
        console.log('[room] ClientVertifySucceedRet.SourceID =', msg.SourceID);
        console.log('[room] ClientVertifySucceedRet.Type =', msg.Type);
        console.log('[room] login to room success!');

        var req = new SpaceUserConnect();
        req.UID = this.user.uid;
        req.SpaceID = this.user.spaceID;
        var reqdata = req.encode();
        this.ws.send(reqdata);
    };

    proto.onmessage_ClientVertifyFailedRet = function (buf) {
        alert('[room] 登录Room失败！');
    };

    proto.onmessage_SpaceUserConnectSucceedRet = function (buf) {
        this.user.RPC(def.ServerTypeServer, "ParachuteReady");
        Page.showPage('room');
    };

    proto.onmessage_AOISyncUserState = function (buf) {
        var msg = new AOISyncUserState();
        msg.decode(buf);
        for (var i = 0; i < msg.Num; i++) {
            var entityID = msg.EIDS[i];
            var entity = null;
            if (entityID in this.entitys) {
                entity = this.entitys[entityID];
            } else {
                entity = new Entity();
                this.entitys[entityID] = entity;
            }
            if (entity != null) {
                entity.entityID = entityID;
                entity.UpdateData(msg.EDS[i], false);
            }
        }
    };

    proto.onmessage_AdjustUserState = function (buf) {
        var msg = new AdjustUserState();
        msg.decode(buf);
        var entity = new Entity();
        entity.entityID = this.user.entityID;
        entity.UpdateData(msg.Data, true);
        this.entitys[this.user.entityID] = entity;
        this.timeStamp = entity.timeStamp;
        this.epoch = entity.Epoch;
    };

    proto.onmessage_EntityAOIS = function (buf) {
        var msg = new EntityAOIS();
        msg.decode(buf);
        for (var i = 0; i < msg.Num; i++) {
            var flag = unpacker.uint8.unpackerl(msg.data[i]);
            if (flag == 1) {
                // enter aoi event
                var e1 = new EventInfo_EnterAOI();
                e1.UpdateData(msg.data[i].slice(1));
            } else if (flag == 0) {
                // leave aoi event
                var e2 = new EventInfo_LeaveAOI();
                e2.UpdateData(msg.data[i].slice(1));
                if (e2.entityID in this.entitys) {
                    delete this.entitys[e2.entityID];
                }
            } else {
                throw 'message data error!!';
            }
        }
    };

    proto.onmessage_RPCMsg = function (buf) {
        var msg = new RPCMsg();
        msg.decode(buf);

        if (msg.MethodName in this.rpcdonothings) {
            return;
        }

        if (this.RecvRPCMsg(msg)) {
            return;
        }

        console.log('[room] RPCMsg.MethodName =', msg.MethodName);
        console.log('[room] RPCMsg.Data =', msg.Data);
    };

    proto.RecvRPCMsg = function (msg) {
        if (this.rpcs == null) {
            this.rpcs = {};
            this.rpcs.UpdateAirLeft = this.RPC_UpdateAirLeft.bind(this);
            this.rpcs.UpdateTotalNum = this.RPC_UpdateTotalNum.bind(this);
            this.rpcs.UpdateKillNum = this.RPC_UpdateKillNum.bind(this);
            this.rpcs.UpdateAliveNum = this.RPC_UpdateAliveNum.bind(this);
        }
        if (msg.MethodName in this.rpcs) {
            this.rpcs[msg.MethodName](msg);
            return true;
        } else {
            return false;
        }
    };

    proto.RPC_UpdateAirLeft = function (msg) {
        this.airlift = unserializer.uint32(msg.Data);
    };
    proto.RPC_UpdateTotalNum = function (msg) {
        this.totalnum = unserializer.uint32(msg.Data);
    };
    proto.RPC_UpdateKillNum = function (msg) {
        this.killnum = unserializer.uint32(msg.Data);
    };
    proto.RPC_UpdateAliveNum = function (msg) {
        this.alivenum = unserializer.uint32(msg.Data);
    };


    proto.Draw = function () {
        var self = this;
        PageRoom.scope.uid = self.user.uid;
        PageRoom.scope.totalnum = self.totalnum;
        PageRoom.scope.alivenum = self.alivenum;
        PageRoom.scope.killnum = self.killnum;
        if (self.airlift > 0) {
            PageRoom.scope.airlift = "跳伞倒计时：" + String(self.airlift);
        } else if (self.airlift < 0) {
            PageRoom.scope.airlift = "准备跳伞中...";
        } else {
            PageRoom.scope.airlift = "";
        }

        if (self.v2d == null) {
            self.v2d = new V2D(self.user, self);
            self.v2d.initCanvas();
        }
        self.v2d.showBG();

        PageRoom.scope.$apply();

        self.Move();
    };

    proto.Move = function () {
        var self = this;
        if (self.downFlagX == 0 && self.downFlagZ == 0) {
            return;
        }
        var selfEntity = self.entitys[self.user.entityID];
        if (selfEntity == null) {
            return;
        }
        var x = selfEntity.posX = selfEntity.posX + 0.05 * self.downFlagX;
        var z = selfEntity.posZ = selfEntity.posZ + 0.05 * self.downFlagZ;
        self.SyncState(x, z, this.timeStamp, this.epoch);


        this.timeStamp++;
    };

})();