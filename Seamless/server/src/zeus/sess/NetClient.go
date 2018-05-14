package sess

import (
	"net"
	"zeus/iserver"
)

// Dial 创建一个连接
func Dial(protocal string, addr string) (iserver.ISess, error) {
	var conn net.Conn
	var err error

	if protocal == "tcp" {
		if conn, err = net.Dial(protocal, addr); err != nil {
			return nil, err
		}
		// return nil, fmt.Errorf("unknown network %s", protocal)
	}

	sess := newSess(conn, nil, false)
	sess.SetVertify()

	return sess, nil
}
