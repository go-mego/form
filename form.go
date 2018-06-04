package form

import (
	"net/http"
	"strings"

	"github.com/go-mego/binding"
	"github.com/go-mego/mego"
)

const (
	// defaultMemory 是預設的表單解析記憶體大小，當表單使用超過此記憶體大小時會回傳錯誤。
	defaultMemory = 32 << 20 // 32 MB
	// mimeMultipartForm 是基本表單的請求 MIME 種類。
	mimeMultipartForm = "multipart/form-data"
	// 欄位標籤名稱。
	fieldTag = "form"
)

// New 能夠建立一個新的表單模組建構體並且分析請求中的傳統表單內容。
func New() mego.HandlerFunc {
	return func(c *mego.Context) {
		c.Map(&Form{
			context: c,
		})
	}
}

// Form 呈現了一個表單模組建構體。
type Form struct {
	// context 是 Mego 引擎的上下文建構體。
	context *mego.Context
}

// Bind 能夠將接收到的表單請求資料映射至本地的變數。
// 當無法映射至建構體的時候會離開請求，並且以 `text/plain` 回傳一個 HTTP 400 錯誤狀態碼。
func (f *Form) Bind(dest interface{}) error {
	err := f.ShouldBind(dest)
	if err != nil {
		f.context.AbortWithError(http.StatusBadRequest, err)
		return err
	}
	return nil
}

// ShouldBind 和 `Bind` 相同，
// 但會在映射失敗的時候不做任何處理（亦即：不會離開請求、不會回傳錯誤狀態碼）。
func (f *Form) ShouldBind(dest interface{}) error {
	if strings.HasPrefix(f.context.ContentType(), mimeMultipartForm) {
		err := f.context.Request.ParseMultipartForm(defaultMemory)
		if err != nil {
			return err
		}
	} else {
		err := f.context.Request.ParseForm()
		if err != nil {
			return err
		}
	}
	err := binding.Bind(dest, f.context.Request.Form, fieldTag)
	if err != nil {
		return err
	}
	return nil
}

// Has 能夠得知一個指定的欄位是否存在於表單中。
func (f *Form) Has(key string) bool {
	v := f.GetMulti(key)
	return len(v) != 0 && len(v[0]) > 0
}

// Get 能從請求中取得的指定欄位的內容，這個請求通常是 POST 的 `urlencoded` 或 `multipart` 表單。
// 當該欄位不存在的時候會回傳一個空白字串。
func (f *Form) Get(key string) string {
	if !f.Has(key) {
		return ""
	}
	return f.GetMulti(key)[0]
}

// GetDefault 會從請求中取得的指定欄位的內容，
// 這和 `Get` 基本相同，但當無指定參數時會以預設值（`defaultValue`）作為回傳結果。
func (f *Form) GetDefault(key, defaultValue string) string {
	v := f.Get(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// GetMulti 能夠回傳一組基於指定欄位的字串切片作為回應。
// 切片的長度基於請求有多少個相同欄位而定。
func (f *Form) GetMulti(key string) []string {
	r := f.context.Request
	err := r.ParseForm()
	if err != nil {
		f.context.AbortWithError(http.StatusBadRequest, err)
		return []string{}
	}
	err = r.ParseMultipartForm(defaultMemory)
	if err != nil {
		f.context.AbortWithError(http.StatusBadRequest, err)
		return []string{}
	}
	if v := r.PostForm[key]; len(v) > 0 {
		return v
	}
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if v := r.MultipartForm.Value[key]; len(v) > 0 {
			return v
		}
	}
	return []string{}
}
