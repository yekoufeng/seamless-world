package gameconfig

import "testing"
import "fmt"

func Test_Config(t *testing.T) {
	config := New("gameconfig_test.json")
	goal := config.Get("game.honor.test")
	fmt.Println(goal.(string))
}
