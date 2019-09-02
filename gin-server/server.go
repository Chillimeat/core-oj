package server

import (
	"context"
	"fmt"

	config "github.com/Myriad-Dreamin/core-oj/config"
	"github.com/Myriad-Dreamin/core-oj/log"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"

	// import driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Server struct {
	engine *xorm.Engine
	kvdb   *kvorm.GobDB
	logger log.TendermintLogger

	coder      MinimumCodeDatabaseProvider
	problemer  MinimumProblemDatabaseProvider
	procStater MinimumProcStateDatabaseProvider

	judgeService   *JudgeService
	codeService    interface{}
	problemService interface{}
}

func NewServer() (srv *Server, err error) {
	srv = new(Server)

	srv.logger, err = log.NewZapColorfulDevelopmentSugarLogger()
	if err != nil {
		return nil, err
	}

	return
}

func (srv *Server) SetProblemService(problemService interface{}) {
	srv.problemService = problemService
}

func (srv *Server) SetCodeService(codeService interface{}) {
	srv.codeService = codeService
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

func (srv *Server) Close() {
	if err := srv.judgeService.Close(); err != nil {
		srv.logger.Debug("close js error", err)
	}
	if err := srv.kvdb.Close(); err != nil {
		srv.logger.Debug("close kvdb error", err)
	}
	return
}

func (srv *Server) Serve(ctx context.Context) {
	srv.judgeService.ProcessAllCodes(ctx)
	return
}

func (srv *Server) DefaultCodeRouter(codeRouter *gin.RouterGroup) {
	if codeService, ok := (srv.codeService).(CodeServiceInterface); ok {
		codeRouter.GET("/:id", codeService.Get)
		codeRouter.GET("/:id/content", codeService.GetContent)
		codeRouter.GET("/:id/result", codeService.GetResult)
		codeRouter.POST("/postform", codeService.PostForm)
		// codeRouter.PUT("/:id/updateform-runtimeid", codeService.UpdateRuntimeID)
		codeRouter.DELETE("/:id", codeService.Delete)
	}
}

func (srv *Server) DefaultProblemRouter(problemRouter *gin.RouterGroup) {
	if problemService, ok := (srv.problemService).(ProblemServiceInterface); ok {
		problemRouter.GET("/:id", problemService.Get)
		problemRouter.POST("/postform", problemService.PostForm)
		// problemRouter.PUT("/:id/updateform-runtimeid", problemService.UpdateRuntimeID)
		problemRouter.DELETE("/:id", problemService.Delete)
	}
}

func (srv *Server) DefaultProblemFSRouter(problemFSRouter *gin.RouterGroup) {

	if problemService, ok := (srv.problemService).(ProblemServiceInterface); ok {
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

func (srv *Server) DefaultPrepare(cfg *config.Configuration) error {
	err := srv.prepareDatabase(config.Config().KVPath, config.Config().DriverName, config.Config().MasterDataSourceName)
	if err != nil {
		return err
	}

	srv.coder, err = morm.NewCoder()
	if err != nil {
		return err
	}
	srv.problemer, err = morm.NewProblemer()
	if err != nil {
		return err
	}
	srv.procStater, err = kvorm.NewProcStater()
	if err != nil {
		return err
	}

	srv.judgeService, err = NewJudgeService(srv.coder, srv.problemer, srv.procStater, srv.logger)
	if err != nil {
		return err
	}

	srv.codeService = NewCodeService(srv.coder, srv.procStater, srv.logger)
	srv.problemService = NewProblemService(srv.problemer, srv.logger)
	return nil
}

func Main() {
	var srv, err = NewServer()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = srv.DefaultPrepare(config.Config())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer srv.Close()

	r := gin.Default()

	codeRouter := r.Group("/code")
	srv.DefaultCodeRouter(codeRouter)
	problemRouter := r.Group("/problem")
	srv.DefaultProblemRouter(problemRouter)
	problemFSRouter := r.Group("/problemfs")
	srv.DefaultProblemFSRouter(problemFSRouter)

	ctx, cancel := context.WithCancel(context.Background())
	go srv.Serve(ctx)
	defer cancel()
	r.Run(":23336")
}
