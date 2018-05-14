package iserver

// IUIDFetcher 获取唯一的ID
// 在分布式情况下，保证每个服务器产生的ID号不会重复
type IUIDFetcher interface {
	FetchTempID() uint64
}
