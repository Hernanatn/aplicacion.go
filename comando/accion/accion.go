package accion

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/hernanatn/aplicacion.go/consola"
)

type Consola = consola.Consola
type Opciones = consola.Opciones
type Parametros = consola.Parametros
type Argumentos = []any

type CodigoError int

const (
	EXITO CodigoError = iota << 0
	ERROR CodigoError = -1
)

/*
	Se utiliza `unsafe.Pointer` para poder pasar funciones donde el resultado `res T`, sea de un tipo arbitrario pero concreto (en vez de `any`)

La función debe tener la siguiente firma:

	func (consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res T, cod CodigoError, err error)
*/
type PunteroAccion struct {
	P unsafe.Pointer
	T reflect.Type
}
type Accion[R any] func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res R, cod CodigoError, err error)

func (a Accion[R]) Puntero() PunteroAccion {
	return PunteroAccion{P: unsafe.Pointer(&a), T: a.R()}
}

func (a Accion[R]) R() reflect.Type {
	return reflect.TypeFor[Accion[R]]()
}

/*
	Crea un `PunteroAccion` a partir de una `func` compatible.
	Se utiliza `unsafe.Pointer` para poder pasar funciones donde el resultado `res T`, sea de un tipo arbitrario pero concreto (en vez de `any`)

La función debe tener la siguiente firma:

	func (consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res T, cod CodigoError, err error)
*/
func NuevoPunteroAccion(f any /* func (consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res T, cod CodigoError, err error) */) PunteroAccion {
	switch {
	case reflect.TypeOf(f).Kind() != reflect.Func:
		panic("f debe ser una función.")
	case reflect.TypeOf(f).NumOut() != 3:
		panic("f debe tener exactamente 3 parámetros de respuesta.")
	}
	return PunteroAccion{P: unsafe.Pointer(&f), T: reflect.TypeOf(f)}
}

/*
	Coerciona un `PunteroAccion` hacia una `Accion[r]` cuyo primer parámetro de retorno `res` es de tipo `T`.

La función subyacente debe tener la siguiente firma:

	func (consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res T, cod CodigoError, err error)

	y `R` debe ser *idéntico*, *convertible*, o *compatible* con el tipo verdaderamente retornado por la función subyacente.
*/
func InstanciarAccion[R any](p PunteroAccion) Accion[R] {

	O := reflect.TypeFor[Accion[R]]()
	switch {
	case p.T.Kind() != reflect.Func:
		panic("p debe apuntar a una función.")
	case p.T.NumOut() != 3:
		panic("p debe apuntar a una función que tiene exactamente 3 parámetros de respuesta.")
	case p.T.Out(0).Kind() == reflect.Struct && O.Out(0).Kind() == reflect.Struct:
		for i := 0; i < O.Out(0).NumField(); i++ {
			c1 := p.T.Out(0).Field(i)
			c2 := O.Out(0).Field(i)
			if !((c1.Type == c2.Type || c1.Type.ConvertibleTo(c2.Type) || c1.Type.AssignableTo(c2.Type)) && c1.Offset == c2.Offset) {
				panic("Tipos incompatibles")
			}
		}
		fallthrough
	case p.T.Out(0).ConvertibleTo(O.Out(0)),
		p.T.AssignableTo(O) || p.T.Out(0).AssignableTo(O.Out(0)),
		O.Size() <= p.T.Size() && O.Align() == p.T.Align():
		return *(*Accion[R])(p.P)
	default:
		panic("Tipos incompatibles")
	}
}

func AccionNula(f func()) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		f()
		return nil, EXITO, nil
	}
	return Accion[any](p)
}

func AccionFalible(f func() error) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return nil, EXITO, f()
	}
	return Accion[any](p)
}

func AccionEntrada(f func(a any)) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		if len(argumentos) < 1 {
			return nil, ERROR, errors.New("no se pudo ejecutar la función asociada a esta acción. La función requiere un argumento, y 0 fueron provistos. ")
		}
		f(argumentos[0])
		return nil, EXITO, nil
	}
	return Accion[any](p)
}

func AccionSalida(f func() any) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return f(), EXITO, nil
	}
	return Accion[any](p)
}

func AccionImprimible(f func(con Consola)) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		f(consola)
		return nil, EXITO, nil
	}
	return Accion[any](p)
}
func AccionImprimibleFalible(f func(con Consola) error) Accion[any] {
	p := func(consola Consola, opciones Opciones, parametros Parametros, argumentos ...any) (res any, cod CodigoError, err error) {
		return nil, EXITO, f(consola)
	}
	return Accion[any](p)
}
