(function () {
    'use strict';

    var serializer = module.exports = {};
    var unpacker = require('../struct/unpacker/unpacker.js');


    serializer.typeUint8 = 1;
    serializer.typeUint16 = 2;
    serializer.typeUint32 = 3;
    serializer.typeUint64 = 4;
    serializer.typeInt8 = 5;
    serializer.typeInt16 = 6;
    serializer.typeInt32 = 7;
    serializer.typeInt64 = 8;
    serializer.typeFloat32 = 9;
    serializer.typeFloat64 = 10;
    serializer.typeString = 11;
    serializer.typeBytes = 12;
    serializer.typeBool = 13;
    serializer.typeProto = 14;

    serializer.Serialize = function () {
        // TODO:
    };

    serializer.uint32 = function (val) {
        var buf1 = unpacker.uint8.packerl(serializer.typeUint32);
        var buf2 = unpacker.uint32.packerl(val);
        return Buffer.concat([buf1, buf2]);
    };


})();