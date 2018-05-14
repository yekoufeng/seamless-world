package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"zeus/login"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}

	file := flag.String("file", "账号表(openid).txt", "账号表")
	flag.Parse()

	fp, err := os.Open(*file)
	if err != nil {
		fmt.Println("Open err: ", err)
		return
	}

	app := &login.App{}
	initGrade := viper.GetInt("InitGrade")

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		openid := scanner.Text()
		_, err = app.DoCreateNewUser(openid, "", uint32(initGrade))
		if err != nil {
			fmt.Println("DoCreateNewUser err: ", err)
			return
		}

		fmt.Println("Success:", openid)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
