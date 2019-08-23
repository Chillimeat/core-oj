package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/Myriad-Dreamin/core-oj/log"
)

const DriverName = "mysql"
const MasterDataSourceName = "coreoj-admin:123456@tcp(127.0.0.1:3306)/coreoj?charset=utf8"

type Code struct {
	ID        int `xorm:"not null pk autoincr"`
	Hash      []byte
	OwnedUID  int
	ProblemID int
	RuntimeID int
}

type Server struct {
	engine *xorm.Engine
	logger *log.TendermintLogger
}

func NewServer() (srv *Server) {
	srv = new(Server)
	return
}

func (srv *Server) prepareDatabase(driver, connection string) {

	srv.engine ,err := xorm.NewEngine(driver, connection)
	if err != nil {

	}	
}

func main() {
	var srv = NewServer()
	srv.prepareDatabase(DriverName, MasterDataSourceName)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("")

	r.Run()
}
