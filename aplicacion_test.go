package aplicacion_test

import (
	"fmt"
	"testing"

	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
)

/*
	"bufio"
	"github.com/hernanatn/aplicacion.go/fuente/aplicacion"
	"github.com/hernanatn/aplicacion.go/fuente/aplicacion/comando"
	"fmt"
	"os"
	"testing"
		"regexp"
*/

func TestCadenas(t *testing.T) {
	var hola cadena.Cadena = "Hola\n"
	var juan cadena.Cadena = "Juan y Pedro\n"
	var como cadena.Cadena = "¿Cómo están?"
	var chau cadena.Cadena = "Chau"

	var tengo cadena.Cadena = hola.Formatear(cadena.Negrita) +
		juan.Formatear(cadena.Italica) + como.Formatear(cadena.Invertida) +
		chau.Formatear(cadena.Coloreador(color.RojoFondo))
	var quiero cadena.Cadena = "" +
		"\033[1m" + hola + "\033[0m" +
		"\033[3m" + juan + "\033[0m" +
		"\033[7m" + como + "\033[0m" +
		"\033[41m" + chau + "\033[0m"
	t.Log(cadena.Cadena("EAEAPP").Fatal(fmt.Errorf("eaea")))
	t.Log(cadena.Cadena("EAEAPP").Advertencia(fmt.Errorf("eaea")))
	t.Log(cadena.Cadena("EAEAPP").Error(fmt.Errorf("eaea")))
	t.Log(cadena.Cadena("EAEAPP asd sad ").Exito())
	t.Log(cadena.Cadena("EAEAPP asd ad as").Ok())
	t.Log(cadena.Cadena("EAEAPP asd adsa asd").Sugerencia())
	t.Log(cadena.Cadena("EAEAPP asd adsa asd").Debug(fmt.Errorf("asd")))
	t.Log(cadena.Tabla([]string{"1234", "32223", "1231211212112235"}, [][]string{
		{"1", "1", "1", "1"},
		{"1", "1", "1"},
		{"1", "1"},
		{"1", "1", "1"},
		{"1", "1222222222222", "1"},
	}))
	if quiero != tengo {
		t.Fatalf("Tengo:\n %s, tenía que ser:\n %s", tengo, quiero)
	}
}
