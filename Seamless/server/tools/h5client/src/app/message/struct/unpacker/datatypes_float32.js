(function () {
    'use strict';

    var float32 = module.exports = {};

    // packer
    float32.packerl = function (val) {
        var buf = Buffer.allocUnsafe(float32.size());
        buf.writeFloatLE(val, 0, true);
        return buf;
    };

    float32.packerb = function (val) {
        var buf = Buffer.allocUnsafe(float32.size());
        buf.writeFloatBE(val, 0, true);
        return buf;
    };

    // unpacker
    float32.unpackerl = function (buf) {
        return buf.readFloatLE(0);
    };

    float32.unpackerb = function (buf) {
        return buf.readFloatBE(0);
    };

    // size
    float32.size = function () {
        return 4;
    };

})();