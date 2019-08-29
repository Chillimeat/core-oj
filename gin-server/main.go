package main

import (
	"context"
	"fmt"
	"os"

	config "github.com/Myriad-Dreamin/core-oj/config"
	"github.com/Myriad-Dreamin/core-oj/log"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Server struct {
	engine *xorm.Engine
	kvdb   *kvorm.GobDB
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

func (srv *Server) prepareDatabase(kvpath, driver, connection string) error {
	var err error
	srv.kvdb, err = kvorm.NewGobLevelDB(kvpath)
	if err != nil {
		return err
	}

	srv.engine, err = xorm.NewEngine(driver, connection)
	if err != nil {
		srv.logger.Error("prepare failed", "error", err)
		return err
	}

	morm.RegisterEngine(srv.engine)
	kvorm.RegisterEngine(srv.kvdb)

	srv.engine.ShowSQL(true)
	return nil
}

func (srv *Server) Close() error {
	return srv.kvdb.Close()
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
	procStater, err := kvorm.NewProcStater()
	if err != nil {
		return err
	}

	judgeService, err := NewJudgeService(coder, problemer, procStater, srv.logger)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("close...")
		if err := judgeService.Close(); err != nil {
			srv.logger.Debug("close error", err)
		}
	}()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	codeRouter := r.Group("/code")
	{
		var codeService = NewCodeService(coder, procStater, srv.logger)
		codeRouter.GET("/:id", codeService.Get)
		codeRouter.GET("/:id/content", codeService.GetContent)
		codeRouter.GET("/:id/result", codeService.GetResult)
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
	err = srv.prepareDatabase(config.Config().KVPath, config.Config().DriverName, config.Config().MasterDataSourceName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := srv.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	if err = srv.Serve(":23336"); err != nil {
		fmt.Println(err)
		return
	}
}
