package server

import (
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"
)

type MinimumCodeDatabaseProvider interface {
	Query(int) (*morm.Code, error)
	QueryHash([]byte) (*morm.Code, error)
	PushTask(*morm.Code) error
	StartToExecuteTask(*morm.Code) error
	SettleTask(int) (bool, error)
	ExposeWaitingCodes() chan *morm.Code
	ExposeRunningCodes() chan *morm.Code
}

type MinimumProcStateDatabaseProvider interface {
	Query(int) ([]kvorm.ProcState, error)
	InsertP(int, []*kvorm.ProcState) error
}

type CodeServiceInterface interface {
	Get(*gin.Context)
	GetContent(*gin.Context)
	GetResult(*gin.Context)
	PostForm(*gin.Context)
	Delete(*gin.Context)
}

type ProblemServiceInterface interface {
	Get(*gin.Context)
	PostForm(*gin.Context)
	Delete(*gin.Context)

	Stat(*gin.Context)
	Mkdir(*gin.Context)
	Ls(*gin.Context)
	Read(*gin.Context)
	Write(*gin.Context)
	Writes(*gin.Context)
	Zip(*gin.Context)
	ReadConfigV2(*gin.Context)
	PutConfig(*gin.Context)
}
