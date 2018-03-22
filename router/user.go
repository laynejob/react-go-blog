package router

import (
    "github.com/gin-gonic/gin"
    "github.com/laynejob/react-go-blog/db"
    "strconv"
    "net/http"
)

func getUser(c *gin.Context) {
    param := c.Query("id")
    id, err := strconv.ParseUint(param, 10, 64)
    if err != nil {
        c.JSON(http.StatusNotAcceptable, gin.H{"err": err.Error()})
    }
    user, err := db.DB.UserTable.SelectById(id)
    if user != nil {
        //c.JSON(200, gin.H{"id": user.Id, "email": user.Email, "name": user.Username,})
        c.JSON(200, user)
    } else {
        c.JSON(http.StatusOK, gin.H{})
    }
}