(function () {
    'use strict';

    var Page = require('./page.js');

    module.exports = PageStage;

    function PageStage() { }

    PageStage.onLoad = function (app) {
        Page.loadPage(app, 'login');
        Page.loadPage(app, 'lobby');
        Page.loadPage(app, 'maploading');
        Page.loadPage(app, 'room');
    };

})();