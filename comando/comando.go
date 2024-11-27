package comando

import (
	"slices"
	"strings"

	"github.com/hernanatn/aplicacion.go/comando/accion"
	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
	"github.com/hernanatn/aplicacion.go/utiles"
)

type Consola = accion.Consola
type Opciones = accion.Opciones
type Parametros = accion.Parametros
type Argumentos = accion.Argumentos

/*
PunteroAccion

	Se utiliza `unsafe.Pointer` para poder pasar funciones donde el resultado `res T`, sea de un tipo arbitrario pero concreto (en vez de `any`)

La función debe tener la siguiente firma:

	func (consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res T, cod CodigoError, err error)
*/
type PunteroAccion = accion.PunteroAccion
type Accion[T any] accion.Accion[T]

type Cadena = consola.Cadena

type CodigoError = accion.CodigoError

const (
	EXITO = accion.EXITO
	ERROR = accion.ERROR
)

type Comando interface {
	Ayuda(con Consola, args ...string)
	TextoAyuda() string
	BuscarComando(nombre string) (Comando, bool)
	DescifrarOpciones(opciones Opciones) (Parametros, Opciones, Argumentos)

	AsignarPadre(Comando)
	EsOculto() bool

	DevolverNombre() string
	DevolverAliases() []string
	Accion() PunteroAccion
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

	accion   PunteroAccion
	comandos []Comando
	padre    Comando
}

func (c comando) TextoAyuda() string {
	nombre := c.Nombre + " (" + strings.Join(c.Aliases, ",") + ") "
	return nombre + cadena.TextoJustificado(c.Descripcion, 40, cadena.OpcionesFormato{Sangria: strings.Repeat(" ", 40-len(nombre)-2), Prefijo: strings.Repeat(" ", 40), Color: color.GrisFuente}) + "\n"
}

func (c comando) Ayuda(con Consola, args ...string) {
	con.EscribirLinea(Cadena("Ayuda").Negrita().Subrayada())
	con.EscribirLinea("Comandos:")

	for _, c := range c.comandos {
		if !c.EsOculto() {
			con.EscribirLinea(Cadena("  " + c.TextoAyuda()))
		}
	}
	if len(c.Opciones) > 0 {
		con.EscribirLinea("Opciones Generales:")
		for _, o := range c.Opciones {
			con.EscribirLinea(Cadena("  " + o))
		}
	}

	con.Imprimir()
}

func (c *comando) RegistrarComando(sub Comando) Comando {
	sub.AsignarPadre(c)
	c.comandos = append(c.comandos, sub)
	return c
}

func (c *comando) Accion() PunteroAccion {
	return c.accion
}

func (c *comando) AsignarPadre(p Comando) {
	c.padre = p
}

func (c *comando) BuscarComando(nombre string) (Comando, bool) {
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

	for i, m := range opciones {
		switch {
		case strings.Contains(m, "--"), strings.Contains(m, "-"):
			switch {
			case slices.Contains(c.Opciones, m):
				banderas = append(opciones, utiles.Limpiar(m))
			default:
				var j int
				for k, p := range opciones[i+1:] {
					if strings.Contains(p, "--") || strings.Contains(p, "-") {
						j = k
						parametros[m] = opciones[i+1 : j+i+1]
						break
					}
				}
			}
		default:
			argumentos = append(argumentos, utiles.Limpiar(m))
		}

	}
	return parametros, banderas, argumentos
}

func Ejecutar[T any](c Comando, consola Consola, opciones ...string) (res T, cod CodigoError, err error) {

	if len(opciones) > 0 {
		sc, existe := c.BuscarComando(opciones[0])
		if existe {
			return Ejecutar[T](sc, consola, opciones[1:]...)
		}
	}
	parametros, banderas, argumentos := c.DescifrarOpciones(opciones)
	if c.Accion().P == nil {
		c.Ayuda(consola)
		return *new(T), EXITO, nil
	}
	return accion.InstanciarAccion[T](c.Accion())(consola, banderas, parametros, argumentos...)
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

func NuevoComando(nombre string, uso string, aliases []string, descripcion string, accion PunteroAccion, opciones []string, config ...Config) Comando {

	cfg := Config{
		EsOculto: false,
	}
	if len(config) > 0 {
		cfg = config[0]
	}
	return &comando{

		Nombre:      nombre,
		Uso:         uso,
		Aliases:     aliases,
		accion:      accion,
		Descripcion: descripcion,
		Opciones:    opciones,
		Oculto:      cfg.EsOculto,
	}
}
