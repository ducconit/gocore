package config

import (
	"time"

	"github.com/spf13/viper"
)

var globalConfig Config

func SetGlobal(cfg Config) {
	globalConfig = cfg
}

func SetGlobalIfMissing(cfg Config) {
	if globalConfig == nil {
		globalConfig = cfg
	}
}

func Get(key string) any {
	return globalConfig.Get(key)
}

func GetString(key string) string {
	return globalConfig.GetString(key)
}

func GetInt(key string) int {
	return globalConfig.GetInt(key)
}

func GetBool(key string) bool {
	return globalConfig.GetBool(key)
}

func GetFloat64(key string) float64 {
	return globalConfig.GetFloat64(key)
}

func GetStringMap(key string) map[string]any {
	return globalConfig.GetStringMap(key)
}

func GetStringSlice(key string) []string {
	return globalConfig.GetStringSlice(key)
}

func GetIntSlice(key string) []int {
	return globalConfig.GetIntSlice(key)
}

func IsSet(key string) bool {
	return globalConfig.IsSet(key)
}

func AllSettings() map[string]any {
	return globalConfig.AllSettings()
}

func AllKeys() []string {
	return globalConfig.AllKeys()
}

func Set(key string, value any) {
	globalConfig.Set(key, value)
}

func SetDefault(key string, value any) {
	globalConfig.SetDefault(key, value)
}

func LoadFromFile(path string, options ...Option) error {
	return globalConfig.LoadFromFile(path, options...)
}

func LoadFromDB(db any, tableName string) error {
	return globalConfig.LoadFromDB(db, tableName)
}

func Reload() error {
	return globalConfig.Reload()
}

func Watch(key string, callback func(any)) {
	globalConfig.Watch(key, callback)
}
func GetDuration(key string) time.Duration {
	return globalConfig.GetDuration(key)
}

func GetTime(key string) time.Time {
	return globalConfig.GetTime(key)
}

func GetUint(key string) uint {
	return globalConfig.GetUint(key)
}

func GetUint32(key string) uint32 {
	return globalConfig.GetUint32(key)
}

func GetUint64(key string) uint64 {
	return globalConfig.GetUint64(key)
}

func GetInt32(key string) int32 {
	return globalConfig.GetInt32(key)
}

func GetInt64(key string) int64 {
	return globalConfig.GetInt64(key)
}

func GetSizeInBytes(key string) uint {
	return globalConfig.GetSizeInBytes(key)
}

func GetStringMapString(key string) map[string]string {
	return globalConfig.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return globalConfig.GetStringMapStringSlice(key)
}

func Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	return globalConfig.Unmarshal(rawVal, opts...)
}

func UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error {
	return globalConfig.UnmarshalKey(key, rawVal, opts...)
}
