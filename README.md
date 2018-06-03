# Form [![GoDoc](https://godoc.org/github.com/go-mego/form?status.svg)](https://godoc.org/github.com/go-mego/form) [![Coverage Status](https://coveralls.io/repos/github/go-mego/form/badge.svg?branch=master)](https://coveralls.io/github/go-mego/form?branch=master) [![Build Status](https://travis-ci.org/go-mego/form.svg?branch=master)](https://travis-ci.org/go-mego/form) [![Go Report Card](https://goreportcard.com/badge/github.com/go-mego/form)](https://goreportcard.com/report/github.com/go-mego/form)

Form 套件用以從 `application/x-www-form-urlencoded` 或 `multipart/form-data` 表單請求中讀取欄位資料，並且能將其表單映射至本地建構體。

# 索引

* [安裝方式](#安裝方式)
* [使用方式](#使用方式)
    * [取得欄位](#取得欄位)
    * [取得欄位切片](#取得欄位切片)
    * [映射表單](#映射表單)
        * [自訂欄位](#自訂欄位)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/go-mego/form
```

# 使用方式

將 `form.New` 傳入 Mego 所提供的 `Use` 來將 Form 中介軟體安插進全域中介軟體並開始使用。

```go
package main

import (
	"github.com/go-mego/form"
	"github.com/go-mego/mego"
)

func main() {
	m := mego.New()
	// 將表單中介軟體作為全域中介軟體就可以在所有路由中使用。
	m.Use(form.New())
	m.Run()
}
```

Form 中介軟體也能夠套用到單一路由並與其他路由區隔驗證邏輯。

```go
func main() {
	m := mego.New()
	// 表單中介軟體可以獨立用於某路由。
	m.POST("/", form.New(), func(f *form.Form) {
		// ...
	})
	m.Run()
}
```

## 取得欄位

透過 `form.Get` 可以從請求表單中取得特定的欄位資料，當該欄位為空、不存在時則會取得空白字串。

```go
func main() {
	m := mego.New()
	m.POST("/", form.New(), func(f *form.Form) {
        // 輸出表單欄位中 `username` 的資料。
		fmt.Println(f.Get("username"))
	})
	m.Run()
}
```

## 取得欄位切片

`form.GetMulti` 能夠以字串切片的方式取得表單中的同個多重欄位（亦即重複欄位），切片的長度取決於表單中有多少筆重複欄位。

```go
func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		// 透過 `GetMulti` 能以字串切片的方式來取得表單中的重複欄位。
		fmt.Println(f.GetMulti("photos"))
	})
	m.Run()
}
```

## 映射表單

`form.Bind` 可以將表單資料映射至本地的建構體。當無法映射至建構體的時候會離開請求，並且以 `text/plain` 回傳一個 HTTP 400 錯誤狀態碼。

```go
type User struct {
	Username string
	Password string
}

func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		var u User
		// 透過 `Bind` 能夠將接收到的表單資料映射至本地的建構體變數。
		err := f.Bind(&u)
		if err != nil {
			// ...
		}
		fmt.Println(u.Username)
	})
	m.Run()
}
```

映射時會忽略欄位的分隔符號（`-`、`_`）與大小寫，這讓你不需要額外處理名稱不同的問題。

```
createdAt       -> CreatedAt
user-id         -> UserID
favorite_photos -> FavoritePhotos
```

### 自訂欄位

有些時候請求的欄位可能與本地建構體不符，這時可以在建構體中使用 `form` 標籤來標明該建構體欄位對應請求表單中的哪個欄位。

```go
type User struct {
	Username  string `form:"id"`
	CreatedAt string `form:"registration_date"`
}

func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		var u User
        f.Bind(&u)
        // ...
	})
	m.Run()
}
```