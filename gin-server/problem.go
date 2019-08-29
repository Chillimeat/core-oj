package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	config "github.com/Myriad-Dreamin/core-oj/config"
	"github.com/Myriad-Dreamin/core-oj/log"
	problemconfig "github.com/Myriad-Dreamin/core-oj/problem-config"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"

	"github.com/gin-gonic/gin"
)

// ProblemService defines handler functions of problem router
type ProblemService struct {
	Problemer   *morm.Problemer
	logger      log.TendermintLogger
	problemPath string
}

// NewProblemService return a pointer of ProblemService
func NewProblemService(problemer *morm.Problemer, logger log.TendermintLogger) *ProblemService {
	return &ProblemService{
		Problemer:   problemer,
		logger:      logger,
		problemPath: config.Config().ProblemPath,
	}
}

// Delete problems from database
func (pr *ProblemService) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	problem, err := pr.Problemer.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	affected, err := problem.Delete()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"id":        problem.ID,
			"name":      problem.Name,
			"owner_uid": problem.OwnerUID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeDeleteFailed,
		})
	}

}

// Get problems from database
func (pr *ProblemService) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	problem, err := pr.Problemer.Query(int(id))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if problem != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"name":      problem.Name,
			"owner_uid": problem.OwnerUID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeNotFound,
		})
	}
}

// PostForm problems to database
func (pr *ProblemService) PostForm(c *gin.Context) {
	problem := new(morm.Problem)
	var ok bool

	if problem.Name, ok = c.GetPostForm("name"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemNameMissing,
		})
		return
	}

	var ownerUID string
	if ownerUID, ok = c.GetPostForm("owneruid"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeOwnerUIDMissing,
		})
		return
	}
	var err error
	var ownerUIDx int64
	ownerUIDx, err = strconv.ParseInt(ownerUID, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeCodeOwnerUIDMissing,
		})
		return
	}
	problem.OwnerUID = int(ownerUIDx)

	affected, err := problem.Insert()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":      CodeOK,
			"id":        problem.ID,
			"name":      problem.Name,
			"owner_uid": problem.OwnerUID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeInsertFailed,
		})
	}
}

type Adaptor struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	ModTime time.Time `json:"modtime"`
}

func adaptToJson(stat os.FileInfo) *Adaptor {
	return &Adaptor{
		Name:    stat.Name(),
		Size:    stat.Size(),
		IsDir:   stat.IsDir(),
		ModTime: stat.ModTime(),
	}
}

func (pr *ProblemService) Stat(c *gin.Context) {
	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}

	path = pr.problemPath + c.Param("id") + path
	var stat os.FileInfo
	var err error
	if stat, err = os.Stat(path); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeStatError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   CodeOK,
		"status": adaptToJson(stat),
	})
}

func (pr *ProblemService) Mkdir(c *gin.Context) {
	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}

	path = pr.problemPath + c.Param("id") + path
	if err := os.Mkdir(path, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": CodeOK,
	})
}

func (pr *ProblemService) Ls(c *gin.Context) {
	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}

	path = pr.problemPath + c.Param("id") + path
	var files []os.FileInfo
	var err error
	if files, err = ioutil.ReadDir(path); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": CodeOK,
		"result": func() (ret []*Adaptor) {
			ret = make([]*Adaptor, 0, len(files))
			for _, stat := range files {
				ret = append(ret, adaptToJson(stat))
			}
			return
		}(),
	})
}

func (pr *ProblemService) Read(c *gin.Context) {

	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}
	c.File(pr.problemPath + c.Param("id") + path)
}

func (pr *ProblemService) ReadConfig(c *gin.Context) {
	path := pr.problemPath + c.Param("id") + "/problem-config"
	configPath, err := problemconfig.Select(path+".json", path+".yml", path+".toml", path+".xml")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	c.File(configPath)
}

func (pr *ProblemService) ReadConfigV2(c *gin.Context) {
	path := pr.problemPath + c.Param("id") + "/problem-config"
	var config problemconfig.ProblemConfig
	err := problemconfig.Load(&config, path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   CodeOK,
		"config": &config,
	})
}

func (pr *ProblemService) PutConfig(c *gin.Context) {
	path := pr.problemPath + c.Param("id") + "/problem-config"
	var config = new(problemconfig.ProblemConfig)
	err := problemconfig.Load(config, path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	var ok bool
	var key, value string
	if key, ok = c.GetPostForm("key"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeConfigKeyMissing,
		})
		return
	}
	if value, ok = c.GetPostForm("value"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeConfigValueMissing,
		})
		return
	}
	err = config.Modify(key, json.RawMessage(value))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeConfigModifyError,
			"err":  err.Error(),
		})
		return
	}
	err = problemconfig.Save(config, path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   CodeOK,
		"config": &config,
	})
}

func (pr *ProblemService) Write(c *gin.Context) {

	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	if err = c.SaveUploadedFile(file, pr.problemPath+c.Param("id")+path+file.Filename); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": CodeOK,
	})
}

func (pr *ProblemService) Writes(c *gin.Context) {

	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	files := form.File["upload"]
	path = pr.problemPath + c.Param("id") + path
	for _, file := range files {
		if err = c.SaveUploadedFile(file, path+file.Filename); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeFSExecError,
				"err":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": CodeOK,
	})
}

func (pr *ProblemService) Zip(c *gin.Context) {

	var (
		path string
		ok   bool
	)
	if path, ok = c.GetPostForm("path"); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeProblemPathMissing,
		})
		return
	}
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}
	path = pr.problemPath + c.Param("id") + path
	zipName := path + file.Filename
	if err = c.SaveUploadedFile(file, zipName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	r, err := zip.OpenReader(zipName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CodeFSExecError,
			"err":  err.Error(),
		})
		return
	}

	var release = func() {
		r.Close()
		err := os.Remove(zipName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	for _, file := range r.File {
		rc, err := file.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeFSExecError,
				"err":  err.Error(),
			})
			release()
			return
		}
		filename := path + file.Name
		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeFSExecError,
				"err":  err.Error(),
			})
			rc.Close()
			release()
			return
		}
		w, err := os.Create(filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeFSExecError,
				"err":  err.Error(),
			})
			rc.Close()
			release()
			return
		}
		_, err = io.Copy(w, rc)
		rc.Close()
		w.Close()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": CodeFSExecError,
				"err":  err.Error(),
			})
			release()
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"code": CodeOK,
	})
	release()
}
