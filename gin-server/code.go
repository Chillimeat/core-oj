package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"

	config "github.com/Myriad-Dreamin/core-oj/config"
	"github.com/Myriad-Dreamin/core-oj/log"
	types "github.com/Myriad-Dreamin/core-oj/types"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	"github.com/gin-gonic/gin"
)

// CodeService defines handler functions of code router
type CodeService struct {
	Coder      *morm.Coder
	ProcStater *kvorm.ProcStater
	logger     log.TendermintLogger
	codePath   string
}

// NewCodeService return a pointer of CodeService
func NewCodeService(coder *morm.Coder, procStater *kvorm.ProcStater, logger log.TendermintLogger) *CodeService {
	return &CodeService{
		Coder:      coder,
		ProcStater: procStater,
		logger:     logger,
		codePath:   config.Config().CodePath,
	}
}

// Delete codes from database
func (cr *CodeService) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	code, err := cr.Coder.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	affected, err := code.Delete()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"id":        code.ID,
			"hash":      code.Hash,
			"owneruid":  code.OwnerUID,
			"problemid": code.ProblemID,
			"status":    code.Status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeDeleteFailed,
		})
	}

}

// Get codes from database
func (cr *CodeService) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	code, err := cr.Coder.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if code != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"hash":      code.Hash,
			"owneruid":  code.OwnerUID,
			"problemid": code.ProblemID,
			"status":    code.Status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeNotFound,
		})
	}
}

// GetContent codes from database with content
func (cr *CodeService) GetContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	code, err := cr.Coder.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if code != nil {

		// f, err := os.Open()
		// if err != nil {
		// 	c.AbortWithError(500, err)
		// 	return
		// }
		// b, err := ioutil.ReadAll(f)
		// f.Close()
		// if err != nil {
		// 	c.AbortWithError(500, err)
		// 	return
		// }

		// c.JSON(http.StatusOK, gin.H{
		// 	"code": CodeOK,
		// 	"hash": code.Hash,
		// 	// todo: hack b to string
		// 	"content":   string(b),
		// 	"owneruid":  code.OwnerUID,
		// 	"problemid": code.ProblemID,
		// 	"status":    code.Status,
		// })
		c.File(cr.codePath + hex.EncodeToString(code.Hash) + "/main.cpp")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeNotFound,
		})
	}
}

// GetResult from kvdb
func (cr *CodeService) GetResult(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	code, err := cr.Coder.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if code == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeNotFound,
		})
		return
	}

	result, err := cr.ProcStater.Query(code.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   CodeOK,
		"result": result,
	})
}

// PostForm codes to database
func (cr *CodeService) PostForm(c *gin.Context) {
	code := new(morm.Code)
	var ok bool

	// rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
	var codeType string
	if codeType, ok = c.GetPostForm("type"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeTypeMissing,
		})
		return
	}

	if code.CodeType, ok = morm.CodeTypeMap[codeType]; !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeTypeUnknown,
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
	problemIDx, err := strconv.ParseInt(problemID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeProblemIDMissing,
		})
		return
	}
	code.ProblemID = int(problemIDx)
	// todo: find problemid

	var ownerUID string
	if ownerUID, ok = c.GetPostForm("owneruid"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeOwnerUIDMissing,
		})
		return
	}
	var ownerUIDx int64
	ownerUIDx, err = strconv.ParseInt(ownerUID, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeOwnerUIDMissing,
		})
		return
	}
	code.OwnerUID = int(ownerUIDx)
	// todo: find problemid

	var body string
	if body, ok = c.GetPostForm("body"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeBodyMissing,
		})
		return
	}

	codeHash := md5.New()
	_, err = codeHash.Write([]byte(body))

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	code.Hash = codeHash.Sum(nil)

	if cx, err := cr.Coder.QueryHash(code.Hash); err != nil {
		c.AbortWithError(500, err)
		return
	} else if cx != nil {
		c.JSON(200, gin.H{
			"code": CodeCodeUploaded,
			"hash": code.Hash,
		})
		return
	}

	var path = cr.codePath + hex.EncodeToString(code.Hash)
	if _, err = os.Stat(path); err != nil && !os.IsExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		err = os.Chmod(path, 0777)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}
	path += "/main.cpp"
	if _, err = os.Stat(path); err != nil && !os.IsExist(err) {
		f, err := os.Create(path)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		_, err = f.WriteString(body)
		f.Close()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}
	err = nil

	code.Status = types.StatusWaitingForJudge

	affected, err := code.Insert()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected != 0 {
		cr.Coder.PushTask(code)
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"id":        code.ID,
			"hash":      code.Hash,
			"owneruid":  code.OwnerUID,
			"problemid": code.ProblemID,
			"status":    code.Status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeInsertFailed,
		})
	}
}
