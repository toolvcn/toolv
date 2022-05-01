// 字符串自定义参数替换
package strReplacer

import (
	"regexp"
	"strings"
)

type (
	// 替换配置
	Replace struct {
		MatchStart  string                        // 匹配开始标记
		MatchEnd    string                        // 匹配结束标记
		ParamsStart string                        // 参数开始标记
		ParamsSplit string                        // 参数分隔符号
		ParamsEnd   string                        // 参数结束标记
		Params      map[string]ReplaceParams      // 普通参数列表
		RegexParams map[string]ReplaceRegexParams // 正则参数列表
	}
	// 普通参数结构体
	ReplaceParams struct {
		Args    bool                        // 是否有参数
		Handler func(args ...string) string // 参数处理函数
	}
	// 正则参数结构体
	ReplaceRegexParams struct {
		Args    bool                                         // 是否有参数
		Handler func(params []string, args ...string) string // 参数处理函数
	}
)

// 默认 Replace 对象
//	 return &Replace{
//	 	MatchStart:  "{#",
//	 	MatchEnd:    "}",
//	 	ParamsStart: `\(`,
//	 	ParamsSplit: ",",
//	 	ParamsEnd:   `\)`,
//	 	Params:      make(map[string]ReplaceParams),
//	 	RegexParams: make(map[string]ReplaceRegexParams),
//	 }
func Default() *Replace {
	return &Replace{
		MatchStart:  "{#",
		MatchEnd:    "}",
		ParamsStart: `\(`,
		ParamsSplit: ",",
		ParamsEnd:   `\)`,
		Params:      make(map[string]ReplaceParams),
		RegexParams: make(map[string]ReplaceRegexParams),
	}
}

// 新建 Replace 对象
//	 return &Replace{
//	 	MatchStart:  "",
//	 	MatchEnd:    "",
//	 	ParamsStart: "",
//	 	ParamsSplit: "",
//	 	ParamsEnd:   "",
//	 	Params:      make(map[string]ReplaceParams),
//	 	RegexParams: make(map[string]ReplaceRegexParams),
//	 }
func New() *Replace {
	return &Replace{
		MatchStart:  "",
		MatchEnd:    "",
		ParamsStart: "",
		ParamsSplit: "",
		ParamsEnd:   "",
		Params:      make(map[string]ReplaceParams),
		RegexParams: make(map[string]ReplaceRegexParams),
	}
}

// 添加普通参数
//	name    - 参数名
//	handler - 参数处理函数
//	args    - 是否有参数
func (r *Replace) AddParams(name string, handler func(...string) string, args bool) {
	if r.Params == nil {
		r.Params = make(map[string]ReplaceParams)
	}
	r.Params[name] = ReplaceParams{Args: args, Handler: handler}
}

// 添加正则参数
//	name    - 参数名是正则表达式
//	handler - 参数处理函数
//	args    - 是否有参数
func (r *Replace) AddRegexParams(name string, handler func([]string, ...string) string, args bool) {
	if r.RegexParams == nil {
		r.RegexParams = make(map[string]ReplaceRegexParams)
	}
	r.RegexParams[name] = ReplaceRegexParams{Args: args, Handler: handler}
}

// 删除普通参数
//	name - 参数名
func (r *Replace) DelParams(name string) {
	delete(r.Params, name)
}

// 删除正则参数
//	name - 参数名是正则表达式
func (r *Replace) DelRegexParams(name string) {
	delete(r.RegexParams, name)
}

// 返回替换后的字符串
//	s - 要替换的字符串
func (r *Replace) String(s string) string {
	return r.replace(s)
}

// 替换字符串
//	s - 要替换的字符串
func (r *Replace) ToString(s *string) {
	*s = r.replace(*s)
}

// replace - 替换
//	s - 要替换的字符串
func (r *Replace) replace(s string) string {
	reg := regexp.MustCompile(r.MatchStart + ".+?" + r.MatchEnd)
	return reg.ReplaceAllStringFunc(s, r.replaceMatch)
}

// 替换函数
//	s - 要替换的字符串
func (r *Replace) replaceMatch(s string) string {
	params, args := r.parseParams(s)
	if params == "" { // 没有参数
		return s
	}
	// 从正则参数中获取
	for k, v := range r.RegexParams {
		res := regexp.MustCompile(k).FindStringSubmatch(params)
		switch {
		case res == nil: // 没有匹配到
			continue
		case v.Args && len(args) > 0: // 有参数,且参数数量大于0
			return v.Handler(res[1:], args...)
		case !v.Args && len(args) == 0: // 没有参数,且参数数量为0
			return v.Handler(res[1:])
		default:
			return s
		}
	}
	// 从普通参数中获取
	for k, v := range r.Params {
		switch {
		case k != params: // 没有匹配到
			continue
		case v.Args && len(args) > 0: // 有参数,且参数数量大于0
			return v.Handler(args...)
		case !v.Args && len(args) == 0: // 没有参数,且参数数量为0
			return v.Handler()
		default:
			return s
		}
	}
	return s
}

// 解析参数
//	s      - 要解析的字符串
//	params - 参数名
//	args   - 参数值
func (r *Replace) parseParams(s string) (params string, args []string) {
	reg := regexp.MustCompile(`^` + r.MatchStart + `([^` + r.ParamsStart + `]+)` + `(?:` + r.ParamsStart + `(.+?)` + r.ParamsEnd + `)?` + r.MatchEnd + `$`)
	res := reg.FindStringSubmatch(s)
	switch {
	case res == nil: // 没有匹配到
		return
	case len(res) >= 2: // 匹配到参数名
		params = res[1]
		fallthrough
	case len(res) == 3: // 匹配到参数值
		s := strings.Split(res[2], r.ParamsSplit)
		for _, v := range s { // 确保参数值不为空
			if v == "" {
				return
			}
		}
		args = s
	}
	return
}
