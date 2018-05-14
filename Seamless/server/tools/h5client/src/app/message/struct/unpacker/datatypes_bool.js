(function () {
    'use strict';

    var bool = module.exports = {};

    // packer
    bool.packerl = function (val) {
        var v = ((!!val) ? 1 : 0) & 0xff;
        var buf = Buffer.allocUnsafe(bool.size());
        buf.writeUInt8(v, 0, true);
        return buf;
    };

    bool.packerb = function (val) {
        var v = ((!!val) ? 1 : 0) & 0xff;
        var buf = Buffer.allocUnsafe(bool.size());
        buf.writeUInt8(v, 0, true);
        return buf;
    };

    // unpacker
    bool.unpackerl = function (buf) {
        return buf.readUInt8(0) != 0;
    };

    bool.unpackerb = function (buf) {
        return buf.readUInt8(0) != 0;
    };

    // size
    bool.size = function () {
        return 1;
    };

})();