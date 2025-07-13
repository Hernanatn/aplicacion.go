package aplicacion

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/hernanatn/aplicacion.go/comando"
	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
	"github.com/hernanatn/aplicacion.go/menu"
	"github.com/hernanatn/aplicacion.go/menu/multimenu"
	"github.com/hernanatn/aplicacion.go/utiles"
)

type Cadena = comando.Cadena
type Consola = comando.Consola
type CodigoError = comando.CodigoError
type Opciones = consola.Opciones
type Parametros = consola.Parametros
type Argumentos = []any
type Accion = comando.Accion
type Comando = comando.Comando

type Menu = menu.Menu
type OpcionMenu = menu.Opcion

const (
	EXITO = comando.EXITO
	ERROR = comando.ERROR
)

var (
	NuevaConsola   = consola.NuevaConsola
	NuevoComando   = comando.NuevoComando
	NuevoMenu      = menu.NuevoMenu
	NuevoMultiMenu = multimenu.NuevoMultiMenu
)

type Aplicacion interface {
	Consola
	Comando

	Correr(args ...string) (r any, err error)

	Inicializar(...string) error
	Limpiar(...string) error
	Finalizar(...string) error

	Leer(Cadena) (Cadena, error)

	RegistrarInicio(f FUN) Aplicacion
	RegistrarLimpieza(f FUN) Aplicacion
	RegistrarFinal(f FUN) Aplicacion
	RegistrarComando(Comando) Aplicacion

	DebeCerrar() bool
}

type aplicacion struct {
	Nombre      string
	Uso         string
	Descripcion string
	accion      comando.Accion
	Opciones    []string

	consola    Consola
	comandos   []Comando
	debeCerrar bool
	ini        FUN
	lim        FUN
	fin        FUN
}

type FUN func(c Aplicacion, args ...string) error

func (a *aplicacion) Inicializar(args ...string) error {
	return a.ini(a, args...)
}
func (a *aplicacion) Limpiar(args ...string) error {
	return a.lim(a, args...)
}
func (a *aplicacion) Finalizar(args ...string) error {
	return a.fin(a, args...)
}

func (a aplicacion) TextoAyuda() string {
	return a.Nombre + cadena.TextoJustificado(a.Descripcion, 40, cadena.OpcionesFormato{Sangria: strings.Repeat(" ", 20-len(a.Nombre)-2), Prefijo: strings.Repeat(" ", 20), Color: color.GrisFuente}) + "\n"
}

func (a *aplicacion) Ayuda(_ Consola, args ...string) {
	if len(args) > 0 {
		c, existe := a.buscarComando(args[0])
		if existe {
			c.Ayuda(a, args[:1]...)
		}
	}
	a.ImprimirCadena(Cadena(cadena.Titulo(a.Nombre)))
	a.ImprimirCadena(Cadena(cadena.Subtitulo(a.Descripcion)))
	a.consola.EscribirLinea(Cadena("Uso:"))
	a.consola.EscribirLinea(Cadena("\t" + a.Uso))
	//a.consola.EscribirLinea(Cadena("Ayuda").Negrita().Subrayada())
	a.consola.EscribirLinea(Cadena("Comandos:"))
	for _, c := range a.comandos {
		if !c.EsOculto() && !slices.Contains(args, "-v") {
			a.consola.EscribirCadena(Cadena("  " + c.TextoAyuda()))
		}
	}
	if len(a.Opciones) > 0 {
		a.consola.EscribirLinea(Cadena("Opciones Generales:"))
		for _, o := range a.Opciones {
			a.consola.EscribirCadena(Cadena("\t" + o))
		}
		a.consola.EscribirLinea("")
	}
	a.consola.Imprimir()
}

func (a aplicacion) Consola() Consola {
	return a
}

func (a aplicacion) DevolverNombre() string {
	return a.Nombre
}
func (a aplicacion) DevolverAliases() []string {
	return []string{a.Nombre}
}

func (a *aplicacion) RegistrarComando(sub Comando) Aplicacion {
	sub.AsignarPadre(a)
	a.comandos = append(a.comandos, sub)
	return a
}

func (a aplicacion) buscarComando(nombre string) (Comando, bool) {
	for _, a := range a.comandos {
		if a.DevolverNombre() == nombre || slices.Contains(a.DevolverAliases(), nombre) {
			return a, true
		}
	}
	return nil, false // [HACER] MEJORAR RETORNO...
}

func (a *aplicacion) AsignarPadre(Comando) {}
func (a aplicacion) DescifrarOpciones(opciones []string) (Parametros, Opciones, Argumentos) {
	parametros := make(comando.Parametros)
	banderas := make([]string, 0)
	argumentos := make([]any, 0)

	for i, m := range opciones {
		switch {
		case strings.Contains(m, "--"), strings.Contains(m, "-"):
			switch {
			case slices.Contains(a.Opciones, m):
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

func (a *aplicacion) Ejecutar(_ Consola, opciones ...string) (res any, cod comando.CodigoError, err error) {

	if len(opciones) > 0 {
		sc, existe := a.buscarComando(opciones[0])
		if existe {
			return sc.Ejecutar(a, opciones[1:]...)
		}
	}
	parametros, banderas, argumentos := a.DescifrarOpciones(opciones)
	if a.accion == nil {
		a.Ayuda(a)
		return nil, comando.EXITO, nil
	}
	return a.accion(a, banderas, parametros, argumentos...)
}

func (a *aplicacion) RegistrarInicio(f FUN) Aplicacion {
	a.ini = FUN(f)
	return a
}
func (a *aplicacion) RegistrarLimpieza(f FUN) Aplicacion {
	a.lim = FUN(f)
	return a
}
func (a *aplicacion) RegistrarFinal(f FUN) Aplicacion {
	a.fin = FUN(f)
	return a
}

func (a aplicacion) Leer(c Cadena) (Cadena, error) {
	return a.consola.Leer(c)
}

func (e aplicacion) Read(p []byte) (n int, err error) {
	return e.consola.Read(p)
}
func (s aplicacion) Write(p []byte) (n int, err error) {
	return s.consola.Write(p)
}

/*
	func (a aplicacion) Escribir(c Cadena) error {
		_, err := a.salida.WriteString(c.S())
		if err != nil {
			return err
		}
		err = a.salida.Flush()
		if err != nil {
			return err
		}
		return nil
	}

	func (a aplicacion) Escribirf(f string, v ...any) error {
		_, err := a.salida.WriteString(fmt.Sprintf(f, v))
		if err != nil {
			return err
		}
		err = a.salida.Flush()
		if err != nil {
			return err
		}
		return nil
	}
*/

func (a aplicacion) BorrarLinea() error {
	return a.consola.BorrarLinea()
}
func (a aplicacion) Imprimir() error {
	return a.consola.Imprimir()
}
func (a aplicacion) ImprimirLinea(c Cadena) error {
	return a.consola.ImprimirLinea(c)
}

func (a aplicacion) ImprimirSeparador() {
	a.consola.ImprimirSeparador()

}
func (a aplicacion) EsOculto() bool {
	return false
}
func (a aplicacion) DebeCerrar() bool {
	return a.debeCerrar
}

func (a *aplicacion) Correr(args ...string) (r any, err error) {
	var res any
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlC
		a.Limpiar(args...)
		a.ImprimirError("Programa terminado por el usuario: [CTRL+C]", nil)
		os.Exit(1)
	}()
	err = a.Inicializar(args...)
	if err != nil {
		a.Limpiar(args...)
		a.ImprimirFatal("No se pudo inicializar la aplicacion", err)
		return nil, *new(error)
	}

	a.Ejecutar(a.consola, args...)
	for !a.DebeCerrar() {
		entrada, err := a.Leer("")
		if err != nil {
			if err == io.EOF {
				a.debeCerrar = true
				a.Limpiar(args...)
				a.ImprimirError("Programa terminado por el usuario: [CTRL+C]", nil)
				os.Exit(1)
			}
			a.Limpiar(args...)
			a.ImprimirFatal("No se pudo leer desde la entrada de la Aplicación", err)
			return nil, err
		}

		argumentos := strings.Split(entrada.Limpiar().S(), " ")
		var com Comando
		nombreComando := argumentos[0]
		com, existe := a.buscarComando(nombreComando)
		opciones := argumentos[1:]
		switch {
		case len(argumentos) < 1 || len(argumentos) == 1 && nombreComando == "":
			continue
		case !existe:
			a.ImprimirError(cadena.Cadena(fmt.Sprintf("Se intento ejecutar el comando: %s. Pero el comando no existe", nombreComando)), nil)
			a.Ayuda(a)
			return nil, nil
		}
		var cod comando.CodigoError
		res, cod, err = com.Ejecutar(a, opciones...)

		if err != nil {
			a.ImprimirFatal(cadena.Cadena(strconv.Itoa(int(cod))), err)
			a.ImprimirFatal("No se pudo ejecutar correctamente la aplicacion", err)
			return nil, err
		}
	}

	err = a.Finalizar(args...)
	if err != nil {
		a.Limpiar(args...)
		a.ImprimirFatal("No se pudo finalizar correctamente la aplicacion", err)
		return nil, *new(error)
	}

	return res, nil
}

func NuevaAplicacion(nombre string, uso string, descripcion string, opciones []string, consola Consola) Aplicacion {

	a := &aplicacion{
		Nombre:      nombre,
		Uso:         uso,
		Descripcion: descripcion,
		Opciones:    opciones,
		consola:     consola,
	}

	a.RegistrarComando(
		comando.NuevoComando(
			"ayuda",
			"ayuda",
			[]string{"-a", "-h"},
			"Imprime la ayuda.",
			comando.Accion(
				func(con Consola, opciones comando.Opciones, parametros comando.Parametros, argumentos ...any) (res any, cod comando.CodigoError, err error) {
					a.Ayuda(con, opciones...)
					return nil, comando.EXITO, nil
				}),
			[]string{}))
	a.RegistrarComando(
		comando.NuevoComando(
			"chau",
			"chau",
			[]string{},
			"Cierra el programa.",
			comando.Accion(
				func(con Consola, opciones comando.Opciones, parametros comando.Parametros, argumentos ...any) (res any, cod comando.CodigoError, err error) {
					a.debeCerrar = true
					return nil, comando.EXITO, nil
				}),
			[]string{},
			comando.Config{
				EsOculto: true,
			}))
	return a
}

func (a aplicacion) LeerContraseña(mensaje Cadena) (Cadena, error) {
	return a.consola.LeerContraseña(mensaje)
}
func (a aplicacion) LeerTecla(b *[]byte) (int, error) {
	return a.consola.LeerTecla(b)
}

// Escribe la Cadena al buffer y llama Imprimir()
func (a aplicacion) ImprimirCadena(cadena Cadena) error {
	return a.consola.ImprimirCadena(cadena)
}

// Escribe los bytes al buffer y llama Imprimir()
func (a aplicacion) ImprimirBytes(b []byte) error {
	return a.consola.ImprimirBytes(b)
}

// Escribe la Cadena al buffer
func (a aplicacion) EscribirCadena(cadena Cadena) error {
	return a.consola.EscribirCadena(cadena)
}

// Escribe la Cadena +\r\n al buffer
func (a aplicacion) EscribirLinea(cadena Cadena) error {
	return a.consola.EscribirLinea(cadena)
}

// Escribe los bytes al buffers
func (a aplicacion) EscribirBytes(b []byte) error {
	return a.consola.EscribirBytes(b)
}

// Escribe la Cadena al buffer, la formatea como Advertencia y llama Imprimir()
func (a aplicacion) ImprimirAdvertencia(cadena Cadena, e error) error {
	return a.consola.ImprimirAdvertencia(cadena, e)
}

// Escribe la Cadena al buffer, la formatea como Error y llama Imprimir()
func (a aplicacion) ImprimirError(cadena Cadena, e error) error {
	return a.consola.ImprimirError(cadena, e)
}

// Escribe la Cadena al buffer, la formatea como Fatal y llama Imprimir()
func (a aplicacion) ImprimirFatal(cadena Cadena, e error) error {
	return a.consola.ImprimirFatal(cadena, e)
}

func (a aplicacion) EsTerminal() bool {
	return a.consola.EsTerminal()
}

func (a aplicacion) FEntrada() *os.File {
	return a.consola.FEntrada()
}
func (a aplicacion) FSalida() *os.File {
	return a.consola.FSalida()
}
