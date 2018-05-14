(function () {
    'use strict';

    var int32 = module.exports = {};

    // packer
    int32.packerl = function (val) {
        var buf = Buffer.allocUnsafe(int32.size());
        buf.writeInt32LE(val, 0, true);
        return buf;
    };

    int32.packerb = function (val) {
        var buf = Buffer.allocUnsafe(int32.size());
        buf.writeInt32BE(val, 0, true);
        return buf;
    };

    // unpacker
    int32.unpackerl = function (buf) {
        return buf.readInt32LE(0);
    };

    int32.unpackerb = function (buf) {
        return buf.readInt32BE(0);
    };

    // size
    int32.size = function () {
        return 4;
    };

})();