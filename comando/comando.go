package comando

import (
	"errors"
	"slices"
	"strings"

	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
	"github.com/hernanatn/aplicacion.go/utiles"
)

type Consola = consola.Consola
type Cadena = consola.Cadena

type Opciones = consola.Opciones
type Parametros = consola.Parametros
type Argumentos = []any

type Accion = func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error)

type CodigoError int

const (
	EXITO CodigoError = iota << 0
	ERROR CodigoError = -1
)

type Comando interface {
	Ejecutar(consola Consola, opciones ...string) (res any, cod CodigoError, err error)

	Ayuda(con Consola, args ...string)
	TextoAyuda() string

	DescifrarOpciones(opciones Opciones) (Parametros, Opciones, Argumentos)

	AsignarPadre(Comando)
	EsOculto() bool

	DevolverNombre() string
	DevolverAliases() []string
}

type Config struct {
	EsOculto bool
}
type comando struct {
	Nombre      string
	Aliases     []string
	Uso         string
	Descripcion string
	Opciones    []string
	Oculto      bool

	accion   Accion
	comandos []Comando
	padre    Comando
}

func (c comando) TextoAyuda() string {
	nombre := c.Nombre + " (" + strings.Join(c.Aliases, ",") + ") "
	return nombre + cadena.TextoJustificado(c.Descripcion, 40, cadena.OpcionesFormato{Sangria: strings.Repeat(" ", 40-len(nombre)-2), Prefijo: strings.Repeat(" ", 40), Color: color.GrisFuente}) + "\n"
}

func (c comando) Ayuda(con Consola, args ...string) {
	con.ImprimirCadena(Cadena(cadena.Titulo(c.Nombre)))
	con.ImprimirCadena(Cadena(cadena.Subtitulo(c.Descripcion)))
	con.EscribirLinea(Cadena("Uso:"))
	con.EscribirLinea(Cadena("\t" + c.Uso))
	//con.EscribirLinea(Cadena("Ayuda").Negrita().Subrayada())
	con.EscribirLinea("Subcomandos:")

	for _, c := range c.comandos {
		if !c.EsOculto() {
			con.EscribirLinea(Cadena("  " + c.TextoAyuda()))
		}
	}
	if len(c.Opciones) > 0 {
		con.EscribirLinea(Cadena("Opciones Generales:"))
		for _, o := range c.Opciones {
			con.EscribirCadena(Cadena("\t" + o))
		}
		con.EscribirLinea("")
	}
	con.Imprimir()
}

func (c *comando) RegistrarComando(sub Comando) Comando {
	sub.AsignarPadre(c)
	c.comandos = append(c.comandos, sub)
	return c
}

func (c *comando) AsignarPadre(p Comando) {
	c.padre = p
}

func (c *comando) buscarSubComando(nombre string) (Comando, bool) {
	for _, c := range c.comandos {
		if c.DevolverNombre() == nombre || slices.Contains(c.DevolverAliases(), nombre) {
			return c, true
		}
	}
	return nil, false // [HACER] MEJORAR RETORNO...
}

func (c *comando) DescifrarOpciones(opciones []string) (Parametros, Opciones, Argumentos) {
	parametros := make(Parametros)
	banderas := make([]string, 0)
	argumentos := make([]any, 0)

	i := 0
	for i < len(opciones) {
		m := opciones[i]
		if strings.HasPrefix(m, "--") || strings.HasPrefix(m, "-") {
			if slices.Contains(c.Opciones, m) {
				banderas = append(banderas, utiles.Limpiar(m))
				i++
			} else {
				j := i + 1
				for j < len(opciones) && !(strings.HasPrefix(opciones[j], "--") || strings.HasPrefix(opciones[j], "-")) {
					j++
				}
				parametros[utiles.Limpiar(m)] = opciones[i+1 : j]
				i = j
			}
		} else {

			argumentos = append(argumentos, utiles.Limpiar(m))
			i++
		}
	}
	return parametros, banderas, argumentos
}

func (c *comando) Ejecutar(consola Consola, opciones ...string) (res any, cod CodigoError, err error) {

	if len(opciones) > 0 {
		sc, existe := c.buscarSubComando(opciones[0])
		if existe {
			return sc.Ejecutar(consola, opciones[1:]...)
		}
	}
	parametros, banderas, argumentos := c.DescifrarOpciones(opciones)
	if c.accion == nil {
		c.Ayuda(consola)
		return nil, EXITO, nil
	}
	return c.accion(consola, banderas, parametros, argumentos...)
}

func (c comando) EsOculto() bool {
	return c.Oculto
}

func (c comando) debeCerrar() bool {
	return false
}

func (c comando) DevolverNombre() string {
	return c.Nombre
}
func (c comando) DevolverAliases() []string {
	return c.Aliases
}

func NuevoComando(nombre string, uso string, aliases []string, descripcion string, accion Accion, opciones []string, config ...Config) *comando {

	cfg := Config{
		EsOculto: false,
	}
	if len(config) > 0 {
		cfg = config[0]
	}
	c := &comando{

		Nombre:      nombre,
		Uso:         uso,
		Aliases:     aliases,
		accion:      accion,
		Descripcion: descripcion,
		Opciones:    opciones,
		Oculto:      cfg.EsOculto,
	}

	c.RegistrarComando(
		&comando{
			Nombre:      "ayuda",
			Uso:         "ayuda",
			Aliases:     []string{"-a", "-h"},
			Descripcion: "Imprime la ayuda.",
			accion: Accion(
				func(con Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
					c.Ayuda(con, opciones...)
					return nil, EXITO, nil
				}),
			Opciones: []string{},
			Oculto:   false})

	return c
}

func AccionNula(f func()) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		f()
		return nil, EXITO, nil
	}
}

func AccionFalible(f func() error) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return nil, EXITO, f()
	}
}

func AccionEntrada(f func(a any)) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		if len(argumentos) < 1 {
			return nil, ERROR, errors.New("no se pudo ejecutar la función asociada a esta acción. La función requiere un argumento, y 0 fueron provistos. ")
		}
		f(argumentos[0])
		return nil, EXITO, nil
	}
}

func AccionSalida(f func() any) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return f(), EXITO, nil
	}
}

func AccionImprimible(f func(con Consola)) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		f(consola)
		return nil, EXITO, nil
	}
}
func AccionImprimibleFalible(f func(con Consola) error) Accion {
	return func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return nil, EXITO, f(consola)
	}
}
