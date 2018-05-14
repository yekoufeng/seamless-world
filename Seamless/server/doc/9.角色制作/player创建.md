### player创建

玩家在LoginServer登录完成后，选择一个GatewayServer并在gateway中创建player

```go
CreateEntityAll("Player", uid, "", true)
```

一共有三个服务器注册过player这个实体，分别为GatewayServer，LobbyServer以及RoomServer，由于RoomServer中的player为空间实体，CreateEntityAll不会直接创建，
故GatewayServer调用CreateEntityAll只会在GatewayServer和LobbyServer中创建player实体。


RoomServer中的player是在消息处理函数MsgProc_EnterSpaceReq中创建的

```go
//room上创建
SpaceMsgProc::MsgProc_EnterSpaceReq(content msgdef.IMsg)
{
  //创建了player, room上也就是roomuser
	err := proc.space.AddEntity(msg.EntityType, msg.EntityID, msg.DBID, params[0], false, false)
	if err != nil {
		proc.space.Error("Add entity error ", err, msg)
		return
	}
}
```
