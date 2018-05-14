(function () {
    'use strict';

    var uint24 = require('./datatypes_uint24.js');

    var int24 = module.exports = {};

    // packer
    int24.packerl = function (val) {
        if (val < 0) {
            return uint24.packerl(0x7FFFFF - val);
        } else {
            return uint24.packerl(val);
        }
    };

    int24.packerb = function (val) {
        if (val < 0) {
            return uint24.packerb(0x7FFFFF - val);
        } else {
            return uint24.packerb(val);
        }
    };

    // unpacker
    int24.unpackerl = function (buf) {
        var v2 = buf.readUInt8(2);
        if ((v2 & 0x80) == 0x80) {
            return 0x7FFFFF - uint24.unpackerl(buf);
        } else {
            return uint24.unpackerl(buf);
        }
    };

    int24.unpackerb = function (buf) {
        var v2 = buf.readUInt8(0);
        if ((v2 & 0x80) == 0x80) {
            return 0x7FFFFF - uint24.unpackerb(buf);
        } else {
            return uint24.unpackerb(buf);
        }
    };

    // size
    int24.size = function () {
        return 3;
    };

})();