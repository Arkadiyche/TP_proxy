package models

import (
	"github.com/Arkadiyche/TP_proxy/utils"
	"github.com/jackc/pgx"
)

type Config struct {
	Port string
	Database pgx.ConnPoolConfig
}
var ServerConfig = Config{
	Port:     ":8080",
	Database: pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Host:                 "localhost",
			Port:                 5432,
			Database:             "secure",
			User:                 "secure",
			Password:             "secure",
			TLSConfig:            nil,
			UseFallbackTLS:       false,
			FallbackTLSConfig:    nil,
			Logger:               nil,
			LogLevel:             0,
			Dial:                 nil,
			RuntimeParams:        nil,
			OnNotice:             nil,
			CustomConnInfo:       nil,
			CustomCancel:         nil,
			PreferSimpleProtocol: false,
			TargetSessionAttrs:   "",
		},
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	},
}

var Params = utils.GetParams()