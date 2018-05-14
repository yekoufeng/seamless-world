(function () {
    'use strict';

    var Page = require('../pages/page.js');
    var Login = require('./login.js');
    var Gateway = require('./gateway.js');
    var Lobby = require('./lobby.js');
    var Room = require('./room.js');

    module.exports = User;

    function User() {
        this.uid = 0;
        this.token = "";
        this.gatewayIP = "";
        this.gatewayPort = 0;
        this.login = new Login(this);
        this.gateway = new Gateway(this);
        this.lobby = new Lobby(this);
        this.entityID = 0;
        this.mapid = 0;
        this.spaceID = 0;
        this.roomIP = "";
        this.roomPort = 0;
        this.room = new Room(this);
        this.window = null; //$window obj
    }

    var proto = User.prototype;

    proto.Login = function (data) {
        console.log("user data = ", JSON.stringify(data));

        this.uid = data.UID;
        this.token = data.Token;
        this.gatewayIP = data.LobbyAddr.split(":")[0];
        this.gatewayPort = parseInt(data.LobbyAddr.split(":")[1]);

        // 登录Gateway
        this.gateway.Login();
    };

    proto.RPC = function (ServerType, MethodName, Data) {
        var RPCMsg = require('../message/struct/RPCMsg.js');
        var msg = new RPCMsg();
        msg.ServerType = ServerType;
        msg.SrcEntityID = this.entityID;
        msg.MethodName = MethodName;
        msg.Data = Data;
        var d = msg.encodeWithHead();
        this.gateway.ws.send(d);
    };


    proto.ShowPage = function (page) {
        Page.showPage(page);
    };

    var u = new User();

    User.initUser = function (app) {
        app.factory('user', obj);
        obj.$inject = [
            '$rootScope',
            '$window'
        ];

        function obj($rootScope, $window) {
            u.window = $window;
            return u;
        }
    };

})();