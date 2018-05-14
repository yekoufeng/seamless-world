(function () {
    'use strict';

    var unpacker = require('../message/struct/unpacker/unpacker.js');

    module.exports = Entity;

    function Entity() {
        this.entityID = 0;
        this.timeStamp = 0;
        this.mask = 0;
        this.posX = 0;
        this.posY = 0;
        this.posZ = 0;
        this.rotaX = 0;
        this.rotaY = 0;
        this.rotaZ = 0;
        this.param1 = 0;
        this.param2 = 0;
        this.events = null;
        this.Epoch = 0;
    }

    var Mask_Pos_X = 0x00000001;
    var Mask_Pos_Y = 0x00000002;
    var Mask_Pos_Z = 0x00000004;
    var Mask_Rota_X = 0x00000008;
    var Mask_Rota_Y = 0x00000010;
    var Mask_Rota_Z = 0x00000020;
    var Mask_Param_1 = 0x00000040;
    var Mask_Param_2 = 0x00000080;
    var Mask_Events = 0x00000100;

    var proto = Entity.prototype;

    proto.UpdateData = function (buf, bEpoch) {
        var pos = 0;
        this.timeStamp = unpacker.uint32.unpackerl(buf.slice(pos));
        pos += unpacker.uint32.size();

        var mask = this.mask = unpacker.uint32.unpackerl(buf.slice(pos));
        pos += unpacker.uint32.size();

        if ((mask & Mask_Pos_X) == Mask_Pos_X) {
            this.posX = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Pos_Y) == Mask_Pos_Y) {
            this.posY = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Pos_Z) == Mask_Pos_Z) {
            this.posZ = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Rota_X) == Mask_Rota_X) {
            this.rotaX = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Rota_Y) == Mask_Rota_Y) {
            this.rotaY = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Rota_Z) == Mask_Rota_Z) {
            this.rotaZ = unpacker.float32.unpackerl(buf.slice(pos));
            pos += unpacker.float32.size();
        }

        if ((mask & Mask_Param_1) == Mask_Param_1) {
            this.param1 = unpacker.uint64.unpackerl(buf.slice(pos));
            pos += unpacker.uint64.size();
        }

        if ((mask & Mask_Param_2) == Mask_Param_2) {
            this.param2 = unpacker.uint64.unpackerl(buf.slice(pos));
            pos += unpacker.uint64.size();
        }

        if ((mask & Mask_Events) == Mask_Events) {
            this.events = unpacker.binary.unpackerl(buf.slice(pos));
            pos += unpacker.binary.sizel(buf.slice(pos));
        }

        if (bEpoch) {
            this.Epoch = unpacker.uint8.unpackerl(buf.slice(pos));
            pos += unpacker.uint8.size();
        }
    };


    proto.PrintInfo = function () {
        console.log('this.entityID =', this.entityID);
        console.log('this.timeStamp =', this.timeStamp);
        console.log('this.mask =', this.mask);
        console.log('this.posX =', this.posX);
        console.log('this.posY =', this.posY);
        console.log('this.posZ =', this.posZ);
        console.log('this.rotaX =', this.rotaX);
        console.log('this.rotaY =', this.rotaY);
        console.log('this.rotaZ =', this.rotaZ);
        console.log('this.param1 =', this.param1);
        console.log('this.param2 =', this.param2);
        console.log('this.events =', this.events);
    };

})();