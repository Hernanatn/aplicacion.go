package comando

import (
	"slices"
	"strings"

	"aplicacion.go/consola"
	"aplicacion.go/consola/cadena"
	"aplicacion.go/consola/color"
	"aplicacion.go/utiles"
)

type Consola = consola.Consola
type Cadena = consola.Cadena

type Parametros = consola.Parametros

type AccionComando interface {
	Ejecutar(consola Consola, parametros Parametros, opciones ...string) (res any, cod CodigoError, err error)
}

type Accion func(consola Consola, parametros Parametros, opciones ...string) (res any, cod CodigoError, err error)

func (f Accion) Ejecutar(consola Consola, parametros Parametros, opciones ...string) (res any, cod CodigoError, err error) {
	return f(consola, parametros, opciones...)
}

type CodigoError int

const (
	EXITO CodigoError = 0
	ERROR CodigoError = -1
)

type Comando interface {
	Ejecutar(consola Consola, opciones ...string) (res any, cod CodigoError, err error)

	Ayuda(con Consola, args ...string)
	TextoAyuda() string

	DescifrarOpciones(opciones []string) (Parametros, []string)

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
	Accion      AccionComando
	Opciones    []string
	Oculto      bool

	comandos []Comando
	padre    Comando
}

type TComando struct {
	*comando
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

func (c *comando) DescifrarOpciones(opciones []string) (Parametros, []string) {
	parametros := make(Parametros)
	banderas := make([]string, 0)

	for i, m := range opciones {
		switch {
		case strings.Contains(m, "--"), strings.Contains(m, "-"):
			switch {
			case slices.Contains(c.Opciones, m):
				banderas = append(opciones, utiles.Limpiar(m))
			default:
				var j int
				for k, p := range opciones[i+1:] {
					switch {
					case strings.Contains(p, "--"), strings.Contains(p, "-"):
						j = k
						break
					}
					parametros[m] = opciones[i+1 : j+i+1]
				}
			}
		}

	}
	return parametros, banderas
}

func (c *comando) Ejecutar(consola Consola, opciones ...string) (res any, cod CodigoError, err error) {

	if len(opciones) > 0 {
		sc, existe := c.buscarSubComando(opciones[0])
		if existe {
			return sc.Ejecutar(consola, opciones[1:]...)
		}
	}
	parametros, banderas := c.DescifrarOpciones(opciones)
	if c.Accion == nil {
		c.Ayuda(consola)
		return nil, EXITO, nil
	}
	return c.Accion.Ejecutar(consola, parametros, banderas...)
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

func NuevoComando(nombre string, uso string, aliases []string, descripcion string, accion AccionComando, opciones []string, config ...Config) Comando {

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
		Accion:      accion,
		Descripcion: descripcion,
		Opciones:    opciones,
		Oculto:      cfg.EsOculto,
	}
}
