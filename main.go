package main

import (

	"github.com/laynejob/react-go-blog/router"
	"github.com/laynejob/react-go-blog/conf"
	"fmt"
	"github.com/golang/glog"
)


func main() {
	c := conf.GetConf()
	fmt.Println(c)
	r := router.SetupRouter()

	glog.Warning("warning")
	glog.Error("err")
	glog.Info("info")
	glog.V(2).Infoln("This line will be printed if you use -v=N with N >= 2.")
	// Listen and Server in 0.0.0.0:8080
	r.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))
}
