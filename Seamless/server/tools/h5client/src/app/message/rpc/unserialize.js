(function () {
    'use strict';

    var unserializer = module.exports = {};
    var unpacker = require('../struct/unpacker/unpacker.js');


    unserializer.typeUint8 = 1;
    unserializer.typeUint16 = 2;
    unserializer.typeUint32 = 3;
    unserializer.typeUint64 = 4;
    unserializer.typeInt8 = 5;
    unserializer.typeInt16 = 6;
    unserializer.typeInt32 = 7;
    unserializer.typeInt64 = 8;
    unserializer.typeFloat32 = 9;
    unserializer.typeFloat64 = 10;
    unserializer.typeString = 11;
    unserializer.typeBytes = 12;
    unserializer.typeBool = 13;
    unserializer.typeProto = 14;

    unserializer.Serialize = function () {
        // TODO:
    };

    unserializer.uint64 = function (buf) {
        return unpacker.uint64.unpackerl(buf.slice(1));
    };
    unserializer.uint32 = function (buf) {
        return unpacker.uint32.unpackerl(buf.slice(1));
    };

})();