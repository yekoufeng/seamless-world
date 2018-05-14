var app = angular.module('baseApp', []);

app.controller('loginCtrl', ['$scope', '$http', '$window', '$location', 
    function ($scope, $http, $window, $location) {
        $scope.hidemessage = true;
        delete $window.sessionStorage.username;
        delete $window.sessionStorage.token;
        delete $window.sessionStorage.adminserver;

        $scope.login = function() {
            var successCallback = function (response) {
                $window.sessionStorage.username = $scope.username;
                $window.sessionStorage.token = response.data.token;
                $window.sessionStorage.adminserver = 'http://' + $scope.adminserver;

                $window.location='/index.html';
            };

            var errorCallback = function (response) {
                delete $window.sessionStorage.token;
                $scope.respmessage = "登录失败";
                $scope.hidemessage = false;
            };
                
            $http({
                method: 'POST',
                url: 'http://' + $scope.adminserver + '/login',
                data: {
                    "username": $scope.username,
                    "password": $scope.password
                }
            }).then(successCallback, errorCallback);
        }

}]);
