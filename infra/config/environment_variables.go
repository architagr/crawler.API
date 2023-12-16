package config

import "fmt"

func GetStringWithEnv(env, str string) string {
	return fmt.Sprintf("%s-%s", env, str)
}
