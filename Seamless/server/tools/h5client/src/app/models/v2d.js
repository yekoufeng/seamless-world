(function () {
    'use strict';

    var Camera = require('./v2d_camera.js');

    module.exports = V2D;

    function V2D(user, room) {
        this.user = user;
        this.room = room;
        this.context = null;
        this.canvas = null;
        this.bg = null;
        this.bgSize = 0;
        this.camera = new Camera();
    }

    var NET_SCALE = 100;
    var MAP_SIZE = 8192;
    var CELL_SIZE = 50;
    var DEFAULT_ENTITY_SIZE = 30;

    var proto = V2D.prototype;

    proto.initCanvas = function () {
        var self = this;
        if (self.canvas == null) {
            self.canvas = document.getElementById('scene');
            self.context = self.canvas.getContext("2d");
            var w = angular.element(self.user.window);
            w.bind("keydown", self.doKeyDown.bind(self));
            w.bind("keyup", self.doKeyUp.bind(self));
            w.bind('resize', self.resizeCanvas.bind(self));
            self.resizeCanvas();
        }
    };

    proto.doKeyDown = function (e) {
        var keyID = e.keyCode ? e.keyCode : e.which;
        if (keyID === 38 || keyID === 87) { // up arrow and W
            this.room.downFlagZ = -1;
        }
        if (keyID === 39 || keyID === 68) { // right arrow and D
            this.room.downFlagX = 1;
        }
        if (keyID === 40 || keyID === 83) { // down arrow and S
            this.room.downFlagZ = 1;
        }
        if (keyID === 37 || keyID === 65) { // left arrow and A
            this.room.downFlagX = -1;
        }
    };

    proto.doKeyUp = function (e) {
        var keyID = e.keyCode ? e.keyCode : e.which;
        if (keyID === 38 || keyID === 87) { // up arrow and W
            this.room.downFlagZ = 0;
        }
        if (keyID === 39 || keyID === 68) { // right arrow and D
            this.room.downFlagX = 0;
        }
        if (keyID === 40 || keyID === 83) { // down arrow and S
            this.room.downFlagZ = 0;
        }
        if (keyID === 37 || keyID === 65) { // left arrow and A
            this.room.downFlagX = 0;
        }
    };

    proto.resizeCanvas = function () {
        var self = this;
        if (self.canvas != null) {
            var w = angular.element(self.user.window);
            if (w.length > 0) {
                self.canvas.width = w[0].innerWidth;
                self.canvas.height = w[0].innerHeight;
                self.bgSize = parseInt(Math.ceil(Math.max(self.canvas.width, self.canvas.height) / CELL_SIZE)) * CELL_SIZE;
                self.initBG();
                self.showBG();
            }
        }
    };

    proto.showBG = function () {
        var self = this;
        var selfEntity = self.room.entitys[self.user.entityID];
        if (selfEntity == null) {
            return;
        }
        var posX = parseInt(selfEntity.posX * NET_SCALE);
        var posY = parseInt(selfEntity.posZ * NET_SCALE);

        var mapsize = MAP_SIZE * NET_SCALE;
        var viewdata = self.camera.getView(
            self.canvas.width,
            self.canvas.height,
            posX,
            posY,
            mapsize,
            mapsize
        );

        if (self.canvas != null & self.context != null && self.bg != null) {
            var bgx = parseInt(Math.floor(viewdata.xbegin / self.bgSize)) * self.bgSize;
            var bgy = parseInt(Math.floor(viewdata.ybegin / self.bgSize)) * self.bgSize;
            self.context.putImageData(self.bg, bgx - viewdata.xbegin, bgy - viewdata.ybegin);
            self.context.putImageData(self.bg, bgx - viewdata.xbegin + self.bgSize, bgy - viewdata.ybegin);
            self.context.putImageData(self.bg, bgx - viewdata.xbegin, bgy - viewdata.ybegin + self.bgSize);
            self.context.putImageData(self.bg, bgx - viewdata.xbegin + self.bgSize, bgy - viewdata.ybegin + self.bgSize);
        }

        for (var key in self.room.entitys) {
            var entity = self.room.entitys[key];
            var tmpposx = parseInt(entity.posX * NET_SCALE) * self.camera.scale - viewdata.xbegin;
            var tmpposy = parseInt(entity.posZ * NET_SCALE) * self.camera.scale - viewdata.ybegin;
            self.drawArc(tmpposx, tmpposy, DEFAULT_ENTITY_SIZE, "#FF0000FF");
        }

        posX = posX * self.camera.scale - viewdata.xbegin;
        posY = posY * self.camera.scale - viewdata.ybegin;
        self.drawArc(posX, posY, DEFAULT_ENTITY_SIZE, "#00FF00FF");
    };

    proto.drawArc = function (x, y, r, color) {
        this.context.beginPath();
        this.context.arc(x, y, r, 0, 360, false);
        this.context.fillStyle = color;
        this.context.closePath();
        this.context.fill();
    };

    proto.initBG = function () {
        var self = this;
        var bg = self.context.createImageData(self.bgSize, self.bgSize);
        var step = CELL_SIZE;
        for (var row = 0; row < self.bgSize; row += step) {
            for (var col = 0; col < self.bgSize; col += step) {
                var R, G, B = 0;
                if ((parseInt(col / step) + parseInt(row / step)) % 2 == 0) {
                    R = 69;
                    G = 139;
                    B = 0;
                } else {
                    R = 0;
                    G = 139;
                    B = 0;
                }
                for (var sr = 0; sr < step; sr++) {
                    for (var sc = 0; sc < step; sc++) {
                        var index = (row + sr) * self.bgSize + (col + sc);
                        if (sr == 0 || sr == step - 1 || sc == 0 || sc == step - 1) {
                            bg.data[index * 4 + 0] = 0;
                            bg.data[index * 4 + 1] = 0;
                            bg.data[index * 4 + 2] = 0;
                            bg.data[index * 4 + 3] = 255;
                        } else {
                            bg.data[index * 4 + 0] = R;
                            bg.data[index * 4 + 1] = G;
                            bg.data[index * 4 + 2] = B;
                            bg.data[index * 4 + 3] = 255;
                        }
                    }
                }
            }
        }
        self.bg = bg;
    };

})();