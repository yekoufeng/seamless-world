(function () {
    'use strict';

    var Uint64LE = require("int64-buffer").Uint64LE;
    var Uint64BE = require("int64-buffer").Uint64BE;

    var uint64 = module.exports = {};

    // packer
    uint64.packerl = function (val) {
        var v = new Uint64LE(val);
        return v.toBuffer();
    };

    uint64.packerb = function (val) {
        var v = new Uint64BE(val);
        return v.toBuffer();
    };

    // unpacker
    uint64.unpackerl = function (buf) {
        var v = new Uint64LE(buf);
        return v.toNumber();
    };

    uint64.unpackerb = function (buf) {
        var v = new Uint64BE(buf);
        return v.toNumber();
    };

    // size
    uint64.size = function () {
        return 8;
    };

})();