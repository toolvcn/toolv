package qq

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type (
	// QQ扫码登录
	qrLoginStruct struct {
		Appid     string        // 应用 ID
		UserAgent string        // 请求 User-Agent 配置
		Timeout   time.Duration // 请求超时时间
	}
	// 图片信息结构
	getQrStruct struct {
		Qrsig string `json:"qrsig"` // Qrsig
		Image string `json:"image"` // base64 图片
	}
	// 登录状态结构
	loginStatusStruct struct {
		Status    int    `json:"-"`    // 状态
		Message   string `json:"-"`    // 消息
		Uin       string `json:"uin"`  // QQ号码
		Name      string `json:"name"` // 昵称
		Skey      string `json:"-"`    // cookie skey
		P_skey    string `json:"-"`    // cookie p_skey
		Pt4_token string `json:"-"`    // cookie pt4_token
		Url       string `json:"-"`    // 登录地址
	}
)

func NewQrLogin() *qrLoginStruct {
	return &qrLoginStruct{
		Appid:     "549000912",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36 Edg/97.0.1072.55",
		Timeout:   time.Second * 10,
	}
}

// 请求客户端
func (login *qrLoginStruct) client() *http.Client {
	return &http.Client{
		Timeout: login.Timeout,
	}
}

// 获取扫码图片地址
func (login *qrLoginStruct) qrShowUrl() (url string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := r.Float32()
	url = fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=%s&e=2&l=M&s=3&d=72&v=4&t=%.16f&daid=8&pt_3rd_aid=0", login.Appid, t)
	return
}

// 获取扫码图片信息
func (login *qrLoginStruct) GetQr() (data getQrStruct, err error) {
	var (
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
	)
	if req, err = http.NewRequest("GET", login.qrShowUrl(), nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", login.UserAgent)
	if resp, err = login.client().Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("响应图片信息失败")
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	for _, v := range resp.Cookies() {
		if v.Name == "qrsig" {
			data.Qrsig = v.Value
		}
	}
	if data.Qrsig == "" {
		err = errors.New("获取qrsig失败")
		return
	}
	data.Image = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(bodyBytes))
	return
}

// 检测登录状态
func (login *qrLoginStruct) Check(qrsig string) (data loginStatusStruct, err error) {
	var (
		ptqrtoken = fmt.Sprintf("%d", login.ptqrtoken(qrsig))
		action    = fmt.Sprintf("%d", time.Now().UnixMilli())
		url       = "https://xui.ptlogin2.qq.com/ssl/ptqrlogin?u1=https%3A%2F%2Fcf.qq.com%2F&ptqrtoken=" + ptqrtoken + "&ptredirect=1&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-" + action + "&js_ver=22011714&js_type=1&login_sig=" + qrsig + "&pt_uistyle=40&aid=" + login.Appid + "&daid=8&"
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
		result    []string
	)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", login.UserAgent)
	req.AddCookie(&http.Cookie{Name: "qrsig", Value: qrsig})
	if resp, err = login.client().Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("响应登录状态失败")
		return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if result = regexp.MustCompile(`ptuiCB\('(.*)','(.*)','(.*)','(.*)','(.*)', '(.*)'\)`).FindStringSubmatch(string(bodyBytes)); result == nil {
		return data, errors.New("查找不到登录信息")
	}
	switch result[1] {
	case "0":
		var (
			req2  *http.Request
			resp2 *http.Response
		)
		data.Name, data.Url = result[6], result[3]
		qqReg := regexp.MustCompile(`&uin=([1-9][0-9]{4,9})&`).FindStringSubmatch(data.Url) // 从URL中提取QQ号
		if qqReg == nil {
			err = errors.New("获取QQ号失败")
			return
		}
		data.Uin = qqReg[1]
		if req2, err = http.NewRequest("GET", data.Url, nil); err != nil { // 获取 Cookie
			return
		}
		if resp2, err = http.DefaultTransport.RoundTrip(req2); err != nil {
			return data, err
		}
		defer resp2.Body.Close()
		for _, v := range resp2.Cookies() {
			switch {
			case v.Domain == "qq.com" && v.Name == "skey":
				data.Skey = v.Value
				fallthrough
			case v.Domain == "game.qq.com" && v.Name == "p_skey":
				data.P_skey = v.Value
				fallthrough
			case v.Domain == "game.qq.com" && v.Name == "pt4_token":
				data.Pt4_token = v.Value
			}
		}
		if len(data.Uin) > 10 || len(data.Skey) != 10 || len(data.P_skey) != 44 || len(data.Pt4_token) != 44 {
			err = errors.New("获取登录信息失败")
			return
		}
		data.Status, data.Message = 0, "登录成功"
	case "7":
		data.Status, data.Message = -1, "提交参数错误，请检查。"
	case "65":
		data.Status, data.Message = -1, "二维码已失效。"
	case "66":
		data.Status, data.Message = 1, "二维码未失效。"
	case "67":
		data.Status, data.Message = 2, "二维码认证中。"
	case "68":
		data.Status, data.Message = -1, "本次登录已被拒绝。"
	default:
		data.Status, data.Message = -2, result[5]
	}
	return
}

// ptqrtoken 计算
func (login *qrLoginStruct) ptqrtoken(qrsig string) int {
	len := len(qrsig)
	hash := 0
	for i := 0; i < len; i++ {
		hash += ((hash << 5 & 2147483647) + int(qrsig[i])) & 2147483647
		hash &= 2147483647
	}
	return hash & 2147483647
}
