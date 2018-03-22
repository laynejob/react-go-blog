package router

import (

    "github.com/golang/glog"
    "github.com/gin-gonic/gin"
    "github.com/laynejob/react-go-blog/conf"
    "github.com/laynejob/react-go-blog/oauth2/google"
    "github.com/szuecs/gin-glog"
    "time"
    "github.com/zalando/gin-oauth2"
    "net/http"
)

func SetupRouter() *gin.Engine {
    c := conf.GetConf()
    // Disable Console Color
    // gin.DisableConsoleColor()
    //gin.SetMode(gin.ReleaseMode)
    //r0 := gin.Default()
    r0 := gin.New()
    r0.Use(ginglog.Logger(3 * time.Second))
    r := r0.Group(c.Prefix)

    // static file folder
    //r.Static("static", "./fe/dist/")
    // the index file position
    r.StaticFile("/", "./fe/dist/index.html")


    glog.Infoln(c)

    // group router
    v1 := r.Group("/v1")
    a := r.Group("/auth")
    a.Use(ginoauth2.RequestLogger([]string{"uid"}, "data"))
    redirectURL := "http://127.0.0.1:8080/"
    credFile := "./oauth2/google1.json"
    scopes := []string{
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/userinfo.profile",
    }
    secret := []byte("secret")
    sessionName := "goquestsession"
    google.Setup(redirectURL, credFile, scopes, secret)
    r0.Use(google.Session(sessionName))
    r0.GET("/login", google.LoginHandler)

    a.Use(google.Auth())
    //a.GET("/", UserInfoHandler)
    a.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": "yes"})
    })
    a.GET("/api", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello from private for groups"})
    })

    v1.GET("/page", func(c *gin.Context) {
        c.String(200, "v1 page")
    })


    // Get user value
    r.GET("/user", getUser)
    //func(c *gin.Context) {
    //    user := c.Params.ByName("name")
    //    value, ok := DB[user]
    //    if ok {
    //        c.JSON(200, gin.H{"user": user, "value": value})
    //    } else {
    //        c.JSON(200, gin.H{"user": user, "status": "no value"})
    //    }
    //}

    a.GET("/page", func(c *gin.Context) {
        //glog.Level.Set("4")

        glog.Info("info")
        glog.Warning("warning")
        glog.Error("error")
        c.String(200, "admin page")
    })
    a.GET("/comment", func(c *gin.Context) {
        c.String(200,"admin comment")
    })

    // Ping test
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })



    // Authorized group (uses gin.BasicAuth() middleware)
    // Same than:
    // authorized := r.Group("/")
    // authorized.Use(gin.BasicAuth(gin.Credentials{
    //      "foo":  "bar",
    //      "manu": "123",
    //}))
    authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
        "foo":  "bar", // user:foo password:bar
        "manu": "123", // user:manu password:123
    }))

    authorized.POST("admin", func(c *gin.Context) {
        //user := c.MustGet(gin.AuthUserKey).(string)

        // Parse JSON
        var json struct {
            Value string `json:"value" binding:"required"`
        }

        if c.Bind(&json) == nil {
            //DB[user] = json.Value
            c.JSON(200, gin.H{"status": "ok"})
        }
    })

    r0.Use(gin.Recovery())
    return r0
}

func UserInfoHandler(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": ctx.MustGet("user").(google.User)})
}