(function () {
    'use strict';

    var uint16 = module.exports = {};

    // packer
    uint16.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint16.size());
        buf.writeUInt16LE(val, 0, true);
        return buf;
    };

    uint16.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint16.size());
        buf.writeUInt16BE(val, 0, true);
        return buf;
    };

    // unpacker
    uint16.unpackerl = function (buf) {
        return buf.readUInt16LE(0);
    };

    uint16.unpackerb = function (buf) {
        return buf.readUInt16BE(0);
    };

    // size
    uint16.size = function () {
        return 2;
    };

})();