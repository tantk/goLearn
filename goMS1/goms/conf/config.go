package conf

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/tkanos/gonfig"
)

//Configuration :
type Configuration struct {
	DbUser   string
	DbPW     string
	DbPort   string
	DbHost   string
	DbName   string
	RESTport string
}

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../")
)

//GetConfig :
func GetConfig(params ...string) Configuration {
	configuration := Configuration{}
	fileName := fmt.Sprintf(Root + "/conf/conf.json")
	gonfig.GetConf(fileName, &configuration)
	return configuration
}
