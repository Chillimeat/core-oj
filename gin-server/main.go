package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Myriad-Dreamin/core-oj/log"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const DriverName = "mysql"
const MasterDataSourceName = "coreoj-admin:123456@tcp(127.0.0.1:3306)/coreoj?charset=utf8"

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

func insertCode() {

}

const (
	CodeOK int = iota
	CodeNotFound
	CodeInsertError
	CodeCodeTypeMissing
	CodeCodeProblemIDMissing
	CodeCodeBodyMissing
	CodeCodeUploaded
)

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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	codeRoute := r.Group("/code")

	codeRoute.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		fmt.Println("id", id)
		code, err := new(morm.Coder).Query(int(id))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if code != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":      CodeOK,
				"hash":      code.Hash,
				"owneruid":  code.OwnedUID,
				"problemid": code.ProblemID,
				"runtimeID": code.RuntimeID,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeNotFound,
			})
		}
	})

	codeRoute.POST("/post", func(c *gin.Context) {
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		code := new(morm.Code)
		var ok bool
		if code.CodeType, ok = c.GetPostForm("type"); !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeCodeTypeMissing,
			})
			return
		}
		var problemID string
		if problemID, ok = c.GetPostForm("problemid"); !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeCodeProblemIDMissing,
			})
			return
		}
		var problemIDx int64
		problemIDx, err = strconv.ParseInt(problemID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeCodeProblemIDMissing,
			})
			return
		}
		code.ProblemID = int(problemIDx)
		// todo: find problemid

		var body string
		if body, ok = c.GetPostForm("body"); !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeCodeBodyMissing,
			})
			return
		}

		codeHash := md5.New()

		buf := bytes.NewBufferString(body)
		var p = make([]byte, 0)
		_, err := io.TeeReader(buf, codeHash).Read(p)
		fmt.Println(p)

		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		code.Hash = codeHash.Sum(nil)

		if cx, err := new(morm.Coder).QueryHash(code.Hash); err != nil {
			c.AbortWithError(500, err)
			return
		} else if cx != nil {
			c.JSON(200, gin.H{
				"code": CodeCodeUploaded,
			})
			return
		}
		code.RuntimeID = 0

		affected, err := code.Insert()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if affected != 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":      CodeOK,
				"id":        code.ID,
				"hash":      code.Hash,
				"owneruid":  code.OwnedUID,
				"problemid": code.ProblemID,
				"runtimeID": code.RuntimeID,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
		}
	})

	r.Run()
}
