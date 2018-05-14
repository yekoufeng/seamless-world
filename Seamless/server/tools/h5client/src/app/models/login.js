(function () {
    'use strict';

    require('sprintf-js');

    module.exports = Login;

    function Login(user) {
        this.user = user;
    }

    var proto = Login.prototype;

    proto.Login = function ($http, account, password, ip, port) {
        var self = this;
        var UserLoginReq = {
            "User": account,
            "Password": password
        };
        var data = JSON.stringify(UserLoginReq);
        var url = sprintf("http://%s:%s/login", ip, port);

        // TODO: 有空闲时，服务器支持GET请求，并这里改成GET方式。
        $http({
            url: url,
            method: 'POST',
            data: data,
            async: false,
            headers: {
                "Access-Control-Allow-Origin": "*",
                'Access-Control-Allow-Methods': 'POST',
                'Access-Control-Allow-Headers': 'Accept,X-Custom-Header,X-Requested-With,Content-Type,Origin'
            }
        }).then(function success(response) {
            console.log("login to Login success!");
            console.log('response:', response);
            self.user.Login(response.data);
        }, function fail(response) {
            console.log("login to Login fail!");
            console.log('response:', response);
            alert("login fail.\nresponse:" + JSON.stringify(response));
        });
    };



})();