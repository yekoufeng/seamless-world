package iserver

// ServerInfo 服务器信息
type ServerInfo struct {
	ServerID     uint64 `redis:"serverid"`
	Type         uint8  `redis:"type"`
	OuterAddress string `redis:"outeraddr"`
	InnerAddress string `redis:"inneraddr"`
	Console      uint64 `redis:"console"`
	Load         int    `redis:"load"`
	Token        string `redis:"token"`
	Status       int    `redis:"status"`
}

// ServerList 服务器列表
type ServerList []*ServerInfo

// Len 实现Len方法
func (list ServerList) Len() int {
	return len(list)
}

// Swap 实现Swap方法
func (list ServerList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// Less 实现Less方法
func (list ServerList) Less(i, j int) bool {
	return list[i].Load < list[j].Load
}
