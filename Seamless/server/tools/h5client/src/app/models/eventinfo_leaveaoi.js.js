(function () {
    'use strict';

    var unpacker = require('../message/struct/unpacker/unpacker.js');

    module.exports = EventInfo_LeaveAOI;

    function EventInfo_LeaveAOI() {
        this.entityID = 0;
    }

    var proto = EventInfo_LeaveAOI.prototype;

    proto.UpdateData = function (buf) {
        var pos = 0;
        this.entityID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();
    };


    proto.PrintInfo = function () {
        console.log('[leave aoi] this.entityID =', this.entityID);
    };

})();