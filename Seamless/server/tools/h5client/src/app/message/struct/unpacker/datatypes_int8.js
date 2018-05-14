(function () {
    'use strict';

    var int8 = module.exports = {};

    // packer
    int8.packerl = function (val) {
        var buf = Buffer.allocUnsafe(int8.size());
        buf.writeInt8(val, 0, true);
        return buf;
    };

    int8.packerb = function (val) {
        var buf = Buffer.allocUnsafe(int8.size());
        buf.writeInt8(val, 0, true);
        return buf;
    };

    // unpacker
    int8.unpackerl = function (buf) {
        return buf.readInt8(0);
    };

    int8.unpackerb = function (buf) {
        return buf.readInt8(0);
    };

    // size
    int8.size = function () {
        return 1;
    };

})();