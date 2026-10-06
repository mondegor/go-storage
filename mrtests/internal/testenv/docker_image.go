package testenv

import (
	"os"
)

// DockerImage - возвращает докер образ из переменной окружения envName,
// а если она не задана или пуста - образ по умолчанию defaultImage.
func DockerImage(envName, defaultImage string) string {
	if image := os.Getenv(envName); image != "" {
		return image
	}

	return defaultImage
}
