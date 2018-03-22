package db

import (
    "github.com/laynejob/react-go-blog/conf"
    "log"
)

var (
    DB          *mysqlDB
    //OAuthConfig *oauth2.Config
    //
    //StorageBucket     *storage.BucketHandle
    //StorageBucketName string
    //
    //SessionStore sessions.Store
    //
    //PubsubClient *pubsub.Client

    // Force import of mgo library.
    //_ mgo.Session
)

const (
    timeFormart = "2006-01-02 15:04:05"
)

func init()  {
    var err error
    c := conf.GetConf()
    DB, err = newMySQLDB(MySQLConfig{
        Username: c.Db.User,
        Password: c.Db.Password,
        Host:     c.Db.Host,
        Port:     c.Db.Port,
        Database: c.Db.Database,
    })

    if err != nil {
        log.Fatal(err)
    }
}
