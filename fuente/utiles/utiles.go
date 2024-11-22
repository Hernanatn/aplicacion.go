package utiles

import (
	"strings"
)

func Limpiar(s string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(s, "\r"), "\n"))
}
