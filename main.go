package main

import (
    "github.com/gin-gonic/gin"
)


func main() {
	
	db.AutoMigrate(&User{})
    r := gin.Default()
    r.POST("/users", CreateUser)
    r.GET("/users", ListUsers)
    r.GET("/users/:id", GetUser)
    r.PUT("/users/:id", UpdateUser)
    r.DELETE("/users/:id", DeleteUser)
    r.Run(":8080")
}
