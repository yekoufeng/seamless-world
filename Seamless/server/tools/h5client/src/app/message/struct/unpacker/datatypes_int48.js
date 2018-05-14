(function () {
    'use strict';

    var uint48 = require('./datatypes_uint48.js');

    var int48 = module.exports = {};

    // packer
    int48.packerl = function (val) {
        if (val < 0) {
            return uint48.packerl(0x7FFFFFFFFFFF - val);
        } else {
            return uint48.packerl(val);
        }
    };

    int48.packerb = function (val) {
        if (val < 0) {
            return uint48.packerb(0x7FFFFFFFFFFF - val);
        } else {
            return uint48.packerb(val);
        }
    };

    // unpacker
    int48.unpackerl = function (buf) {
        var v5 = buf.readUInt8(5);
        if ((v5 & 0x80) == 0x80) {
            return 0x7FFFFFFFFFFF - uint48.unpackerl(buf);
        } else {
            return uint48.unpackerl(buf);
        }
    };

    int48.unpackerb = function (buf) {
        var v5 = buf.readUInt8(0);
        if ((v5 & 0x80) == 0x80) {
            return 0x7FFFFFFFFFFF - uint48.unpackerb(buf);
        } else {
            return uint48.unpackerb(buf);
        }
    };

    // size
    int48.size = function () {
        return 6;
    };

})();