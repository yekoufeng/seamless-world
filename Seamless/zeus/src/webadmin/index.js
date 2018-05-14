var app = angular.module('baseApp', ['ui.grid', 'ui.grid.selection']);

app.controller('mainCtrl', ['$scope', '$http', '$window','uiGridConstants', 
    function($scope, $http, $window, uiGridConstants) {
        if (!$window.sessionStorage.username || !$window.sessionStorage.token) {
            $window.location='/login.html';
        }

        $scope.gridOptions = {
            enableRowSelection: true,
            enableRowHeaderSelection: false
        };
        $scope.gridOptions.multiSelect = false;
        $scope.gridOptions.columnDefs = [
            { field: '服务器编号' },
            { field: '类型' },
            { field: '内网地址' },
            { field: '外网地址' },
            { field: '当前负载' },
            { field: '状态' },
        ];

        var serverTypes = {
            0:'未知',
            1:'网关',
            11:'匹配',
            12:'聊天',
            13:'场景'
        };

        $scope.gridOptions.onRegisterApi = function(gridApi){
        //set gridApi on scope
        $scope.gridApi = gridApi;
        gridApi.selection.on.rowSelectionChanged($scope,function(row){
            $scope.selectedServerID = row.entity['服务器编号'];
            $scope.selectedType = row.entity['类型'];
            $scope.selectedStatus = row.entity['状态'];
        });
        };

        $scope.hideResp = true
        $scope.adminCmd = function() {
            $scope.hideResp = true
            var successCallback = function (response) {
                $scope.cmdresp = response.data
                $scope.hideResp = false

                $scope.getServers()
            };

            var errorCallback = function (response) {
                $scope.cmdresp = "请求失败"
                $scope.hideResp = false
            };
            
            $http({
                method: 'POST',
                url: $window.sessionStorage.adminserver+'/cmds',
                data: {
                    "ServerID": $scope.selectedServerID,
                    "Command": $scope.cmd
                }
            }).then(successCallback, errorCallback);
        }

        $scope.getServers = function() {
            var successCallback = function (response) {
                var gridData = [];
                for (var key in response.data) {
                    var type = serverTypes[response.data[key].Type];
                    if (type == undefined) {
                        type = "未知";
                    }
                    gridData.push({
                        "服务器编号": response.data[key].ServerID,
                        "类型": type,
                        "内网地址": response.data[key].InnerAddress,
                        "外网地址": response.data[key].OuterAddress,
                        "当前负载": response.data[key].Load,
                        "状态": response.data[key].Status
                    });
                }

                $scope.gridOptions.data = gridData;
            };

            var errorCallback = function (response) {
                $scope.cmdresp = "请求失败"
                $scope.hideResp = false
            };

            $http({
                method: 'GET',
                url: $window.sessionStorage.adminserver+'/servers'
            }).then(successCallback, errorCallback);
        };

        $scope.getServers()

}]);

app.factory('authInterceptor', function ($rootScope, $q, $window) {
    return {
        request: function (config) {
          config.headers = config.headers || {};
          if ($window.sessionStorage.token) {
              config.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
          }
          return config;
        },
        response: function (response) {
            console.log(response.status)
            if (response.status === 401) {
                $window.location='/login.html';
            }
            return response || $q.when(response);
        }
    };
});

app.config(function ($httpProvider) {
    $httpProvider.interceptors.push('authInterceptor');
});