# aplicacion.go

Librería en GO para crear aplicaciones con interfaz basada en la Terminal o Línea de Comandos.

La librería se presenta aún en versión beta - distribución `v 0.2.4-beta`, no está particularmente optimizada y puede presentar errores.

[![Hecho por Chaska](https://img.shields.io/badge/hecho_por-Ch'aska-303030.svg)](https://cajadeideas.ar)
[![Versión: Beta v0.1](https://img.shields.io/badge/version-Beta_v0.2.4-orange.svg)](https://github.com/hernanatn/github.com/hernanatn/aplicacion.go/releases/latest)
[![Verisón de Go: 1.22.5](https://img.shields.io/badge/Go-1.22.5-blue?logo=go)](https://go.dev/doc/go1.22)
[![Licencia: CC BY-NC 4.0](https://img.shields.io/badge/Licencia-CC_BY--SA_4.0-lightgrey.svg)](LICENSE)
## Documentación
De momento se ofrece como documentación únicamente la generada automáticamente por [`gomarkdoc`](<https://github.com/princjef/gomarkdoc>), a la cual se puede acceder desde [/documentación]().

Se ofrece, además, un [ejemplo de úso básico](#uso-básico) y una pequeña descripción de las [Interfaces Públicas](#interfaz-pública-simple) provistas.

Si se desea conocer más acerca de la implementación y funcionamiento de la librería, se recomienda revisar el código fuente, el cual está escrito en español y con nombres *relativamente* descriptivos.

Los comentarios son inconsistentes y no deben ser toados como documentación.


## Uso Básico:
1. Creamos una nueva `Aplicacion` a partir de una Consola instanciada con `os.Stdin` y `os.Stdout`
```go
	var app aplicacion.Aplicacion = aplicacion.NuevaAplicacion(
		"App", // nombre
		"app", // uso
		"app / v 0.1", // descripción
		make([]string, 0), // opciones
		aplicacion.NuevaConsola(os.Stdin, os.Stdout), // consola
	)
```
2. Creamos y registramos los comandos relevantes 
```go
	var prueba aplicacion.Comando = aplicacion.NuevoComando(
		"prueba",
		"prueba -p <PROYOECTO> -r <RAIZ PROYECTO> [OPCIONES] --> [OPCIONES flask]",
		[]string{"run"},
		"Corre el servidor con flask.",
		func(con aplicacion.Consola, parametros aplicacion.Parametros, opciones ...string) (res any, cod aplicacion.CodigoError, err error) {
			con.ImprimirLinea(aplicacion.Cadena("Comando de prueba"))
			return nil, aplicacion.EXITO, nil
		},
		make([]string, 0),
	)

	app.RegistrarComando(prueba)

	// ... otros comandos
```
3. (Opcional) registramos las acciones de inicialización, limpieza, y finalización del programa
```go

	func ini(a aplicacion.Aplicacion, args ...string) error {
		a.ImprimirLinea(aplicacion.Cadena("¡Hola!"))
		return nil
	}
	func lim(a aplicacion.Aplicacion, args ...string) error {
		a.ImprimirLinea(aplicacion.Cadena("¡Adiós!"))
		return nil
	}
	func fin(a aplicacion.Aplicacion, args ...string) error {
		a.ImprimirLinea(aplicacion.Cadena("Limpiando..."))
		return nil
	}

	app.
		RegistrarInicio(ini).
		RegistrarLimpieza(lim).
		RegistrarFinal(fin)
```

4. Corremos la aplicación:
```go
	app.Correr(os.Args[1:]...)
```

## Interfaz Pública Simple:
Se ofrece la interfaz `Aplicacion`, y sus interfaces asociadas: `Consola`, `Comando`. 
Una `Aplicacion` consiste de un `Comando` y una `Consola`.
Un `Comando` es funcionalidad (ejecutable mediante el método `Ejecutar` conforme la interfaz definida) asociada a un nombre y aliases, y que puede contener un arreglo de subcomandos.

```go

func NuevaConsola(fe *os.File, fs *os.File) *Consola

func NuevoComando(nombre string, uso string, aliases []string, descripcion string, accion AccionComando, opciones []string, config ...Config) Comando

func NuevaAplicacion(nombre string, uso string, descripcion string, opciones []string, consola *Consola) Aplicacion


```
### Las tres interfaces principales son `Aplicacion`, `Comando` y `Consola`
```go
type Aplicacion interface {
	Consola
	Comando

	Correr(args ...string) (r any, err error)

	Inicializar(...string) error
	Limpiar(...string) error
	Finalizar(...string) error

	Leer(Cadena) (Cadena, error)

	Consola() Consola

	RegistrarInicio(f FUN) Aplicacion
	RegistrarLimpieza(f FUN) Aplicacion
	RegistrarFinal(f FUN) Aplicacion
	RegistrarComando(Comando) Aplicacion

	DebeCerrar() bool
}
```
```go
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
```
```go
type Consola interface {
	Leer(Cadena) (Cadena, error)
	LeerContraseña(Cadena) (Cadena, error)
	LeerTecla(*[]byte) (int, error)

	Imprimir() error
	ImprimirLinea(Cadena) error
	ImprimirCadena(Cadena) error
	BorrarLinea() error
	ImprimirBytes([]byte) error
	EscribirLinea(Cadena) error
	EscribirCadena(Cadena) error
	ImprimirSeparador()
	EscribirBytes([]byte) error
	EsTerminal() bool
}
```

## Interfaz Pública Avanzada
La librería incluye módulos segregados y con interfaz pública propia:
- `consola` operaciones de lectura y escritura hacia y desde un búfer compatible con una terminal;
- `menu` funcionalidad de menu de opciones renderizado en línea de comandos;
- `cadena` alias de `string` con métodos personalizados para formato compatible con ANSI

Además de `Aplicacion`, `Consola` y `Comando`, se ofrecen los tipos `Menu`, `MultiMenu` y `Cadena`:
```go

type Opcion struct {
	Nombre string
	Accion Accion //comando.Accion
}

type Menu struct {
	Opciones []*Opcion
	Consola  consola.Consola

	cursor       rune
	seleccionada int
	largo        int
	debeCerrar   bool
}

func NuevoMenu(con consola.Consola, cur rune) *Menu

func (m *Menu) RegistrarOpcion(o *Opcion) *Menu
func (m *Menu) Correr() (*Opcion, error)
func (m Menu) DebeCerrar() bool 
```
```go
type MultiMenu struct {
	Opciones []*Opcion
	Consola  consola.Consola

	cursor        rune
	enfocada      int
	seleccionadas []int
	largo         int
	debeCerrar    bool
}

func NuevoMultiMenu(con consola.Consola, cur rune) *Menu

func (m *Menu) RegistrarOpcion(o *Opcion) *Menu
func (m *Menu) Correr() ([]*Opcion, error)
func (m Menu) DebeCerrar() bool 
```
```go
type Cadena string
func (c Cadena) Colorear(col color.Color) Cadena
func (c Cadena) Imprimir(f *bufio.Writer)
func (c Cadena) Invertida() Cadena
func (c Cadena) Italica() Cadena
func (c Cadena) Limpiar() Cadena
func (c Cadena) Negrita() Cadena
func (c Cadena) Subrayada() Cadena
func (c Cadena) Formatear(formatos ...Formateador) Cadena
func (c Cadena) S() string 
func (c Cadena) String() string // Alias de S()

func DesdeArchivo(nombre string) (Cadena, error)
```

Por simplicidad se optó por importar al espacio de nombres `aplicacion` todos los tipos, interfaces y funciones públicas releventes. 

Los nombres expuestos por el módulo `aplicacion` son:
```go
type Cadena // cadena.Cadena
type Consola interface // consola.Consola
type CodigoError // comando.CodigoError
type Parametros // consola.Parametros
type Opciones // consola.Opciones
type Argumentos // comando.Argumentos
type Accion // comando.Accion
type Comando interface // comando.Comando
type Menu interface // menu.Menu
type OpcionMenu // menu.Opcion
type FUN
type Aplicacion interface

const EXITO // 0
const ERROR // 1

// consola.NuevaConsola
func NuevaConsola(fe *os.File, fs *os.File) *Consola 

// comando.NuevoComando
func NuevoComando(nombre string, uso string, aliases []string, descripcion string, accion AccionComando, opciones []string, config ...Config) Comando 

// menu.NuevoMenu
func NuevoMenu(con consola.Consola, cur rune) *Menu 

// multimenu.NuevoMultimenMenu
func NuevoMultiMenu(con consola.Consola, cur rune) *Menu 

func NuevaAplicacion(nombre string, uso string, descripcion string, opciones []string, consola *Consola) Aplicacion
 
```

Librería Desarrollada por Hernán A.T.N. para Ch'aska S.R.L. y distribuída bajo [Licencia CC BY-SA 4.0][cc-by-sa].  Derechos de autor (c) 2023 Ch'aska S.R.L. 

---

[![CC BY-SA 4.0][cc-by-sa-image]][cc-by-sa] [CH'ASKA](https://cajadeideas.ar)

[cc-by-sa]: http://creativecommons.org/licenses/by-sa/4.0/
[cc-by-sa-image]: https://licensebuttons.net/l/by-sa/4.0/88x31.png
[cc-by-sa-shield]: https://img.shields.io/badge/License-CC%20BY--SA%204.0-lightgrey.svg