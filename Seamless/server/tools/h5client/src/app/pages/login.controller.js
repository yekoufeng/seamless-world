(function () {
    'use strict';

    module.exports = PageLogin;

    function PageLogin() { }

    PageLogin.onController = function ($scope, $http, user) {
        $scope.enable = true;
        $scope.txtaccount = 'test2';
        $scope.txtpassword = '123456';
        $scope.txtip = '127.0.0.1';
        $scope.txtport = 8080;
        $scope.click = function () {
            onClick();
        };

        function onClick() {
            console.log('txtaccount:', $scope.txtaccount);
            console.log('txtpassword:', $scope.txtpassword);
            console.log('txtip:', $scope.txtip);
            console.log('txtport:', $scope.txtport);

            if ($scope.txtaccount == "") {
                alert("账号名不能为空！");
                return;
            }
            if ($scope.txtpassword == "") {
                alert("密码不能为空！");
                return;
            }

            user.login.Login($http, $scope.txtaccount, $scope.txtpassword, $scope.txtip, $scope.txtport);
        }
    };

})();