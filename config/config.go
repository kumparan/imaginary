package config

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil && Env() != "test" {
		log.Warningf("%v", err)
	}
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// RedisCacheHost :nodoc:
func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

// RedisLockHost :nodoc:
func RedisLockHost() string {
	return viper.GetString("redis.lock_host")
}

// DisableCaching :nodoc:
func DisableCaching() bool {
	return viper.GetBool("disable_caching")
}

// CacheTTL :nodoc:
func CacheTTL() time.Duration {
	if !viper.IsSet("cache_ttl") {
		return DefaultCacheTTL
	}

	return time.Duration(viper.GetInt("cache_ttl")) * time.Millisecond
}

// AWSRegion :nodoc:
func AWSRegion() string {
	return viper.GetString("aws.region")
}

// AWSS3Bucket :nodoc:
func AWSS3Bucket() string {
	return viper.GetString("aws.s3_bucket")
}

// AWSS3Key :nodoc:
func AWSS3Key() string {
	return viper.GetString("aws.s3_key")
}

// AWSS3Secret :nodoc:
func AWSS3Secret() string {
	return viper.GetString("aws.s3_secret")
}
