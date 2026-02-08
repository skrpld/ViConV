package config

import (
	"log"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/skrpld/NearBeee/internal/database/mongodb"
	"github.com/skrpld/NearBeee/internal/database/postgres"
	"github.com/skrpld/NearBeee/internal/logger"
	"github.com/skrpld/NearBeee/internal/transport/grpc/servers"
	"github.com/skrpld/NearBeee/pkg/utils/jwt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	servers.NearBeeeServerConfig `mapstructure:",squash"`
	mongodb.MongoDBConfig        `mapstructure:",squash"`
	postgres.PostgresConfig      `mapstructure:",squash"`
	logger.LoggerConfig          `mapstructure:",squash"`
	jwt.JWTConfig                `mapstructure:",squash"`
}

var (
	currentConfig  atomic.Pointer[Config]
	lastReloadTime atomic.Int64
)

func InitConfig() error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := reloadConfig(); err != nil {
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		now := time.Now().UnixNano()
		last := lastReloadTime.Load()

		if now-last < int64(500*time.Millisecond) {
			return
		}

		time.Sleep(100 * time.Millisecond) //TODO: чекнуть мб ремув
		if err := reloadConfig(); err != nil {
			log.Printf("Error reading config file, %s", err)
			return
		}

		lastReloadTime.Store(time.Now().UnixNano())

		log.Println("Configuration updated")
	})
	viper.WatchConfig()

	return nil
}

func reloadConfig() error {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	bindStructDefaults(v, Config{})

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

func GetConfig() *Config {
	return currentConfig.Load()
}

func bindStructDefaults(v *viper.Viper, s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Anonymous || field.Tag.Get("mapstructure") == ",squash" {
			bindStructDefaults(v, val.Field(i).Interface())
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
