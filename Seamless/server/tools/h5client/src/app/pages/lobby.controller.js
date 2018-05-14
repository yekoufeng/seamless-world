(function () {
    'use strict';

    module.exports = PageLobby;

    function PageLobby() { }

    PageLobby.onController = function ($scope, $http, user) {
        $scope.txtuid = 0;
        $scope.txtentityid = 0;
        $scope.txtmatch = "空闲";
        $scope.clickMatch1 = function () {
            onClickMatch1();
        };
        $scope.clickMatch2 = function () {
            onClickMatch1();
        };
        $scope.clickMatch4 = function () {
            onClickMatch1();
        };

        function onClickMatch1() {
            $scope.txtmatch = "开始匹配...";
            user.lobby.StartMatch();
        }
        function onClickMatch2() {
        }
        function onClickMatch4() {
        }
    };

    PageLobby.onShow = function () {
        if (PageLobby.scope.enable) {
            PageLobby.scope.txtuid = PageLobby.user.uid;
            PageLobby.scope.txtentityid = PageLobby.user.entityID;
        }
    };

})();