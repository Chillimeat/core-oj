package language

const (
	InternalCodeTypeCpp11 int64 = iota
	InternalCodeTypeCpp14
	InternalCodeTypeCpp17
	InternalCodeTypeC
	InternalCodeTypePython2
	InternalCodeTypePython3
	InternalCodeTypeJava
)

var CodeTypeMap map[string]int64
var ReverseCodeTypeMap map[int64]string
