# 队伍消息处理

列出消息流程，每一步列出消息的名字，发送方和接收方，并简要说明，可能添加发送方的函数名。

## Lobby\TeamMgrProc.go

func (proc *LobbyUserMsgProc)
* RPC_QuickEnterTeam(mapid uint32, teamtype uint32, automatch uint32) {
* RPC_EnterTeamRet(ret uint32, teamID uint64) {
* RPC_SyncTeamInfoRet(msg *protoMsg.SyncTeamInfoRet) {
* RPC_ChangeTeamType(teamtype uint32) {
* RPC_ChangeAutoMatch(automatch uint32) {
* RPC_ChangeMap(mapid uint32) {
* RPC_QuitTeam() {
* RPC_QuitTeamRet(result uint32) {
* RPC_ConfirmTeamMatch() {
* RPC_InviteReq(otherid uint64) {
* RPC_InviteRsp(teamid uint64) {
* RPC_InviteRspRet(result uint32, oldteam uint64) {
* RPC_AutoEnterTeamRet(result uint32) {
* RPC_AutoEnterTeam(teamid uint64) {
	+ Todo: 从 LobbyUserMsgProc.go 移到 TeamMgrProc.go

## 邀请好友

1. Clt1 -> Lobby: InviteReq, 邀请另一玩家加入队伍，`LobbyUserMsgProc.RPC_InviteReq(otherid uint64)`
1. Lobby -> Clt2: Invite, 邀请
1. Clt2 -> Lobby: InviteRsp, 接受邀请 `onConfirmInvite()`
1. Lobby -> Match: InviteRsp, `LobbyUserMsgProc.RPC_InviteRsp(teamid uint64)`
1. Match -> Lobby: InviteRspRet, 加入队伍
	+ Lobby -> Match: InviteLeaveTeam, 离开旧队伍

## 切换队伍类型
1. Clt -> Lobby: ChangeTeamType
1. Lobby -> Match: ChangeTeamType
	1. Match -> Lobby: SyncTeamInfoRet, 广播队伍信息
	1. Lobby -> Clt: SyncTeamInfoRet

## 快速进入队伍
1. Clt -> Lobby: QuickEnterTeam, onMultiMatchType() 当进入多人模式时发送？ 由单人转成队伍？
1. Lobby -> Match: EnterTeam，新队伍ID, 请求进入队伍, 带旧队伍ID
1. Match -> Lobby: EnterTeamRet, 如果是新队伍才发送

## 退出队伍
1. Clt -> Lobby: QuitTeam 多人模式转单人模式时退出队伍
1. Lobby -> Match: LeaveTeam
1. Match -> Lobby: QuitTeamRet
1. Lobby -> Match: QuitTeamRet

## 自动匹配
1. Clt -> Lobby: ChangeAutoMatch 自动匹配按钮？
1. Lobby -> Match:

## 开始匹配
1. Clt -> Lobby: ConfirmTeamMatch 如果是单人模式，则是 EnterRoomReq 消息流程
1. Lobby -> Match: ConfirmTeamMatch
1. Match -> Match: EnterDuoQueue, 仅队伍未满需自动匹配时，否则直接加入队伍列表

## 切换地图
1. Clt -> Lobby -> Match ChangeMap 客户端未实现？

## 自动进入队伍
1. Clt -> Lobby: AutoEnterTeam 被平台好友邀请组队，点击分享平台唤醒后，自动进入房间
1. Lobby -> Match: AutoEnterTeam
1. Match -> Lobby: AutoEnterTeamRet
