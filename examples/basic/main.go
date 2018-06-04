package main

import (
	"github.com/go-mego/form"
	"github.com/go-mego/mego"
)

func main() {
	e := mego.Default()
	e.POST("/", form.New(), func(c *mego.Context, f *form.Form) {
		c.String(200, "%s", f.Get("username"))
	})
	e.Run()
}
