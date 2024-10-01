package assets

import (
	"embed"
)

//go:embed *.png
var Icons embed.FS

func GetIcon(name string) ([]byte, error) {
	return Icons.ReadFile(name)
}
