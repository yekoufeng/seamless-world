(function () {
    'use strict';

    var uint48 = module.exports = {};

    // packer
    uint48.packerl = function (val) {
        var buf = Buffer.allocUnsafe(uint48.size());
        var v0 = val & 0xFF;
        var v1 = (val >> 8) & 0xFF;
        var v2 = (val >> 16) & 0xFF;
        var v3 = (val >> 24) & 0xFF;
        var v4 = (val / 0x100000000) & 0xFF;
        var v5 = (val / 0x10000000000) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        buf.writeUInt8(v3, 3, true);
        buf.writeUInt8(v4, 4, true);
        buf.writeUInt8(v5, 5, true);
        return buf;
    };

    uint48.packerb = function (val) {
        var buf = Buffer.allocUnsafe(uint48.size());
        var v5 = val & 0xFF;
        var v4 = (val >> 8) & 0xFF;
        var v3 = (val >> 16) & 0xFF;
        var v2 = (val >> 24) & 0xFF;
        var v1 = (val / 0x100000000) & 0xFF;
        var v0 = (val / 0x10000000000) & 0xFF;
        buf.writeUInt8(v0, 0, true);
        buf.writeUInt8(v1, 1, true);
        buf.writeUInt8(v2, 2, true);
        buf.writeUInt8(v3, 3, true);
        buf.writeUInt8(v4, 4, true);
        buf.writeUInt8(v5, 5, true);
        return buf;
    };

    // unpacker
    uint48.unpackerl = function (buf) {
        var v0 = buf.readUInt8(0);
        var v1 = buf.readUInt8(1);
        var v2 = buf.readUInt8(2);
        var v3 = buf.readUInt8(3);
        var v4 = buf.readUInt8(4);
        var v5 = buf.readUInt8(5);
        var v = (v3 << 24) + (v2 << 16) + (v1 << 8) + v0;
        v >>>= 0;
        v += (v5 * 0x10000000000) + (v4 * 0x100000000);
        return v;
    };

    uint48.unpackerb = function (buf) {
        var v5 = buf.readUInt8(0);
        var v4 = buf.readUInt8(1);
        var v3 = buf.readUInt8(2);
        var v2 = buf.readUInt8(3);
        var v1 = buf.readUInt8(4);
        var v0 = buf.readUInt8(5);
        var v = (v3 << 24) + (v2 << 16) + (v1 << 8) + v0;
        v >>>= 0;
        v += (v5 * 0x10000000000) + (v4 * 0x100000000);
        return v;
    };

    // size
    uint48.size = function () {
        return 6;
    };

})();