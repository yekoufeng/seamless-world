(function () {
    'use strict';

    module.exports = Camera;

    function Camera() {
        this.scale = 1;
    }

    var proto = Camera.prototype;

    proto.getView = function (viewWidth, viewHeight, posX, posY, mapWidth, mapHeight) {
        var self = this;
        var viewdata = {};
        viewdata.width = viewWidth;
        viewdata.height = viewHeight;

        posX = posX * self.scale;
        posY = posY * self.scale;
        mapWidth = mapWidth * self.scale;
        mapHeight = mapHeight * self.scale;

        var centerX = viewdata.width / 2;
        var centerY = viewdata.height / 2;

        if (posX <= centerX) {
            viewdata.xbegin = 0;
            viewdata.xend = centerX * 2 + 1;
        } else if (mapWidth - posX <= centerX) {
            viewdata.xbegin = mapWidth - centerX * 2;
            viewdata.xend = mapWidth;
        } else {
            viewdata.xbegin = posX - centerX;
            viewdata.xend = posX + centerX + 1;
        }
        if (posY <= centerY) {
            viewdata.ybegin = 0;
            viewdata.yend = centerY * 2 + 1;
        } else if (mapHeight - posY <= centerY) {
            viewdata.ybegin = mapHeight - centerY * 2;
            viewdata.yend = mapHeight;
        } else {
            viewdata.ybegin = posY - centerY;
            viewdata.yend = posY + centerY + 1;
        }
        viewdata.xbegin = Math.ceil(viewdata.xbegin);
        viewdata.xend = Math.ceil(viewdata.xend);
        viewdata.ybegin = Math.ceil(viewdata.ybegin);
        viewdata.yend = Math.ceil(viewdata.yend);

        return viewdata;
    };

})();