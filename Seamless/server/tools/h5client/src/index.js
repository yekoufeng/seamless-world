(function () {
    'use strict';

    require('angular');
    var Page = require('./app/pages/page.js');
    var User = require('./app/models/user.js');

    function startApp() {
        var app = angular.module("app", []);
        closePreload(app);
        User.initUser(app);
        Page.loadPage(app, 'stage');
    }

    startApp();

})();