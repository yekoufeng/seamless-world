(function () {
    'use strict';

    var unpacker = require('../message/struct/unpacker/unpacker.js');

    module.exports = EventInfo_EnterAOI;

    function EventInfo_EnterAOI() {
        this.entityID = 0;
        this.entityType = "";
        this.state = null;
        this.propNum = 0;
        this.properties = null;
        this.baseProps = null;
    }

    var proto = EventInfo_EnterAOI.prototype;

    proto.UpdateData = function (buf) {
        var pos = 0;
        this.entityID = unpacker.uint64.unpackerl(buf.slice(pos));
        pos += unpacker.uint64.size();

        this.entityType = unpacker.string.unpackerl(buf.slice(pos));
        pos += unpacker.string.sizel(buf.slice(pos));

        this.state = unpacker.binary.unpackerl(buf.slice(pos));
        pos += unpacker.binary.sizel(buf.slice(pos));

        this.propNum = unpacker.uint16.unpackerl(buf.slice(pos));
        pos += unpacker.uint16.size();

        this.properties = unpacker.binary.unpackerl(buf.slice(pos));
        pos += unpacker.binary.sizel(buf.slice(pos));
    };


    proto.PrintInfo = function () {
        console.log('[enter aoi] this.entityID =', this.entityID);
        console.log('[enter aoi] this.entityType =', this.entityType);
        console.log('[enter aoi] this.state =', this.state);
        console.log('[enter aoi] this.propNum =', this.propNum);
        console.log('[enter aoi] this.properties =', this.properties);
    };

})();