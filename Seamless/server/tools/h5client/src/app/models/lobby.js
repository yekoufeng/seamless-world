(function () {
    'use strict';

    var def = require('../message/def.js');
    var serializer = require('../message/rpc/serialize.js');
    var unserializer = require('../message/rpc/unserialize.js');
    var _ = require('../message/proto/game_pb.js');
    var PageLobby = require('../pages/lobby.controller.js');
    var Page = require('../pages/page.js');
    var unpacker = require('../message/struct/unpacker/unpacker.js');

    module.exports = Lobby;

    function Lobby(user) {
        this.user = user;
        this.rpcs = null;

        this.needDraw = 0;
        this.expectTime = 0;
        this.waitingNum = 0;
        this.timer = setInterval(this.Draw.bind(this), 1000);
    }

    var proto = Lobby.prototype;

    proto.Login = function () {
        var pm = new _.PlayerLogin();
        pm.setLevel(1);
        var pmdata = pm.serializeBinary();
        var data = Buffer.concat([unpacker.uint8.packerl(14), unpacker.uint16.packerl(440), unpacker.uint16.packerl(pmdata.length), Buffer.from(pmdata)]);
        this.user.RPC(def.ServerTypeLobby, "PlayerLogin", data);
        this.user.ShowPage('lobby');
    };

    proto.StartMatch = function () {
        var data = serializer.uint32(1);
        this.user.RPC(def.ServerTypeLobby, "EnterRoomReq", data);
    };

    proto.RecvRPCMsg = function (msg) {
        if (this.rpcs == null) {
            this.rpcs = {};
            this.rpcs.ExpectTime = this.RPC_ExpectTime.bind(this);
            this.rpcs.NotifyWaitingNums = this.RPC_NotifyWaitingNums.bind(this);
            this.rpcs.MatchSuccess = this.RPC_MatchSuccess.bind(this);
            // this.rpcs.UpdateAirLeft = this.user.room.RPC_UpdateAirLeft.bind(this.user.room);
            this.rpcs.UpdateTotalNum = this.user.room.RPC_UpdateTotalNum.bind(this.user.room);
            this.rpcs.UpdateKillNum = this.user.room.RPC_UpdateKillNum.bind(this.user.room);
            this.rpcs.UpdateAliveNum = this.user.room.RPC_UpdateAliveNum.bind(this.user.room);
        }
        if (msg.MethodName in this.rpcs) {
            this.rpcs[msg.MethodName](msg);
            return true;
        } else {
            return false;
        }
    };

    proto.RPC_ExpectTime = function (msg) {
        var t = unserializer.uint64(msg.Data);
        var now = new Date().getTime();
        this.expectTime = now + t * 1000;
        this.needDraw = 1;
    };


    proto.RPC_NotifyWaitingNums = function (msg) {
        this.waitingNum = unserializer.uint32(msg.Data);
        this.needDraw = 1;
    };

    proto.RPC_MatchSuccess = function (msg) {
        var mapid = this.user.mapid = unserializer.uint32(msg.Data);
        console.log('mapid =', mapid);
        PageLobby.scope.txtmatch = "开始加载地图，地图ID: " + String(mapid);
        Page.showPage('maploading');
    };


    proto.Draw = function (msg) {
        if (this.needDraw == 0) {
            return;
        }
        var now = new Date().getTime();
        var lift = this.expectTime - now;
        if (lift < 0) {
            if (this.timer != 0) {
                clearInterval(this.timer);
                this.timer = 0;
            }
            PageLobby.scope.txtmatch = "即将开始, 当前人数: " + String(this.waitingNum);
            PageLobby.scope.$apply();
        } else {
            PageLobby.scope.txtmatch = "匹配时间:" + String(parseInt((this.expectTime - now) / 1000)) + ", 当前人数:" + String(this.waitingNum);
            PageLobby.scope.$apply();
        }
    };

})();