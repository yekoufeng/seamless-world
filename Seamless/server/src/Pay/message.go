package main

// SDK appid
const (
	SDKAPPID = "e5a68bc8f6ac565d0f866d7e78a2613"
)

// 米大师appid
const (
	AndroidAppid = "1450013440"
	IOSAppid     = "1450013441"
)

// 根据登录类型获取midas需要用的appid
func GetMidasAppidByType(typ string) string {
	if typ == "android" {
		return AndroidAppid
	} else if typ == "iap" {
		return IOSAppid
	}

	return ""
}

// 根据操作系统类型获取msdk需要用的appid
func GetMSDKAppidByType(typ string) string {

	if typ == "desktop_m_qq" {
		return "1106393072"
	} else if typ == "desktop_m_wx" {
		return "wxa916d09c4b4ef98f"
	}

	return ""
}

const (
	// ErrPayJSONDecodeFailed 解析json出错
	ErrPayJSONDecodeFailed = -102

	// ErrPayGetPfAndPfkeyFailed 获取pf和pfkey失败
	ErrPayGetPfAndPfkeyFailed = -103

	// ErrPayNewRequestFailed 生成向midas的请求失败
	ErrPayNewRequestFailed = -104

	// ErrPayReqMidasFailed 请求midas失败
	ErrPayReqMidasFailed = -105

	// ErrPayReadMidasRetFailed 读取midas返回数据失败
	ErrPayReadMidasRetFailed = -106

	// ErrPayBuyPropIDFailed 购买物品时物品id错误
	ErrPayBuyPropIDFailed = -107
)

// QueryBalace 客户端查询余额请求
type QueryBalaceByClient struct {
	SessionID      string `json:"SessionID"`      // 用户账户类型
	SessionType    string `json:"SessionType"`    // session类型
	Openid         string `json:"Openid"`         // openid
	Openkey        string `json:"Openkey"`        // openkey
	PayToken       string `json:"PayToken"`       // pay_token
	AccessToken    string `json:"AccessToken"`    // 登录态(qq使用paytoken，微信使用accesstoken)
	Platform       string `json:"Platform"`       //平台标识(一般情况下：qq对应值为desktop_m_qq，wx对应值为desktop_m_wx)
	RegChannel     string `json:"RegChannel"`     //注册渠道
	Os             string `json:"Os"`             //系统(安卓对应android，ios对应iap)
	Installchannel string `json:"Installchannel"` //安装渠道
	Offerid        string `json:"Offerid"`        //支付的appid
}

// MidasQueryBalanceTssList midas查询虚拟币余额结果月卡信息
type MidasQueryBalanceTssList struct {
	innerproductid         uint32 `json:"innerproductid"`         // 用户开通的订阅物品id
	begintime              int64  `json:"begintime"`              // 用户订阅的开始时间
	endtime                int64  `json:"endtime"`                // 用户订阅的结束时间
	paychan                string `json:"paychan"`                // 用户订阅该物品id最后一次的支付渠道
	paysubchan             string `json:"paysubchan"`             //用户订阅该物品 id 最后一次的支付子渠道 id
	autopaychan            string `json:"autopaychan"`            //预留扩展字段，目前没有使用
	autopaysubchan         string `json:"autopaysubchan"`         //预留扩展字段，目前没有使用
	grandtotal_opendays    uint32 `json:"grandtotal_opendays"`    //用户订阅累计开通天数
	grandtotal_presentdays uint32 `json:"grandtotal_presentdays"` //用户订阅累计赠送天数
	first_buy_time         int64  `json:"first_buy_time"`         //首充开通时间
	extend                 string `json:"extend"`                 //预留扩展字段，目前没有使用
}

// MidasQueryBalanceResult midas查询虚拟币余额结果
type MidasQueryBalanceResult struct {
	Ret         uint32                      `json:"ret"`         // 返回码
	Balance     uint32                      `json:"balance"`     // 虚拟币个数
	Gen_balance uint32                      `json:"gen_balance"` // 赠送虚拟币个数
	Save_amt    uint32                      `json:"save_amt"`    // 累计充值金额的虚拟币数量
	Gen_expire  uint32                      `json:"gen_expire"`  // 该字段已作废
	Tiss_list   []*MidasQueryBalanceTssList `json:"tiss_list"`   // 月卡信息字段, 如果没有月卡该字段值为空
}

// RetQueryBalanceResultToClient 返回查询余额结果至客户端
type RetQueryBalanceResultToClient struct {
	Ret     uint32 `json:"Ret"`     // 查询结果
	Balance uint32 `json:"Balance"` // 虚拟币个数
}

// DeductVirtualCoin 客户端扣除虚拟货币请求
type DeductVirtualCoin struct {
	SessionID      string `json:"SessionID"`      // 用户账户类型
	SessionType    string `json:"SessionType"`    // session类型
	Openid         string `json:"Openid"`         // openid
	Openkey        string `json:"Openkey"`        // openkey
	PayToken       string `json:"PayToken"`       // pay_token
	AccessToken    string `json:"AccessToken"`    // 登录态(qq使用paytoken，微信使用accesstoken)
	Platform       string `json:"Platform"`       //平台标识(一般情况下：qq对应值为desktop_m_qq，wx对应值为desktop_m_wx)
	RegChannel     string `json:"RegChannel"`     //注册渠道
	Os             string `json:"Os"`             //系统(安卓对应android，ios对应iap)
	Installchannel string `json:"Installchannel"` //安装渠道
	Offerid        string `json:"Offerid"`        //支付的appid
	productID      uint32 `json:"productID"`      //产品id
}

// MidasDeductVirtualCoinResult midas扣除虚拟货币结果
type MidasDeductVirtualCoinResult struct {
	Ret          uint32 // `json:"ret"`返回码
	Billno       string // `json:"billno"`预扣流水号
	Balance      uint32 // `json:"balance"`预扣后的余额
	Used_gen_amt uint32 // `json:"used_gen_amt"`支付使用赠送币金额
}

// RetDeductVirtualCoinResultToClient 返回扣除虚拟币结果至客户端
type RetDeductVirtualCoinResultToClient struct {
	Ret     uint32 `json:"Ret"`     // 结果
	Balance uint32 `json:"Balance"` // 预扣后的余额
}

// GetpfAndpfkeyData 获取pf和pfkey需要的数据
type GetpfAndpfkeyData struct {
	Appid          string `json:"appid"`          // 游戏的唯一标识
	Openid         string `json:"openid"`         // 用户的唯一标识
	AccessToken    string `json:"accessToken"`    // openkey
	Platform       string `json:"platform"`       // pay_token
	RegChannel     string `json:"regChannel"`     // 登录态(qq使用paytoken，微信使用accesstoken)
	Os             string `json:"os"`             //系统(安卓对应android，ios对应iap)
	Installchannel string `json:"installchannel"` //安装渠道
	Offerid        string `json:"offerid"`        //支付的appid
}

// GetpfAndpfkeyRet 获取pf和pfkey返回信息
type GetpfAndpfkeyRet struct {
	Ret int `json:"ret"` // 返回码 0：正确，其它：失败

	Msg   string `json:"msg"`   // ret非0，则表示“错误码，错误提示”，详细注释参见错误码描述
	Pf    string `json:"pf"`    // 对应的pf值
	PfKey string `json:"pfKey"` // 对应的pfKey值
}
