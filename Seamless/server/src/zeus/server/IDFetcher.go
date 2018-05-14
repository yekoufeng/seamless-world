package server

import "sync/atomic"
import "zeus/dbservice"
import "github.com/cihub/seelog"

/*
	临时ID组成部分

	27位的 srvID
	5位的 startID
	32位的 seed

	|-- srvID 27位 --|-- startID 5位 --|-- seed 20位 --|
	XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
*/

// IDFetcher 获取唯一的TempID号
type IDFetcher struct {
	baseID  uint64
	startID uint64
	seed    uint64
}

func newIDFetcher(srvID uint64) *IDFetcher {
	return &IDFetcher{
		baseID: srvID,
		seed:   0,
	}
}

func (srv *IDFetcher) init() error {
	var err error
	srv.startID, err = dbservice.SrvIDUtil(srv.baseID).GetStartID()
	if err != nil {
		seelog.Error("fetch id error , server ID ", srv.baseID)
		return err
	}

	seelog.Debug("server ", srv.baseID, " start , fetch start ID ", srv.startID)
	return nil
}

// FetchTempID 获取一个唯一ID
func (srv *IDFetcher) FetchTempID() uint64 {
	seed := srv.incSeed()
	id := srv.baseID<<25 | srv.startID<<20 | seed
	return id
}

func (srv *IDFetcher) incSeed() uint64 {
	var n, v uint64
	for {
		v = atomic.LoadUint64(&srv.seed)
		n = v + 1
		if n > 0xFFFFF {
			n = 0
		}
		if atomic.CompareAndSwapUint64(&srv.seed, v, n) {
			break
		}
	}
	return n
}
