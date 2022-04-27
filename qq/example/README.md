# 例子

```go
package main

import (
    "fmt"
    "time"

    "github.com/toolvcn/toolv/qq"
)

func main() {
    q := qq.NewQrLogin()
    // 1.获取 Qr 图片 信息
    qr, err := q.GetQr()
    if err != nil { // 获取二维码失败错误信息
        panic(err)
    }
    fmt.Printf("%v\n", qr.Image)
    // 你大概会得到像这样的结果
    // data:image/png;base64,iVBOR...
    // 复制这条结果，然后到：https://www.it399.com/image/base64 进行还原
    // 2.获取二维码状态
    for {
        time.Sleep(time.Second * 5)
        qrStatus, err := q.LoginStatus(qr.Qrsig)
        if err != nil { // 获取二维码状态失败错误信息
            panic(err)
        }
        fmt.Printf("%+v\n", qrStatus)
        if qrStatus.Status == 0 {
            fmt.Println("登录成功")
            break
        }
    }
}
```
