/*
Utilidades de formateo de cadenas para escribir a terminales.
*/
package consola

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hernanatn/aplicacion.go/consola/cadena"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
)

type Cadena = cadena.Cadena
type Progreso = progressbar.ProgressBar

// Entrada es un *bufio.Reader que guarda una referencia al *os.File subyacente utilzado para constuir el Reader, junto con una bandera que determina si ese File es una terminal.
// Implementa una serie de métodos como comodidades para leer desde ese Reader:
type Entrada struct {
	*bufio.Reader
	f          *os.File
	esTerminal bool
}

type Salida struct {
	*bufio.Writer
	f          *os.File
	esTerminal bool
}

type EntradaSalida struct {
	Entrada
	Salida
	esTerminal bool
}

type Consola interface {
	io.ReadWriter
	Leer(Cadena) (Cadena, error)
	LeerContraseña(Cadena) (Cadena, error)
	LeerTecla(*[]byte) (int, error)

	Imprimir() error
	ImprimirLinea(Cadena) error
	ImprimirCadena(Cadena) error
	BorrarLinea() error
	//ImprimirString(string) error
	ImprimirBytes([]byte) error
	EscribirLinea(Cadena) error
	EscribirCadena(Cadena) error
	//EscribirString(string) error
	ImprimirSeparador()
	EscribirBytes([]byte) error
	EsTerminal() bool
}

type Parametros map[string]any

type consola struct {
	EntradaSalida
}

func (c consola) Leer(mensaje Cadena) (Cadena, error) {
	c.ImprimirCadena(cadena.Señalador(">") + mensaje + Cadena(": "))
	s, err := c.Entrada.ReadString('\n')
	if err != nil {
		return Cadena("\n"), err
	}

	return Cadena(s).Limpiar(), nil
}
func (c consola) LeerContraseña(mensaje Cadena) (Cadena, error) {
	c.ImprimirCadena(cadena.Señalador(">") + mensaje + Cadena(": "))
	viejo, _ := term.MakeRaw(int(c.EntradaSalida.Entrada.Fd()))
	t := term.NewTerminal(c.EntradaSalida, "")
	defer term.Restore(int(c.EntradaSalida.Entrada.Fd()), viejo)
	contraseña, err := t.ReadPassword("")
	return Cadena(contraseña), err
}

func (c consola) LeerTecla(b *[]byte) (int, error) {
	viejo, _ := term.MakeRaw(int(c.EntradaSalida.Entrada.Fd()))
	defer term.Restore(int(c.EntradaSalida.Entrada.Fd()), viejo)
	return c.Entrada.f.Read(*b)
}

// Llama Flush() en la Salida subyacente
func (c consola) Imprimir() error {
	return c.Salida.Flush()
}

// Escribe la Cadena al buffer y llama Imprimir()
func (c consola) ImprimirCadena(cadena Cadena) error {
	err1 := c.EscribirCadena(cadena)
	err2 := c.Imprimir()
	return errors.Join(err1, err2)
}

// Escribe la Cadena al buffer y llama Imprimir()
func (c consola) ImprimirLinea(cadena Cadena) error {
	err1 := c.EscribirCadena(cadena + Cadena("\r\n"))
	err2 := c.Imprimir()
	return errors.Join(err1, err2)
}

// Escribe los bytes al buffer y llama Imprimir()
func (c consola) ImprimirBytes(b []byte) error {
	err1 := c.EscribirBytes(b)
	err2 := c.Imprimir()
	return errors.Join(err1, err2)
}

// Escribe la Cadena al buffer
func (c consola) EscribirCadena(cadena Cadena) error {
	_, err := c.Writer.WriteString(cadena.S())
	return err
}

// Escribe la Cadena +\r\n al buffer
func (c consola) EscribirLinea(cadena Cadena) error {
	_, err := c.Writer.WriteString(cadena.S() + "\r\n")
	return err
}

// Escribe los bytes al buffer
func (c consola) EscribirBytes(b []byte) error {
	_, err := c.Writer.Write(b)
	return err
}

func (c consola) EsTerminal() bool {
	return true
}
func NuevaEntrada(f *os.File) *Entrada {
	return &Entrada{
		bufio.NewReader(f),
		f,
		term.IsTerminal(int(f.Fd())),
	}
}

// Crea una *Salida a partir de un *os.File
func NuevaSalida(f *os.File) *Salida {
	return &Salida{
		bufio.NewWriter(f),
		f,
		term.IsTerminal(int(f.Fd())),
	}
}

// Crea una *Salida a partir de un *bufio.Writer y un *os.File (pensado para utilizar con io.MultiWriter declarando el File autoritativo)
//
// # Ejemplo:
//
//	var buf bytes.Buffer
//	lector := bufio.NewWriter(io.MultiWriter(&buf, os.Stdout))
//	multiSalida := NuevaMultiSalida(lector, os.Stdout)
//
// Todas las operaciones de escritura se duplicaran entre buf y os.Stdout,
// pero se considerará os.Stdout como la salida subyacente autoritativa para todos los procedimientos que dependen de Salida.f
func NuevaMultiSalida(w *bufio.Writer, f *os.File) *Salida {
	return &Salida{
		w,
		f,
		term.IsTerminal(int(f.Fd())),
	}
}

func NuevaEntradaSalida(
	fe *os.File,
	fs *os.File,
) *EntradaSalida {
	return &EntradaSalida{

		*NuevaEntrada(fe),
		*NuevaSalida(fs),
		term.IsTerminal(int(fe.Fd())) || term.IsTerminal(int(fs.Fd())),
	}
}

func NuevaConsola(fe *os.File, fs *os.File) *consola {
	return &consola{
		EntradaSalida: *NuevaEntradaSalida(fe, fs),
	}
}

func NuevaEntradaMultiSalida(
	fe *os.File,
	w *bufio.Writer,
	fs *os.File,
) *EntradaSalida {
	return &EntradaSalida{

		*NuevaEntrada(fe),
		*NuevaMultiSalida(w, fs),
		term.IsTerminal(int(fe.Fd())) || term.IsTerminal(int(fs.Fd())),
	}
}

func (s Salida) EsTerminal() bool {
	return s.esTerminal
}

// Devuelve el tamaño de la terminal asociada a s.f
// Si s.f no es una terminal, devulve 0,0 para el tamaño, y err.
func (s Salida) DevolverTamaño() (int, int, error) {
	if !s.esTerminal {
		return 0, 0, *new(error)
	}
	ancho, alto, err := term.GetSize(int(s.f.Fd()))
	return ancho, alto, err
}

func (e Entrada) EsTerminal() bool {
	return e.esTerminal
}

// Devuelve el tamaño de la terminal asociada a e.f
// Si e.f no es una terminal, devulve 0,0 para el tamaño, y err.
func (e Entrada) DevolverTamaño() (int, int, error) {
	if !e.esTerminal {
		return 0, 0, *new(error)
	}
	ancho, alto, err := term.GetSize(int(e.f.Fd()))
	return ancho, alto, err
}

func Limpiar() {
	cmd := exec.Command("clear")
	_, _ = cmd.Output()
}

func (s Salida) BorrarLinea() error {
	ancho, _, err := s.DevolverTamaño()
	if err != nil {
		fmt.Printf("\r%s\r", strings.Repeat(" ", 56))
	}
	return s.EscribirCadena(Cadena(fmt.Sprintf("\r%s\r", strings.Repeat(" ", ancho))))

}

func (s Salida) ImprimirSeparador() {
	ancho, _, err := s.DevolverTamaño()
	if err != nil {
		fmt.Printf("\n%s\n", strings.Repeat("-", 56))
	}
	fmt.Printf("\n%s\n", strings.Repeat("-", ancho))
}

func (s *Salida) EscribirCadena(c Cadena) error {
	_, err := s.WriteString(c.S())
	if err != nil {
		return err
	}
	return nil
}
func (s *Salida) ImprimirCadena(c Cadena) error {
	err1 := s.EscribirCadena(c)
	err2 := s.Flush()

	err := errors.Join(err1, err2)
	return err
}

func Separador(f *os.File) string {
	err := *new(error)
	if !term.IsTerminal(int(f.Fd())) {
		return fmt.Sprintf("\n%s\n", strings.Repeat("-", 56))
	}
	ancho, _, err := term.GetSize(int(f.Fd()))
	if err != nil {
		return fmt.Sprintf("\n%s\n", strings.Repeat("-", 56))
	}
	return fmt.Sprintf("\n%s\n", strings.Repeat("-", ancho))
}

func ImprimirCadena(c Cadena, s *Salida) error {
	_, err := s.WriteString(c.S())

	if err != nil {
		return err
	}
	s.Flush()
	return nil
}

func (s Salida) Fd() uintptr {
	return s.f.Fd()
}

func (e Entrada) Fd() uintptr {
	return e.f.Fd()
}

func (e consola) Read(p []byte) (n int, err error) {
	return e.Entrada.Read(p)
}
func (s consola) Write(p []byte) (n int, err error) {
	return s.Salida.Write(p)
}
