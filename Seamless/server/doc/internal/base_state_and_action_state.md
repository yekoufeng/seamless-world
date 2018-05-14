# BaseState 和 ActionState

角色有2种状态：BaseState 和 ActionState,
 BaseState 对应下肢的动作，ActionState 对应上肢的动作。

```csharp
public enum ActionState : byte
{
    none = 0,
    rescue = 1, //复活
    potion = 2, //吃药
    shoot = 3, //射击
    shootPrepare = 4, //射击准备
    aim = 5, //瞄准
    throwPrepare = 6, //投掷准备
    Throw = 7, //投掷
    melee = 8, //近战
    meleePrepare = 9, //近战准备
    change = 10, //换武器
    reload = 11, //换弹
}
```

```go
const (
	//RoomPlayerBaseState_Inplane 跳伞准备
	RoomPlayerBaseState_Inplane = 1
	//RoomPlayerBaseState_Glide 跳伞俯冲
	RoomPlayerBaseState_Glide = 2
	//RoomPlayerBaseState_Parachute 跳伞
	RoomPlayerBaseState_Parachute = 3
	//RoomPlayerBaseState_Stand 正常状态(站立，移动)
	RoomPlayerBaseState_Stand = 4
	//RoomPlayerBaseState_Down 匍匐
	RoomPlayerBaseState_Down = 5
	//RoomPlayerBaseState_Ride 载具
	RoomPlayerBaseState_Ride = 6
	//RoomPlayerBaseState_Passenger 乘客
	RoomPlayerBaseState_Passenger = 7
	//RoomPlayerBaseState_Swim 游泳
	RoomPlayerBaseState_Swim = 8
	//RoomPlayerBaseState_Dead 死亡
	RoomPlayerBaseState_Dead = 9
	//RoomPlayerBaseState_WillDie 被击倒
	RoomPlayerBaseState_WillDie = 10
	//RoomPlayerBaseState_Watch 观战
	RoomPlayerBaseState_Watch = 11
	//RoomPlayerBaseState_Crouch 蹲
	RoomPlayerBaseState_Crouch = 12
	//RoomPlayerBaseState_Fall 跌落
	RoomPlayerBaseState_Fall = 13
	//RoomPlayerBaseState_Jump 跳跃
	RoomPlayerBaseState_Jump = 14
	//RoomPlayerBaseState_LeaveMap 离开地图
	RoomPlayerBaseState_LeaveMap = 15
	//RoomPlayerBaseState_LoadingMap 加载地图
	RoomPlayerBaseState_LoadingMap = 100
)
```
cehua\excel\行为状态表(action).xlsx 
定义了动作状态允许的转移状态。

