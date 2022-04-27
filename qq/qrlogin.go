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
func (login *qrLoginStruct) qrShowUrl() string {
	rand.Seed(time.Now().UnixNano())
	t := rand.Float32()
	url := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=%s&e=2&l=M&s=3&d=72&v=4&t=%.16f&daid=8&pt_3rd_aid=0", login.Appid, t)
	return url
}

// 获取扫码图片信息
func (login *qrLoginStruct) GetQr() (getQrStruct, error) {
	var data getQrStruct
	req, err := http.NewRequest("GET", login.qrShowUrl(), nil)
	if err != nil {
		return data, err
	}
	req.Header.Set("User-Agent", login.UserAgent)
	resp, err := login.client().Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return data, errors.New("响应图片信息失败")
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	for _, v := range resp.Cookies() {
		if v.Name == "qrsig" {
			data.Qrsig = v.Value
		}
	}
	if data.Qrsig == "" {
		return data, errors.New("获取qrsig失败")
	}
	data.Image = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(respByte))
	return data, nil
}

// 检测登录状态
func (login *qrLoginStruct) LoginStatus(qrsig string) (loginStatusStruct, error) {
	var data loginStatusStruct
	ptqrtoken := fmt.Sprintf("%d", login.ptqrtoken(qrsig))
	action := fmt.Sprintf("%d", time.Now().UnixMilli())
	url := "https://xui.ptlogin2.qq.com/ssl/ptqrlogin?u1=https%3A%2F%2Fcf.qq.com%2F&ptqrtoken=" + ptqrtoken + "&ptredirect=1&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-" + action + "&js_ver=22011714&js_type=1&login_sig=" + qrsig + "&pt_uistyle=40&aid=" + login.Appid + "&daid=8&"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Set("User-Agent", login.UserAgent)
	req.AddCookie(&http.Cookie{Name: "qrsig", Value: qrsig})
	resp, err := login.client().Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return data, errors.New("响应登录信息失败")
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	reg := regexp.MustCompile(`ptuiCB\('(.*)','(.*)','(.*)','(.*)','(.*)', '(.*)'\)`)
	result := reg.FindStringSubmatch(string(respByte))
	if result == nil {
		return data, errors.New("查找不到登录信息")
	}
	switch result[1] {
	case "0":
		data.Name = result[6]
		data.Url = result[3]
		// 从URL中提取QQ号
		reg2 := regexp.MustCompile(`&uin=([1-9][0-9]{4,9})&`)
		qqReg := reg2.FindStringSubmatch(data.Url)
		if qqReg == nil {
			return data, errors.New("获取QQ号失败")
		}
		data.Uin = qqReg[1]
		// 获取 Cookie
		req2, err := http.NewRequest("GET", data.Url, nil)
		if err != nil {
			return data, err
		}
		resp2, err := http.DefaultTransport.RoundTrip(req2)
		if err != nil {
			return data, err
		}
		defer resp2.Body.Close()
		for _, v := range resp2.Cookies() {
			if v.Name == "skey" && v.Domain == "qq.com" {
				data.Skey = v.Value
			}
			if v.Domain == "game.qq.com" {
				if v.Name == "p_skey" {
					data.P_skey = v.Value
				}
				if v.Name == "pt4_token" {
					data.Pt4_token = v.Value
				}
			}
		}
		if len(data.Uin) > 10 || len(data.Skey) != 10 || len(data.P_skey) != 44 || len(data.Pt4_token) != 44 {
			return data, errors.New("获取登录信息失败")
		}
		data.Status = 0
		data.Message = "登录成功"
	case "7":
		data.Status = -1
		data.Message = "提交参数错误，请检查。"
	case "65":
		data.Status = -1
		data.Message = "二维码已失效。"
	case "66":
		data.Status = 1
		data.Message = "二维码未失效。"
	case "67":
		data.Status = 2
		data.Message = "二维码认证中。"
	case "68":
		data.Status = -1
		data.Message = "本次登录已被拒绝。"
	default:
		data.Status = -2
		data.Message = result[5]
	}
	return data, nil
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
