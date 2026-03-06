package config

import (
	"log"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/skrpld/NearBeee/internal/core/database/mongodb"
	"github.com/skrpld/NearBeee/internal/core/database/postgres"
	"github.com/skrpld/NearBeee/internal/core/logger"
	"github.com/skrpld/NearBeee/internal/transport/rest/servers"
	"github.com/skrpld/NearBeee/pkg/utils/jwt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	servers.HttpServerConfig `mapstructure:",squash"`
	mongodb.MongoDBConfig    `mapstructure:",squash"`
	postgres.PostgresConfig  `mapstructure:",squash"`
	logger.LoggerConfig      `mapstructure:",squash"`
	jwt.JWTConfig            `mapstructure:",squash"`
}

var (
	currentConfig atomic.Pointer[Config]
	timer         *time.Timer
)

func InitConfig() error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := reloadConfig(); err != nil {
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(200*time.Millisecond, func() {
			if err := reloadConfig(); err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Println("Configuration updated")
		})
	})
	viper.WatchConfig()

	return nil
}

func GetConfig() *Config {
	return currentConfig.Load()
}

func reloadConfig() error {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	setStructDefaults(v, Config{})

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var newCfg Config
	if err := v.Unmarshal(&newCfg); err != nil {
		return err
	}

	currentConfig.Store(&newCfg)

	jwt.UpdateJWTConfig(&newCfg.JWTConfig)

	return nil
}

func setStructDefaults(v *viper.Viper, s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Anonymous || field.Tag.Get("mapstructure") == ",squash" {
			setStructDefaults(v, val.Field(i).Interface())
			continue
		}

		defaultVal := field.Tag.Get("env-default")
		if defaultVal != "" {
			key := field.Tag.Get("mapstructure")
			if key == "" {
				key = strings.ToUpper(field.Name)
			}

			v.SetDefault(key, defaultVal)
		}
	}
}
