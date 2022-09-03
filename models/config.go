package models

import (
	"fmt"
)

var Config Configuration

type DatabaseConf struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
}

func (d *DatabaseConf) URI() string {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
	return uri
}

type APIConf struct {
	JWTSecretKey    string `mapstructure:"jwt_secret_key"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	ErrorWebhookUrl string `mapstructure:"error_webhook_url"`
	FrontendUrl     string `mapstructure:"frontend_url"`
	HopToken        string `mapstructure:"hop_token"`
	HopProjectID    string `mapstructure:"hop_project_id"`
}

type Configuration struct {
	API      *APIConf      `mapstructure:"api"`
	Database *DatabaseConf `mapstructure:"database"`
}
