package main

import "zeus/tsssdk"

func main() {

	// clientip := "192.168.22.44:44"

	// ipStr := strings.Split(clientip, ":")[0]
	// bits := strings.Split(ipStr, ".")
	// b0, _ := strconv.Atoi(bits[0])
	// b1, _ := strconv.Atoi(bits[1])
	// b2, _ := strconv.Atoi(bits[2])
	// b3, _ := strconv.Atoi(bits[3])

	// var ip uint32
	// ip += uint32(b0) << 24
	// ip += uint32(b1) << 16
	// ip += uint32(b2) << 8
	// ip += uint32(b3)

	// fmt.Println(ip)

	tsssdk.Init(5)

	// time.Sleep(1 * time.Second)

	// if err := tsssdk.OnPlayerLogin(1024, "test", 0, 0, 10, 24, "192.168.10.10:6666", "哈哈"); err != nil {
	// 	fmt.Println(err)
	// }

	// tsssdk.OnRecvAntiData(1024, "test", 0, 0, 10, []byte("testantidata"))

	// sig := make(chan bool, 1)
	// <-sig

	// tsssdk.Destroy()

	// if err := tsssdk.OnPlayerLogout(1024, "test", 0, 0, 10); err != nil {
	// 	fmt.Println(err)
	// }
}
