package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"zeus/gameconfig"
)

// ListDir 获取目录下所有文件名
func ListDir(dirPth string) (files []os.FileInfo, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			// files1 = append(files1, dirPth+PthSep+fi.Name())
			// ListDir(dirPth + PthSep + fi.Name())
		} else {
			files = append(files, fi)
		}
	}
	return files, nil
}
func main() {
	err := filepath.Walk(".", func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			fmt.Println("Processing", file.Name())

			defInfo := gameconfig.New(file.Name())
			name := defInfo.Get("name")
			t := NewTemplate(name.(string))

			props := defInfo.Get("props").(map[string]interface{})
			for prop, info := range props {
				infoMap := info.(map[string]interface{})
				typ := infoMap["type"].(string)
				t.AddType(prop, typ)
			}

			targetFileName := fmt.Sprintf("../../src/entitydef/%sDef.go", name)
			f, err := os.Create(targetFileName)
			if err != nil {
				fmt.Println(err)
				return err
			}
			_, err = f.WriteString(t.String())
			if err != nil {
				fmt.Println(err)
				return err
			}
			f.Sync()
			f.Close()

			fmt.Println("Process Done")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
