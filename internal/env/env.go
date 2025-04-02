package env

import "os"

func GetMqUser() string {
	return os.Getenv("MQ_USER")
}

func GetMqPassword() string {
	return os.Getenv("MQ_PASSWORD")
}
