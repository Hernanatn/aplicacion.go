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
	"github.com/hernanatn/aplicacion.go/utiles"
)

type Cadena = comando.Cadena
type Consola = comando.Consola
type CodigoError = comando.CodigoError
type Opciones = comando.Opciones
type Parametros = comando.Parametros

type Menu = menu.Menu

const (
	EXITO = comando.EXITO
	ERROR = comando.ERROR
)

var (
	NuevaConsola = consola.NuevaConsola
	NuevoMenu    = menu.NuevoMenu
)

type Aplicacion[T any] interface {
	Consola
	comando.Comando[T]

	Correr(args ...string) (r any, err error)

	Inicializar(...string) error
	Limpiar(...string) error
	Finalizar(...string) error

	Leer(Cadena) (Cadena, error)

	RegistrarInicio(f FUN[T]) Aplicacion[T]
	RegistrarLimpieza(f FUN[T]) Aplicacion[T]
	RegistrarFinal(f FUN[T]) Aplicacion[T]
	RegistrarComando(comando.Comando[T]) Aplicacion[T]

	DebeCerrar() bool
}

type aplicacion[T any] struct {
	Nombre      string
	Uso         string
	Descripcion string
	accion      comando.Accion[T]
	Opciones    []string

	consola    Consola
	comandos   []comando.Comando[T]
	debeCerrar bool
	ini        FUN[T]
	lim        FUN[T]
	fin        FUN[T]
}

type FUN[T any] func(c Aplicacion[T], args ...string) error

func (a *aplicacion[T]) Inicializar(args ...string) error {
	return a.ini(a, args...)
}
func (a *aplicacion[T]) Limpiar(args ...string) error {
	return a.lim(a, args...)
}
func (a *aplicacion[T]) Finalizar(args ...string) error {
	return a.fin(a, args...)
}

func (a aplicacion[T]) TextoAyuda() string {
	return a.Nombre + cadena.TextoJustificado(a.Descripcion, 40, cadena.OpcionesFormato{Sangria: strings.Repeat(" ", 20-len(a.Nombre)-2), Prefijo: strings.Repeat(" ", 20), Color: color.GrisFuente}) + "\n"
}

func (a *aplicacion[T]) Ayuda(_ Consola, args ...string) {
	a.ImprimirCadena(Cadena(cadena.Titulo(a.Nombre)))
	a.ImprimirCadena(Cadena(cadena.Subtitulo(a.Descripcion)))
	a.consola.EscribirLinea(Cadena("Ayuda").Negrita().Subrayada())
	a.consola.EscribirLinea(Cadena("Comandos:"))
	for _, c := range a.comandos {
		if !c.EsOculto() && !slices.Contains(args, "-v") {
			a.consola.EscribirCadena(Cadena("  " + c.TextoAyuda()))
		}
	}
	if len(a.Opciones) > 0 {
		a.consola.EscribirLinea(Cadena("Opciones Generales:"))
		for _, o := range a.Opciones {
			a.consola.EscribirCadena(Cadena("  " + o))
		}
	}
	a.consola.Imprimir()
}

func (a aplicacion[T]) Consola() Consola {
	return a
}

func (a aplicacion[T]) DevolverNombre() string {
	return a.Nombre
}
func (a aplicacion[T]) DevolverAliases() []string {
	return []string{a.Nombre}
}

func (a *aplicacion[T]) RegistrarComando(sub comando.Comando[T]) Aplicacion[T] {
	sub.AsignarPadre(a)
	a.comandos = append(a.comandos, sub)
	return a
}

func (a aplicacion[T]) buscarComando(nombre string) (comando.Comando[T], bool) {
	for _, a := range a.comandos {
		if a.DevolverNombre() == nombre || slices.Contains(a.DevolverAliases(), nombre) {
			return a, true
		}
	}
	return nil, false // [HACER] MEJORAR RETORNO...
}

func (a *aplicacion[T]) AsignarPadre(comando.Comando[T]) {}
func (a aplicacion[T]) DescifrarOpciones(opciones []string) (comando.Parametros, []string) {
	parametros := make(comando.Parametros)
	banderas := make([]string, 0)

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
		}

	}
	return parametros, banderas
}

func (a *aplicacion[T]) Ejecutar(_ Consola, opciones ...string) (res T, cod comando.CodigoError, err error) {

	if len(opciones) > 1 {
		sc, existe := a.buscarComando(opciones[1])
		if existe {
			return sc.Ejecutar(a, opciones[1:]...)
		}
	}
	parametros, banderas := a.DescifrarOpciones(opciones)
	if a.accion == nil {
		a.Ayuda(a)
		return *new(T), comando.EXITO, nil
	}
	return a.accion(a, banderas, parametros)
}

func (a *aplicacion[T]) RegistrarInicio(f FUN[T]) Aplicacion[T] {
	a.ini = FUN[T](f)
	return a
}
func (a *aplicacion[T]) RegistrarLimpieza(f FUN[T]) Aplicacion[T] {
	a.lim = FUN[T](f)
	return a
}
func (a *aplicacion[T]) RegistrarFinal(f FUN[T]) Aplicacion[T] {
	a.fin = FUN[T](f)
	return a
}

func (a aplicacion[T]) Leer(c Cadena) (Cadena, error) {
	return a.consola.Leer(c)
}

func (e aplicacion[T]) Read(p []byte) (n int, err error) {
	return e.consola.Read(p)
}
func (s aplicacion[T]) Write(p []byte) (n int, err error) {
	return s.consola.Write(p)
}

/*
	func (a aplicacion[T]) Escribir(c Cadena) error {
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

	func (a aplicacion[T]) Escribirf(f string, v ...any) error {
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

func (a aplicacion[T]) BorrarLinea() error {
	return a.consola.BorrarLinea()
}
func (a aplicacion[T]) Imprimir() error {
	return a.consola.Imprimir()
}
func (a aplicacion[T]) ImprimirLinea(c Cadena) error {
	return a.consola.ImprimirLinea(c)
}

func (a aplicacion[T]) ImprimirSeparador() {
	a.consola.ImprimirSeparador()

}
func (a aplicacion[T]) EsOculto() bool {
	return false
}
func (a aplicacion[T]) DebeCerrar() bool {
	return a.debeCerrar
}

func (a *aplicacion[T]) Correr(args ...string) (r any, err error) {
	var res any
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlC
		a.Limpiar(args...)
		a.ImprimirLinea(cadena.Cadena("\n" + cadena.Error("Programa terminado por el usuario: [CTRL+C]", nil)))
		os.Exit(1)
	}()
	err = a.Inicializar(args...)
	if err != nil {
		a.Limpiar(args...)
		a.ImprimirCadena(Cadena(cadena.Fatal("No se pudo inicializar la aplicacion[T]", err)))
		return nil, *new(error)
	}

	a.Ejecutar(a.consola, args...)
	for !a.DebeCerrar() {
		entrada, err := a.Leer("")
		if err != nil {
			if err == io.EOF {
				a.debeCerrar = true
				a.Limpiar(args...)
				a.ImprimirLinea(cadena.Cadena("\n" + cadena.Error("Programa terminado por el usuario: [CTRL+C]", nil)))
				os.Exit(1)
			}
			a.Limpiar(args...)
			a.ImprimirCadena(Cadena(cadena.Fatal("No se pudo leer desde la entrada de la Aplicación", err)))
			return nil, err
		}

		argumentos := strings.Split(entrada.Limpiar().S(), " ")
		var com comando.Comando[T]
		nombreComando := argumentos[0]
		com, existe := a.buscarComando(nombreComando)
		opciones := argumentos[1:]
		switch {
		case len(argumentos) < 1 || len(argumentos) == 1 && nombreComando == "":
			continue
		case !existe:
			a.ImprimirCadena(Cadena(cadena.Error(fmt.Sprintf("Se intento ejecutar el comando: %s. Pero el comando no existe", nombreComando), nil)))
			a.Ayuda(a)
			return nil, nil
		}
		var cod comando.CodigoError
		res, cod, err = com.Ejecutar(a, opciones...)

		if err != nil {
			a.ImprimirCadena(Cadena(cadena.Fatal(strconv.Itoa(int(cod)), err)))
			a.ImprimirCadena(Cadena(cadena.Fatal("No se pudo ejecutar correctamente la aplicacion[T]", err)))
			return nil, err
		}
	}

	err = a.Finalizar(args...)
	if err != nil {
		a.Limpiar(args...)
		a.ImprimirCadena(Cadena(cadena.Fatal("No se pudo finalizar correctamente la aplicacion[T]", err)))
		return nil, *new(error)
	}

	return res, nil
}

func NuevaAplicacion[T any](nombre string, uso string, descripcion string, opciones []string, consola Consola) Aplicacion[T] {

	a := &aplicacion[T]{
		Nombre:      nombre,
		Uso:         uso,
		Descripcion: descripcion,
		Opciones:    opciones,
		consola:     consola,
	}

	a.RegistrarComando(
		comando.NuevoComando[T](
			"ayuda",
			"ayuda",
			[]string{"-a", "-h"},
			"Imprime la ayuda.",
			comando.AccionImprimibleNula[T](
				func(con Consola) {
					a.Ayuda(con, opciones...)
				}),
			[]string{}))
	a.RegistrarComando(
		comando.NuevoComando[T](
			"chau",
			"chau",
			[]string{},
			"Cierra el programa.",
			comando.AccionNula[T](func() { a.debeCerrar = true }),
			[]string{},
			comando.Config{
				EsOculto: true,
			}))
	return a
}

func (a aplicacion[T]) LeerContraseña(mensaje Cadena) (Cadena, error) {
	return a.consola.LeerContraseña(mensaje)
}
func (a aplicacion[T]) LeerTecla(b *[]byte) (int, error) {
	return a.consola.LeerTecla(b)
}

// Escribe la Cadena al buffer y llama Imprimir()
func (a aplicacion[T]) ImprimirCadena(cadena Cadena) error {
	return a.consola.ImprimirCadena(cadena)
}

// Escribe los bytes al buffer y llama Imprimir()
func (a aplicacion[T]) ImprimirBytes(b []byte) error {
	return a.consola.ImprimirBytes(b)
}

// Escribe la Cadena al buffer
func (a aplicacion[T]) EscribirCadena(cadena Cadena) error {
	return a.consola.EscribirCadena(cadena)
}

// Escribe la Cadena +\r\n al buffer
func (a aplicacion[T]) EscribirLinea(cadena Cadena) error {
	return a.consola.EscribirLinea(cadena)
}

// Escribe los bytes al buffers
func (a aplicacion[T]) EscribirBytes(b []byte) error {
	return a.consola.EscribirBytes(b)
}

func (a aplicacion[T]) EsTerminal() bool {
	return a.consola.EsTerminal()
}

func (a aplicacion[T]) FEntrada() *os.File {
	return a.consola.FEntrada()
}
func (a aplicacion[T]) FSalida() *os.File {
	return a.consola.FSalida()
}
