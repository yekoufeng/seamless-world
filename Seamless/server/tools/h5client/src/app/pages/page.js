(function () {
    'use strict';

    var Util = require('../util/util.js');

    module.exports = Page;

    function Page() { }

    Page.scopes = {};

    Page.loadPage = function (app, page) {

        app.directive('runoob' + Util.toUpper(page), function () {
            return {
                templateUrl: 'app/pages/' + page + '.html'
            };
        });

        var pageX = require('./' + page + '.controller.js');

        var onLoad = pageX.onLoad;
        if (onLoad != null) {
            onLoad(app);
        }

        function ctrl($scope, $http, user) {
            pageX.scope = $scope;
            pageX.user = user;
            Page.scopes[page] = $scope;
            $scope.enable = false;
            onController($scope, $http, user);
        }
        var onController = pageX.onController;
        if (onController != null) {
            app.controller(page, ctrl);
            ctrl.$inject = [
                '$scope',
                '$http',
                'user'
            ];
        }
    };

    Page.showPage = function (page) {
        var pageX = require('./' + page + '.controller.js');
        for (var key in Page.scopes) {
            Page.scopes[key].enable = false;
            Page.scopes[key].$apply();
            if (key != page) {
                var onHide = pageX.onHide;
                if (onHide) {
                    onHide();
                }
            }
        }
        Page.scopes[page].enable = true;
        var onShow = pageX.onShow;
        if (onShow) {
            onShow();
        }
        Page.scopes[page].$apply();
    };

})();