package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"time"
)


type kiloByte = int

type TaskConfig struct {
	Name        string        `json:"name" yaml:"name" toml:"name" xml:"name"`
	CaseCount   int           `json:"case-count" yaml:"case-count" toml:"case-count" xml:"case-count"`
	TimeLimit   time.Duration `json:"time-limit" yaml:"time-limit" toml:"time-limit" xml:"time-limit"`
	MemoryLimit kiloByte      `json:"memory-limit" yaml:"memory-limit" toml:"memory-limit" xml:"memory-limit"`
	InputPath   string        `json:"input-path" yaml:"input-path" toml:"input-path" xml:"input-path"`
	OutputPath  string        `json:"output-path" yaml:"output-path" toml:"output-path" xml:"output-path"`
}

type JudgeConfig struct {
	Type  string       `json:"judge-type" yaml:"judge-type" toml:"judge-type" xml:"judge-type"`
	Tasks []TaskConfig `json:"tasks" yaml:"tasks" toml:"tasks" xml:"tasks"`
}

type SpecialJudgeConfig struct {
	SpecialJudge bool   `json:"enable" yaml:"enable" toml:"enable" xml:"enable"`
	LanguageType string `json:"language-type" yaml:"language-type" toml:"language-type" xml:"language-type"`
	FilePath     string `json:"file-path" yaml:"file-path" toml:"file-path" xml:"file-path"`
}

type ProblemConfig struct {
	LoadType           string             `json:"-"`
	Name               xml.Name           `json:"-" xml:"Problem-Config"`
	JudgeConfig        JudgeConfig        `json:"judge" yaml:"judge" toml:"judge" xml:"judge"`
	SpecialJudgeConfig SpecialJudgeConfig `json:"special-judge" yaml:"special-judge" toml:"special-judge" xml:"special-judge"`
}

func (c *ProblemConfig) Modify(path string, val json.RawMessage) error {
	if path == "" {
		err := json.Unmarshal(val, c)
		if err != nil {
			return err
		}
		return nil
	}
	paths := strings.Split(path, ".")
	switch paths[0] {
	case "judge":
		return c.JudgeConfig.Modifys(paths[1:], val)
	case "special-judge":
		return c.SpecialJudgeConfig.Modifys(paths[1:], val)
	default:
		return errors.New("property missing")
	}
}

func (c *ProblemConfig) Modifys(paths []string, val json.RawMessage) error {
	if len(paths) == 0 {
		err := json.Unmarshal(val, c)
		if err != nil {
			return err
		}
		return nil
	}
	switch paths[0] {
	case "judge-type":
		return c.JudgeConfig.Modifys(paths[1:], val)
	case "special-judge":
		return c.SpecialJudgeConfig.Modifys(paths[1:], val)
	default:
		return errors.New("property missing")
	}
}

func (c *JudgeConfig) Modifys(paths []string, val json.RawMessage) error {
	if len(paths) == 0 {
		err := json.Unmarshal(val, c)
		if err != nil {
			return err
		}
		return nil
	}
	switch paths[0] {
	case "judge-type":
		err := json.Unmarshal(val, &c.Type)
		if err != nil {
			return err
		}
		return nil
	case "tasks":
		if len(paths) < 2 {
			return errors.New("path consumed")
		}
		index, err := strconv.Atoi(paths[1])
		if err != nil {
			return err
		}
		if len(c.Tasks) <= index {
			return errors.New("index overflow")
		}
		return c.Tasks[index].Modifys(paths[2:], val)
	default:
		return errors.New("property missing")
	}
}

func (c *SpecialJudgeConfig) Modifys(paths []string, val json.RawMessage) error {
	if len(paths) == 0 {
		err := json.Unmarshal(val, c)
		if err != nil {
			return err
		}
		return nil
	}
	switch paths[0] {
	case "enable":
		err := json.Unmarshal(val, &c.SpecialJudge)
		if err != nil {
			return err
		}
		return nil
	case "language-type":
		err := json.Unmarshal(val, &c.LanguageType)
		if err != nil {
			return err
		}
		return nil
	case "file-path":
		err := json.Unmarshal(val, &c.FilePath)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("property missing")
	}
}

func (c *TaskConfig) Modifys(paths []string, val json.RawMessage) error {
	if len(paths) == 0 {
		err := json.Unmarshal(val, c)
		if err != nil {
			return err
		}
		return nil
	}
	switch paths[0] {
	case "name":
		err := json.Unmarshal(val, &c.Name)
		if err != nil {
			return err
		}
		return nil
	case "time-limit":
		err := json.Unmarshal(val, &c.TimeLimit)
		if err != nil {
			return err
		}
		return nil
	case "memory-limit":
		err := json.Unmarshal(val, &c.MemoryLimit)
		if err != nil {
			return err
		}
		return nil
	case "input-path":
		err := json.Unmarshal(val, &c.InputPath)
		if err != nil {
			return err
		}
		return nil
	case "output-path":
		err := json.Unmarshal(val, &c.OutputPath)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("property missing")
	}
}
