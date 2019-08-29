package config

var (
	defaultConfiguration = &Configuration{
		DriverName:           "mysql",
		MasterDataSourceName: "coreoj-admin:123456@tcp(127.0.0.1:3306)/coreoj?charset=utf8",
		CodePath:             "/home/kamiyoru/data/test/",
		ProblemPath:          "/home/kamiyoru/data/problems/",
		KVPath:               "/home/kamiyoru/data/kv.db/",
	}
)

type Configuration struct {
	DriverName           string
	MasterDataSourceName string
	CodePath             string
	ProblemPath          string
	KVPath               string
}

func Config() *Configuration {
	return defaultConfiguration
}
