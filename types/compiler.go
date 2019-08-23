package types

type Compiler interface {
	Compile(string, string) error
}
