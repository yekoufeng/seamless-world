package main

import (
	"db"
	"fmt"

	"github.com/spf13/viper"
	"github.com/xuri/excelize"
)

func main() {
	viper.SetConfigFile("./config.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("加载 config.json 配置文件失败")
		return
	}

	cdkeyFile := viper.GetString("CdkeyFile")
	LoadcdkeyInfo(cdkeyFile)
}

func LoadcdkeyInfo(cdkeyFile string) {
	fmt.Println("激活码开始加载")

	xlsx, err := excelize.OpenFile(cdkeyFile)
	if err != nil {
		fmt.Println("OpenFile err: ", err)
		return
	}

	fillSucSum := 0 // 填充成功数量
	repeatSum := 0  // 重复填充数量
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows[3:] {
		if row[1] == "" {
			break
		}

		ret := db.FillCDKEY(row[1])
		if ret {
			fillSucSum++
		} else {
			repeatSum++
		}

	}
	fmt.Printf("激活码路径 :%s, 填充成功数量 :%d, 重复填充数量 :%d", cdkeyFile, fillSucSum, repeatSum)
}
