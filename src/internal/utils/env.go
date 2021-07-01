package utils

import (
	"fmt"
	"os"
)

func CheckEnv(keys ...string) {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			panic(fmt.Sprintf("%s not found on environment!", key))
		}
	}
}
