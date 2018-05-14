(function () {
    'use strict';

    var uint32 = module.exports = {};

    // packer
    uint32.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint32.size());
        buf.writeUInt32LE(val, 0, true);
        return buf;
    };

    uint32.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint32.size());
        buf.writeUInt32BE(val, 0, true);
        return buf;
    };

    // unpacker
    uint32.unpackerl = function (buf) {
        return buf.readUInt32LE(0);
    };

    uint32.unpackerb = function (buf) {
        return buf.readUInt32BE(0);
    };

    // size
    uint32.size = function () {
        return 4;
    };

})();