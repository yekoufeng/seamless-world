(function () {
    'use strict';

    var uint40 = module.exports = {};

    // packer
    uint40.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint40.size());
        var v0 = val & 0xFF;
        var v1 = (val >> 8) & 0xFF;
        var v2 = (val >> 16) & 0xFF;
        var v3 = (val >> 24) & 0xFF;
        var v4 = (val / 0x100000000) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        buf.writeUInt8(v3, 3, true);
        buf.writeUInt8(v4, 4, true);
        return buf;
    };

    uint40.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint40.size());
        var v4 = val & 0xFF;
        var v3 = (val >> 8) & 0xFF;
        var v2 = (val >> 16) & 0xFF;
        var v1 = (val >> 24) & 0xFF;
        var v0 = (val / 0x100000000) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        buf.writeUInt8(v3, 3, true);
        buf.writeUInt8(v4, 4, true);
        return buf;
    };

    // unpacker
    uint40.unpackerl = function (buf) {
        var v0 = buf.readUInt8(0);
        var v1 = buf.readUInt8(1);
        var v2 = buf.readUInt8(2);
        var v3 = buf.readUInt8(3);
        var v4 = buf.readUInt8(4);
        var v = (v3 << 24) + (v2 << 16) + (v1 << 8) + v0;
        v >>>= 0;
        v += v4 * 0x100000000;
        return v;
    };

    uint40.unpackerb = function (buf) {
        var v4 = buf.readUInt8(0);
        var v3 = buf.readUInt8(1);
        var v2 = buf.readUInt8(2);
        var v1 = buf.readUInt8(3);
        var v0 = buf.readUInt8(4);
        var v = (v3 << 24) + (v2 << 16) + (v1 << 8) + v0;
        v >>>= 0;
        v += v4 * 0x100000000;
        return v;
    };

    // size
    uint40.size = function () {
        return 5;
    };

})();