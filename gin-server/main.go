package main

import (
	"fmt"

	"github.com/Myriad-Dreamin/core-oj/log"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const DriverName = "mysql"
const MasterDataSourceName = "coreoj-admin:123456@tcp(127.0.0.1:3306)/coreoj?charset=utf8"
const codepath = "/home/kamiyoru/data/test/"

type Server struct {
	engine *xorm.Engine
	logger log.TendermintLogger
}

func NewServer() (srv *Server, err error) {
	srv = new(Server)

	srv.logger, err = log.NewZapColorfulDevelopmentSugarLogger()
	if err != nil {
		return nil, err
	}

	return
}

func (srv *Server) prepareDatabase(driver, connection string) error {
	var err error
	srv.engine, err = xorm.NewEngine(driver, connection)
	if err != nil {
		srv.logger.Error("prepare failed", "error", err)
		return err
	}

	morm.RegisterEngine(srv.engine)

	srv.engine.ShowSQL(true)
	return nil
}

func (srv *Server) Serve(port string) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	codeRouter := r.Group("/code")
	{
		var cs = NewCodeService()
		codeRouter.GET("/:id", cs.Get)
		codeRouter.GET("/:id/content", cs.GetContent)
		codeRouter.POST("/postform", cs.PostForm)
		codeRouter.PUT("/:id/updateform-runtimeid", cs.UpdateRuntimeID)
		codeRouter.DELETE("/:id", cs.Delete)
	}

	return r.Run(port)
}

func main() {
	var srv, err = NewServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.prepareDatabase(DriverName, MasterDataSourceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = srv.Serve(":23336"); err != nil {
		fmt.Println(err)
		return
	}
}
