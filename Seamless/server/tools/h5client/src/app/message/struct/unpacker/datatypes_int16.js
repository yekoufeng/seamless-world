(function () {
    'use strict';

    var int16 = module.exports = {};

    // packer
    int16.packerl = function (val) {
        var buf = Buffer.allocUnsafe(int16.size());
        buf.writeInt16LE(val, 0, true);
        return buf;
    };

    int16.packerb = function (val) {
        var buf = Buffer.allocUnsafe(int16.size());
        buf.writeInt16BE(val, 0, true);
        return buf;
    };

    // unpacker
    int16.unpackerl = function (buf) {
        return buf.readInt16LE(0);
    };

    int16.unpackerb = function (buf) {
        return buf.readInt16BE(0);
    };

    // size
    int16.size = function () {
        return 2;
    };

})();