package env

import (
	"log"

	"github.com/spf13/viper"
)

type ENV struct{
	MONGODB_URI           string `mapstructure:"MONGODB_URI"`
}

func NewEnv() *ENV{
	env:=ENV{}
	viper.SetConfigFile(".env")
	err :=viper.ReadInConfig()
	if err!=nil{
		log.Fatal("Can't find the file .env", err)
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env

}