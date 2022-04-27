# 例子

## 字符串处理

### 字符串

```t
写法参考
 不要参数:
  会处理: {#name}
  不处理: {#name()} {#name(myName)}
 要参数:
  会处理: {#number(1,2,3)}
  不处理: {#number} {#number()} {#number(,)} {#number(1,)} {#number(1,,3)} 注意参数值不要为空
正则表达式匹配参考（优先处理）
 名字: {#name666}
 网址: {#urltoolv}
 多个: {#0123456789abcABC}
 复制: {#display(这是我要显示的信息)}
```

### 处理后

```t
写法参考
 不要参数:
  会处理: 相思
  不处理: {#name()} {#name(myName)}
 要参数:
  会处理: 获取到的参数是 [1 2 3] 长度 3
  不处理: {#number} {#number()} {#number(,)} {#number(1,)} {#number(1,,3)} 注意参数值不要为空
正则表达式匹配参考（优先处理）
 名字: 获取到的参数是 [666] 长度 1
 网址: https://www.toolv.cn
 多个: 数字是 0123456789 字母是 abc 大写是 ABC
 复制: 获取到的参数是 [display] 长度 1 值是[这是我要显示的信息] 长度1
```

## 代码示例

```go
package main

import (
    "fmt"

    "github.com/toolvcn/toolv/strReplacer"
)

func main() {
    str := `
写法参考
 不要参数:
  会处理: {#name}
  不处理: {#name()} {#name(myName)}
 要参数:
  会处理: {#number(1,2,3)}
  不处理: {#number} {#number()} {#number(,)} {#number(1,)} {#number(1,,3)} 注意参数值不要为空
正则表达式匹配参考（优先处理）
 名字: {#name666}
 网址: {#urltoolv}
 多个: {#0123456789abcABC}
 复制: {#display(这是我要显示的信息)}
`
    r := strReplacer.Default()
    r.AddParams("name", func(args ...string) string {
        return "相思"
    }, false)
    r.AddParams("number", func(args ...string) string {
        return fmt.Sprintf("获取到的参数是 %v 长度 %d", args, len(args))
    }, true)
    r.AddRegexParams(`name(\d+)`, func(params []string, args ...string) string {
        return fmt.Sprintf("获取到的参数是 %v 长度 %d", params, len(params))
    }, false)
    r.AddRegexParams(`url([^\s]+)`, func(params []string, args ...string) string {
        return "https://www." + params[0] + ".cn"
    }, false)
    r.AddRegexParams(`([0-9]+)([a-z]+)([A-Z]+)`, func(params []string, args ...string) string {
        return "数字是 " + params[0] + " 字母是 " + params[1] + " 大写是 " + params[2]
    }, false)
    r.AddRegexParams(`(display)`, func(params []string, args ...string) string {
        return fmt.Sprintf("获取到的参数是 %v 长度 %d 值是%v 长度%d", params, len(params), args, len(args))
    }, true)
    newStr := r.String(str)
    fmt.Printf("oldStr: %v\nnewStr: %v\n", str, newStr)
}
```
