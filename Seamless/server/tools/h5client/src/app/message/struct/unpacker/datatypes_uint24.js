(function () {
    'use strict';

    var uint24 = module.exports = {};

    // packer
    uint24.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint24.size());
        var v0 = val & 0xFF;
        var v1 = (val >> 8) & 0xFF;
        var v2 = (val >> 16) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        return buf;
    };

    uint24.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint24.size());
        var v2 = val & 0xFF;
        var v1 = (val >> 8) & 0xFF;
        var v0 = (val >> 16) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        return buf;
    };

    // unpacker
    uint24.unpackerl = function (buf) {
        var v0 = buf.readUInt8(0);
        var v1 = buf.readUInt8(1);
        var v2 = buf.readUInt8(2);
        return (v2 << 16) + (v1 << 8) + v0;
    };

    uint24.unpackerb = function (buf) {
        var v2 = buf.readUInt8(0);
        var v1 = buf.readUInt8(1);
        var v0 = buf.readUInt8(2);
        return (v2 << 16) + (v1 << 8) + v0;
    };

    // size
    uint24.size = function () {
        return 3;
    };

})();