package types

import "time"

type kiloByte = int

// TestConfig makes configuration of profiling program
type TestConfig struct {
}

// TestCase makes configuration of testcase
type TestCase struct {
	CaseNumber int `json:"cn"`

	TestPath string `json:"tp"`

	OptionStream int `json:"ops"`

	InputPath string `json:"inp"`

	StdOutputPath string `json:"soup"`

	TimeLimit time.Duration `json:"tl"`

	MemoryLimit kiloByte `json:"ml"`

	SpecialJudge     bool   `json:"spj"`
	SpecialJudgePath string `json:"spjp"`
}
