(function () {
    'use strict';

    var uint40 = require('./datatypes_uint40.js');

    var int40 = module.exports = {};

    // packer
    int40.packerl = function (val) {
        if (val < 0) {
            return uint40.packerl(0x7FFFFFFFFF - val);
        } else {
            return uint40.packerl(val);
        }
    };

    int40.packerb = function (val) {
        if (val < 0) {
            return uint40.packerb(0x7FFFFFFFFF - val);
        } else {
            return uint40.packerb(val);
        }
    };

    // unpacker
    int40.unpackerl = function (buf) {
        var v4 = buf.readUInt8(4);
        if ((v4 & 0x80) == 0x80) {
            return 0x7FFFFFFFFF - uint40.unpackerl(buf);
        } else {
            return uint40.unpackerl(buf);
        }
    };

    int40.unpackerb = function (buf) {
        var v4 = buf.readUInt8(0);
        if ((v4 & 0x80) == 0x80) {
            return 0x7FFFFFFFFF - uint40.unpackerb(buf);
        } else {
            return uint40.unpackerb(buf);
        }
    };

    // size
    int40.size = function () {
        return 5;
    };

})();