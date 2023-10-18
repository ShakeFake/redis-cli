package random

import "github.com/google/uuid"

func GetRandomUUID() string {
	return uuid.New().String()
}
