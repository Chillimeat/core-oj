package language

import (
	argiexecutor "github.com/Myriad-Dreamin/core-oj/language/argi-executor"
	"github.com/Myriad-Dreamin/core-oj/language/cpcompiler"
	iexecutor "github.com/Myriad-Dreamin/core-oj/language/i-executor"
	"github.com/Myriad-Dreamin/core-oj/language/icompiler"
	iocompiler "github.com/Myriad-Dreamin/core-oj/language/iocompiler"
	pureexecutor "github.com/Myriad-Dreamin/core-oj/language/pure-executor"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

var Compilers map[string]types.Compiler
var ReverseCompilers map[int64]types.Compiler

var Executors map[string]types.Executor
var ReverseExecutors map[int64]types.Executor

type Config struct {
	SourcePath string
	TargetPath string
}

var Configs map[string]*Config
var ReverseConfigs map[int64]*Config

const (
	CompilerTypeIOCompiler int = iota
	CompilerTypeICompiler
	CompilerTypeCopyCompiler
)

const (
	ExecutorTypePureExecutor int = iota
	ExecutorTypeIExecutor
	ExecutorTypeArgIExecutor
)

func constructCompiler(cpType int, path, src, dst string, args []string) types.Compiler {
	switch cpType {
	case CompilerTypeIOCompiler:
		return &iocompiler.Compiler{
			Path:   path,
			Args:   args,
			Source: src,
			Target: dst,
		}
	case CompilerTypeICompiler:
		return &icompiler.Compiler{
			Path:   path,
			Args:   args,
			Source: src,
		}
	case CompilerTypeCopyCompiler:
		return &cpcompiler.Compiler{
			Source: src,
			Target: dst,
		}
	default:
		panic("error type of compiler")
	}
}

func constructExecutor(cpType int, path string, dest string, args []string) types.Executor {
	switch cpType {
	case ExecutorTypeIExecutor:
		return &iexecutor.Executor{
			Path:     path,
			DestName: dest,
		}
	case ExecutorTypePureExecutor:
		return &pureexecutor.Executor{DestName: dest}
	case ExecutorTypeArgIExecutor:
		return &argiexecutor.Executor{
			Path:     path,
			DestName: dest,
			Regp:     args[0],
		}
	default:
		panic("error type of compiler")
	}
}

func init() {
	Compilers = make(map[string]types.Compiler)
	ReverseCompilers = make(map[int64]types.Compiler)
	Executors = make(map[string]types.Executor)
	ReverseExecutors = make(map[int64]types.Executor)
	Configs = make(map[string]*Config)
	ReverseConfigs = make(map[int64]*Config)
	CodeTypeMap = make(map[string]int64)
	ReverseCodeTypeMap = make(map[int64]string)

	for _, config := range []struct {
		TypeID                                             int64
		TypeString                                         string
		CompilerType, ExecutorType                         int
		SourcePath, TargetPath, CompilerPath, ExecutorPath string
		CompilerArgs, ExecutorArgs                         []string
	}{
		{
			InternalCodeTypeCpp11, "c++11",
			CompilerTypeIOCompiler, ExecutorTypePureExecutor,
			"/main.cpp", "/main", "g++", "",
			[]string{"-std=c++11", "-fmax-errors=5"}, nil,
		},
		{
			InternalCodeTypeCpp14, "c++14",
			CompilerTypeIOCompiler, ExecutorTypePureExecutor,
			"/main.cpp", "/main", "g++", "",
			[]string{"-std=c++14", "-fmax-errors=5"}, nil,
		},
		{
			InternalCodeTypeCpp17, "c++17",
			CompilerTypeIOCompiler, ExecutorTypePureExecutor,
			"/main.cpp", "/main", "g++", "",
			[]string{"-std=c++17", "-fmax-errors=5"}, nil,
		},
		{
			InternalCodeTypeC, "c",
			CompilerTypeIOCompiler, ExecutorTypePureExecutor,
			"/main.c", "/main", "g++", "",
			[]string{"-fmax-errors=5"}, nil,
		},
		{
			InternalCodeTypePython2, "python2",
			CompilerTypeCopyCompiler, ExecutorTypeIExecutor,
			"/main.py", "/main", "", "python2", nil, nil,
		},
		{
			InternalCodeTypePython3, "python3",
			CompilerTypeCopyCompiler, ExecutorTypeIExecutor,
			"/main.py", "/main", "", "python3", nil, nil,
		},
		{
			InternalCodeTypeJava, "java",
			CompilerTypeICompiler, ExecutorTypeIExecutor,
			"/Main.java", "Main.Class", "javac", "java",
			nil, []string{"-cp {w} {d}"},
		},
	} {
		CodeTypeMap[config.TypeString] = config.TypeID
		ReverseCodeTypeMap[config.TypeID] = config.TypeString
		Compilers[config.TypeString] = constructCompiler(config.CompilerType, config.CompilerPath, config.SourcePath, config.TargetPath, config.CompilerArgs)
		Executors[config.TypeString] = constructExecutor(config.ExecutorType, config.ExecutorPath, config.TargetPath, config.ExecutorArgs)
		Configs[config.TypeString] = &Config{
			SourcePath: config.SourcePath,
			TargetPath: config.TargetPath,
		}
		ReverseConfigs[config.TypeID] = Configs[config.TypeString]
		ReverseCompilers[config.TypeID] = Compilers[config.TypeString]
	}
}
