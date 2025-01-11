package cadena

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hernanatn/aplicacion.go/consola/color"
)

// Alias de string con métodos de formato útiles para consolas.
//
// Uso:
//
//	c := Cadena("Una cadena")
//	c.Negrita().Italica()
type Cadena string

type TipoAlineado int

const (
	IZQUIERDA = iota
	DERECHA
	CENTRO
	JUSTIFICADO
)

type Formateador func(string) string

type OpcionesFormato struct {
	Color    color.ColorFuente
	Fondo    color.ColorFondo
	Alineado TipoAlineado
	Sangria  string
	Prefijo  string
	Sufijo   string
}

func TextoJustificado(s string, ancho int, opciones ...OpcionesFormato) string {
	var o OpcionesFormato
	if opciones != nil {
		o = opciones[0]
	}
	if len(s) == 0 {
		return ""
	}
	if ancho >= len(s) {
		return Colorear(o.Sangria+s+o.Sufijo, o.Color)
	}
	var chunks []string = make([]string, 0, (len(s)-1)/ancho+1)
	largoFila := 0
	inicioFila := 0
	for i := range s {
		if largoFila == ancho {
			idxPrimerCar := inicioFila
			idxUltimoCar := i
			primerCar := string(s[inicioFila])
			if primerCar == " " {
				idxPrimerCar++
				idxUltimoCar++
			}
			largoFila = 0
			inicioFila = idxUltimoCar
			quiebre := ""
			if !((string(s[i-1]) == " " && string(s[i]) == " ") || string(s[idxUltimoCar]) == " " || string(s[idxUltimoCar-1]) == " ") {
				quiebre = "-"
				idxUltimoCar--
				largoFila++
				inicioFila--
			}
			chunks = append(chunks, s[idxPrimerCar:idxUltimoCar]+quiebre)
		}
		largoFila++
	}
	chunks = append(chunks, s[inicioFila:])
	return strings.ReplaceAll(
		Colorear(o.Sangria+strings.Join(chunks, fmt.Sprintf("%s\n%s", o.Sufijo, o.Prefijo))+o.Sufijo, o.Color),
		" -",
		"  ",
	)
	//}
}

func Negrita(s string) string {
	return fmt.Sprintf("\033[1m%s%s", s, color.Resetear)
}
func Italica(s string) string {
	return fmt.Sprintf("\033[3m%s%s", s, color.Resetear)
}
func Subrayada(s string) string {
	return fmt.Sprintf("\033[4m%s%s", s, color.Resetear)
}
func Invertida(s string) string {
	return fmt.Sprintf("\033[7m%s%s", s, color.Resetear)
}
func Colorear(s string, c color.Color) string {
	return fmt.Sprintf("%s%s%s", c, s, color.Resetear)
}

func Coloreador(c color.Color) func(string) string {
	f := func(s string) string {
		return Colorear(s, c)
	}
	return f
}

func Limpiar(s string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(s, "\r"), "\n"))
}

func (c Cadena) Negrita() Cadena {
	return Cadena(Negrita(string(c)))
}
func (c Cadena) Italica() Cadena {
	return Cadena(Italica(string(c)))
}
func (c Cadena) Subrayada() Cadena {
	return Cadena(Subrayada(string(c)))
}
func (c Cadena) Invertida() Cadena {
	return Cadena(Invertida(string(c)))
}
func (c Cadena) Colorear(col color.Color) Cadena {
	return Cadena(Colorear(string(c), col))
}

func (c Cadena) Limpiar() Cadena {
	return Cadena(strings.TrimSpace(strings.Trim(strings.Trim(c.S(), "\r"), "\n")))

}

func (c Cadena) Formatear(formatos ...Formateador) Cadena {
	var cad Cadena = c
	for _, f := range formatos {
		cad = Cadena(f(cad.S()))
	}
	return cad
}

func Sugerencia(msg string) string {
	return Italica(Colorear(fmt.Sprintf("%s.", msg), color.GrisFuente)) + "\n"

}
func Debug(msg string, err error) string {
	return Colorear(fmt.Sprintf(Negrita("[DEBUG]")+"\t%s. err: %v.", msg, err), color.GrisFuente) + "\n"
}
func Ok(msg string) string {
	return Colorear(fmt.Sprintf("%s.", msg), color.VerdeFuente) + "\n"
}
func Exito(msg string) string {
	return Negrita(Colorear(fmt.Sprintf("%s.", msg), color.VerdeFondo)) + "\n"
}
func Advertencia(msg string, err error) string {
	return Colorear(fmt.Sprintf(Negrita("[ADVERTENCIA]")+"\t%s. err: %v.", msg, err), color.AmarilloFuente) + "\n"
}
func Error(msg string, err error) string {
	return Colorear(Negrita("[ERROR]"), color.RojoFuente) + Colorear(fmt.Sprintf("\t%s. err: %v.", msg, err), color.RojoFuente) + "\n"
}
func Fatal(msg string, err error) string {
	return Colorear(fmt.Sprintf(Negrita("[FATAL]")+"\t%s. err: %v.", msg, err), color.RojoFondo) + "\n"
}

func ImprimirTitulo(s string) {
	fmt.Println(Cadena(s).Negrita().Colorear(color.CyanFuente))
}
func Titulo(s string) Cadena {
	return Cadena(s).Negrita().Colorear(color.CyanFuente) + "\n"
}

func ImprimirSubtitulo(s string) {
	fmt.Println(Cadena(s).Italica().Colorear(color.GrisFuente))
}
func Subtitulo(s string) Cadena {
	return Cadena(s).Italica().Negrita().Colorear(color.GrisFuente) + "\n"
}
func Señalador(s string) Cadena {
	return Cadena(s).Colorear(color.CyanFondo) + " "
}

func (c Cadena) Imprimir(f *bufio.Writer) {
	if f == nil {
		f = bufio.NewWriter(os.Stdout)
	}

	f.WriteString(c.S())
	f.Flush()
}

func (c Cadena) String() string {
	return string(c)
}

func (c Cadena) S() string {
	return c.String()
}

func DesdeArchivo(nombre string) (Cadena, error) {
	data, err := os.ReadFile(nombre)
	if err != nil {
		return "", err
	}
	return Cadena(string(data)), nil
}
