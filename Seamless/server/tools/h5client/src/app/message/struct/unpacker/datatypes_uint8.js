(function () {
    'use strict';

    var uint8 = module.exports = {};

    // packer
    uint8.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint8.size());
        buf.writeUInt8(val, 0, true);
        return buf;
    };

    uint8.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint8.size());
        buf.writeUInt8(val, 0, true);
        return buf;
    };

    // unpacker
    uint8.unpackerl = function (buf) {
        return buf.readUInt8(0);
    };

    uint8.unpackerb = function (buf) {
        return buf.readUInt8(0);
    };

    // size
    uint8.size = function () {
        return 1;
    };

})();