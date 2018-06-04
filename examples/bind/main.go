package main

import (
	"github.com/go-mego/form"
	"github.com/go-mego/mego"
)

type User struct {
	Username string
	Password string
}

func main() {
	e := mego.Default()
	e.POST("/", form.New(), func(c *mego.Context, f *form.Form) {
		var u User
		err := f.Bind(&u)
		if err != nil {
			panic(err)
		}
		c.String(200, "%+v", u)
	})
	e.Run()
}
