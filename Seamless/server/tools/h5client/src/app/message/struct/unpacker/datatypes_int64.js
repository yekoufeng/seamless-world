(function () {
    'use strict';

    var Int64LE = require("int64-buffer").Int64LE;
    var Int64BE = require("int64-buffer").Int64BE;

    var uint64 = module.exports = {};

    // packer
    uint64.packerl = function (val) {
        var v = new Int64LE(val);
        return v.toBuffer();
    };

    uint64.packerb = function (val) {
        var v = new Int64BE(val);
        return v.toBuffer();
    };

    // unpacker
    uint64.unpackerl = function (buf) {
        var v = new Int64LE(buf);
        return v.toNumber();
    };

    uint64.unpackerb = function (buf) {
        var v = new Int64BE(buf);
        return v.toNumber();
    };

    // size
    uint64.size = function () {
        return 8;
    };

})();