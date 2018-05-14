(function () {
    'use strict';

    var uint16 = require('./datatypes_uint16.js');

    var string = module.exports = {};

    // packer
    string.packerl = function (val) {
        if (val == null) {
            val = "";
        }
        var len = val.length;
        if (!len) {
            val = String(val);
            len = val.length;
        }
        var buf1 = uint16.packerl(len);
        var buf2 = Buffer.from(val);
        return Buffer.concat([buf1, buf2]);
    };

    string.packerb = function (val) {
        if (val == null) {
            val = "";
        }
        var len = val.length;
        if (!len) {
            val = String(val);
            len = val.length;
        }
        var buf1 = uint16.packerb(len);
        var buf2 = Buffer.from(val);
        return Buffer.concat([buf1, buf2]);
    };

    // unpacker
    string.unpackerl = function (buf) {
        var len = uint16.unpackerl(buf);
        var buf2 = buf.slice(uint16.size(), uint16.size() + len);
        return buf2.toString();
    };

    string.unpackerb = function (buf) {
        var len = uint16.unpackerb(buf);
        var buf2 = buf.slice(uint16.size(), uint16.size() + len);
        return buf2.toString();
    };

    // size
    string.sizel = function (buf) {
        var len = uint16.unpackerl(buf);
        return uint16.size() + len;
    };

    string.sizeb = function (buf) {
        var len = uint16.unpackerb(buf);
        return uint16.size() + len;
    };

})();