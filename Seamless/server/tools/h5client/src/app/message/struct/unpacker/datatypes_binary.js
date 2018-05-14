(function () {
    'use strict';

    var uint16 = require('./datatypes_uint16.js');

    var binary = module.exports = {};

    // packer
    binary.packerl = function (val) {
        if (val == null) {
            val = "";
        }
        var buf2 = Buffer.from(val);
        var len = buf2.length;
        var buf1 = uint16.packerl(len);
        return Buffer.concat([buf1, buf2]);
    };

    binary.packerb = function (val) {
        if (val == null) {
            val = "";
        }
        var buf2 = Buffer.from(val);
        var len = buf2.length;
        var buf1 = uint16.packerb(len);
        return Buffer.concat([buf1, buf2]);
    };

    // unpacker
    binary.unpackerl = function (buf) {
        var len = uint16.unpackerl(buf);
        return buf.slice(uint16.size(), uint16.size() + len);
    };

    binary.unpackerb = function (buf) {
        var len = uint16.unpackerb(buf);
        return buf.slice(uint16.size(), uint16.size() + len);
    };

    // size
    binary.sizel = function (buf) {
        var len = uint16.unpackerl(buf);
        return uint16.size() + len;
    };

    binary.sizeb = function (buf) {
        var len = uint16.unpackerb(buf);
        return uint16.size() + len;
    };

})();