(function () {
    'use strict';

    // 所有已接协议

    module.exports = NetMsgMsgId;

    function NetMsgMsgId() { }

    NetMsgMsgId.ClientVertifyReqMsgID = 1;              // S -> C
    NetMsgMsgId.ClientVertifySucceedRetMsgID = 2;       // S -> C
    NetMsgMsgId.ClientVertifyFailedRetMsgID = 3;        // S -> C
    NetMsgMsgId.HeartBeatMsgID = 4;                     // S -> C
    NetMsgMsgId.ProtoSyncMsgID = 10;                    // S -> C
    NetMsgMsgId.MRolePropsSyncClientMsgID = 32;         // S -> C
    NetMsgMsgId.EnterSpaceMsgID = 45;                   // S -> C
    NetMsgMsgId.SpaceEntityMsgID = 48;                  // C -> S
    NetMsgMsgId.RPCMsgID = 58;                          // S -> C; C -> S
    NetMsgMsgId.SpaceUserConnectMsgID = 61;             // C -> S
    NetMsgMsgId.SpaceUserConnectSucceedRetMsgID = 62;   // S -> C
    NetMsgMsgId.SyncUserStateMsgID = 63;                // C -> S
    NetMsgMsgId.AOISyncUserStateMsgID = 64;             // S -> C
    NetMsgMsgId.AdjustUserStateMsgID = 65;              // S -> C
    NetMsgMsgId.EntityAOISMsgID = 66;                   // S -> C


    /*
        Protobuf, C -> S :

            PlayerLogin : C -> Lobby
            ParachuteReady : C -> Room

    */


    /*
        RPC,  S -> C :
    
            ExpectTime : uint64
            NotifyWaitingNums : uint32
            MatchSuccess : uint32, uint32
    
    */

})();