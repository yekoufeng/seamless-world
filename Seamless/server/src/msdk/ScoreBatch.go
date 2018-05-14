package msdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

func QQScoreBatchList(appid, msdkKey, openid, accessToken string, platid uint32, name string, uid uint64, lst []*Param) {

	body := &Body{
		Appid:       fmt.Sprintf("%v", appid),
		AccessToken: accessToken,
		Openid:      openid,
		Param: []*Param{
			&Param{
				Tp:      12,
				BCover:  1,
				Data:    fmt.Sprintf("%v", platid),
				Expires: "不过期",
			},
			// &Param{
			// 	Tp:      26,
			// 	BCover:  1,
			// 	Data:    "1",
			// 	Expires: "不过期",
			// },
			// &Param{
			// 	Tp:      27,
			// 	BCover:  1,
			// 	Data:    "1",
			// 	Expires: "不过期",
			// },
			&Param{
				Tp:      28,
				BCover:  1,
				Data:    fmt.Sprintf("%v", uid),
				Expires: "不过期",
			},
			// &Param{
			// 	Tp:      29,
			// 	BCover:  1,
			// 	Data:    name,
			// 	Expires: "不过期",
			// },

			//TODO: qqscorebatch 上报大区ID？
			//TODO: qqscorebatch 服务器信息？
			//qqscorebatch 角色ID？
			//qqscorebatch 角色名称？
		},
	}
	if name != "" {
		body.Param = append(body.Param, &Param{
			Tp:      29,
			BCover:  1,
			Data:    name,
			Expires: "不过期",
		})
	}
	body.Param = append(body.Param, lst...)

	// logstr := fmt.Sprintf("qqscorebatch:\nappid:%v\nmsdkkey:%v\nopenid:%v\naccessToken:%v\nplatid:%v\nname:%v\nuid:%v\nlst:%v",
	// 	appid, msdkKey, openid, accessToken, platid, name, uid, body.Param)
	// log.Debug(logstr)
	// for _, v := range body.Param {
	// 	log.Debug(fmt.Sprintf("tp:%v bcover:%v data:%v expires:%v", v.Tp, v.BCover, v.Data, v.Expires))
	// }

	jsonReader, err := body.JsonReader()
	if err != nil {
		log.Error(err)
		return
	}
	url := viper.GetString("Config.MSDKAddr") + "/profile/qqscore_batch/"
	now := time.Now().Unix()

	sigKey := fmt.Sprintf("%v%v", msdkKey, now)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sigKey))
	cipherStr := md5Ctx.Sum(nil)
	sig := hex.EncodeToString(cipherStr)

	url = fmt.Sprintf("%v?timestamp=%v&appid=%v&sig=%v&openid=%v&encode=2&conn=1", url, now, appid, sig, openid)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", jsonReader)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	// log.Debug("qqscorebatch result ", string(bodyData))
}
func QQScoreBatch(appid, msdkKey, openid, accessToken string, platid uint32, name string, uid uint64, tp, bcover int, data, expires string) {
	var lst []*Param
	if tp != 0 {
		lst = append(lst, &Param{
			Tp:      tp,
			BCover:  bcover,
			Data:    data,
			Expires: expires,
		})
	}
	QQScoreBatchList(appid, msdkKey, openid, accessToken, platid, name, uid, lst)
}

type Param struct {
	Tp      int    `json:"type"`
	BCover  int    `json:"bcover"`
	Data    string `json:"data"`
	Expires string `json:"expires"`
}

type Body struct {
	Appid       string   `json:"appid"`
	AccessToken string   `json:"accessToken"`
	Openid      string   `json:"openid"`
	Param       []*Param `json:"param"`
}

func (body *Body) JsonReader() (*strings.Reader, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(data)), nil
}
