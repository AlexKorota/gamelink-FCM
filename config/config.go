package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

var (
	//NatsDialAddress - dial address for nats
	NatsDialAddress string
	//NatsFirebaseChan NATS Chan for firebase push
	NatsFirebaseChan string
	//ServiceKeyPAth - path to json service google key
	ServiceKeyPath string
)

const (
	modeKey             = "MODE"
	devMode             = "development"
	natsDial            = "NATSDIAL"
	natsFirebaseChannel = "NATSCHANFIREBASE"
	servicekeypath      = "SKEYPATH"
)

func init() {
	LoadEnvironment()
}

//GetEnvironment - this function returns mode string of the os environment or "deployment.yaml" mode if empty or not defined
func GetEnvironment() string {
	var env string
	if env = os.Getenv(modeKey); env == "" {
		return devMode
	}
	return env
}

//IsDevelopmentEnv - this function try to get mode environment and check it is deployment.yaml
func IsDevelopmentEnv() bool { return GetEnvironment() == devMode }

func LoadEnvironment() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = godotenv.Load(path.Join(wd, strings.ToLower(GetEnvironment())+".env"))
	if err != nil {
		log.Warning(err.Error())
	}
	NatsDialAddress = os.Getenv(natsDial)
	if NatsDialAddress == "" {
		log.Fatal("nats dial address must be set")
	}
	NatsFirebaseChan = os.Getenv(natsFirebaseChannel)
	if NatsFirebaseChan == "" {
		log.Fatal("nats firebase chan must be set")
	}
	ServiceKeyPath = os.Getenv(servicekeypath)
	if ServiceKeyPath == "" {
		log.Fatal("path to service key must be set")
	}
}
