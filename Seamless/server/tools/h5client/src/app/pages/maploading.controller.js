(function () {
    'use strict';

    var Page = require('./page.js');

    module.exports = PageMaploading;

    function PageMaploading() { }

    PageMaploading.onController = function ($scope, $http, user) {
        $scope.txttips = "地图加载中...";
    };

})();