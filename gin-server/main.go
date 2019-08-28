package main

import (
	"context"
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
const problempath = "/home/kamiyoru/data/problems/"

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

	coder, err := morm.NewCoder()
	if err != nil {
		return err
	}
	problemer, err := morm.NewProblemer()
	if err != nil {
		return err
	}

	judgeService, err := NewJudgeService(coder, problemer, srv.logger)
	if err != nil {
		return err
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	codeRouter := r.Group("/code")
	{
		var codeService = NewCodeService(coder, srv.logger)
		codeRouter.GET("/:id", codeService.Get)
		codeRouter.GET("/:id/content", codeService.GetContent)
		codeRouter.POST("/postform", codeService.PostForm)
		// codeRouter.PUT("/:id/updateform-runtimeid", codeService.UpdateRuntimeID)
		codeRouter.DELETE("/:id", codeService.Delete)
	}

	var problemService = NewProblemService(problemer, srv.logger)
	{
		problemRouter := r.Group("/problem")
		{

			problemRouter.GET("/:id", problemService.Get)
			problemRouter.POST("/postform", problemService.PostForm)
			// problemRouter.PUT("/:id/updateform-runtimeid", problemService.UpdateRuntimeID)
			problemRouter.DELETE("/:id", problemService.Delete)

		}

		problemFSRouter := r.Group("/problemfs")
		{
			problemFSRouter.GET("/:id/stat", problemService.Stat)
			problemFSRouter.PUT("/:id/mkdir", problemService.Mkdir)
			problemFSRouter.GET("/:id/ls", problemService.Ls)
			problemFSRouter.GET("/:id/read", problemService.Read)
			problemFSRouter.POST("/:id/write", problemService.Write)
			problemFSRouter.POST("/:id/writes", problemService.Writes)
			problemFSRouter.POST("/:id/zip", problemService.Zip)
			problemFSRouter.GET("/:id/config", problemService.ReadConfigV2)
			problemFSRouter.PUT("/:id/config", problemService.PutConfig)

		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go judgeService.ProcessAllCodes(ctx)
	defer cancel()
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
