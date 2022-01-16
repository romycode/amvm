package env

import "os"

func Get(name, fallback string) string {
	if env, ok := os.LookupEnv(name); ok {
		return env
	}

	return fallback
}
