(function () {
    'use strict';

    var unpacker = module.exports = {};

    unpacker.binary = require('./datatypes_binary.js');
    unpacker.bool = require('./datatypes_bool.js');
    unpacker.float32 = require('./datatypes_float32.js');
    unpacker.int8 = require('./datatypes_int8.js');
    unpacker.int16 = require('./datatypes_int16.js');
    unpacker.int24 = require('./datatypes_int24.js');
    unpacker.int32 = require('./datatypes_int32.js');
    unpacker.int40 = require('./datatypes_int40.js');
    unpacker.int48 = require('./datatypes_int48.js');
    unpacker.int64 = require('./datatypes_int64.js');
    unpacker.string = require('./datatypes_string.js');
    unpacker.uint8 = require('./datatypes_uint8.js');
    unpacker.uint16 = require('./datatypes_uint16.js');
    unpacker.uint24 = require('./datatypes_uint24.js');
    unpacker.uint32 = require('./datatypes_uint32.js');
    unpacker.uint40 = require('./datatypes_uint40.js');
    unpacker.uint48 = require('./datatypes_uint48.js');
    unpacker.uint64 = require('./datatypes_uint64.js');

})();