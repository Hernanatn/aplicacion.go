package aplicacion_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/hernanatn/aplicacion.go"
	"github.com/hernanatn/aplicacion.go/comando"
	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNuevaAplicacion prueba la creación y configuración básica de una aplicación
func TestNuevaAplicacion(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{"-v", "--version"},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	assert.NotNil(t, app)
	assert.Equal(t, "app-prueba", app.DevolverNombre())
	assert.False(t, app.EsOculto())
}

// TestRegistroDeComandos prueba la funcionalidad de registro de comandos
func TestRegistroDeComandos(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	ejecutado := false
	cmdPrueba := comando.NuevoComando(
		"cmd-prueba",
		"uso de prueba",
		[]string{"cp"},
		"Comando de Prueba",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			ejecutado = true
			return nil, comando.EXITO, nil
		},
		[]string{},
	)

	app.RegistrarComando(cmdPrueba)

	_, codigo, err := app.Ejecutar(nil, "cmd-prueba")
	assert.Nil(t, err)
	assert.Equal(t, comando.EXITO, codigo)
	assert.True(t, ejecutado)
}

// TestAliasDeComandos prueba que los alias de comandos funcionen correctamente
func TestAliasDeComandos(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	ejecutado := false
	cmdPrueba := comando.NuevoComando(
		"cmd-prueba",
		"uso de prueba",
		[]string{"cp", "prueba"},
		"Comando de Prueba",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			ejecutado = true
			return nil, comando.EXITO, nil
		},
		[]string{},
	)

	app.RegistrarComando(cmdPrueba)

	// Probar comando principal
	ejecutado = false
	_, _, _ = app.Ejecutar(nil, "cmd-prueba")
	assert.True(t, ejecutado)

	// Probar alias
	ejecutado = false
	_, _, _ = app.Ejecutar(nil, "cp")
	assert.True(t, ejecutado)

	ejecutado = false
	_, _, _ = app.Ejecutar(nil, "prueba")
	assert.True(t, ejecutado)
}

// TestOpcionesDeComandos prueba el manejo de opciones y parámetros de comandos
func TestOpcionesDeComandos(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	var opcionesCapturadas comando.Opciones
	var parametrosCapturados comando.Parametros

	cmdPrueba := comando.NuevoComando(
		"cmd-prueba",
		"uso de prueba",
		[]string{},
		"Comando de Prueba",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			opcionesCapturadas = opt
			parametrosCapturados = params
			return nil, comando.EXITO, nil
		},
		[]string{"--verbose"},
	)

	app.RegistrarComando(cmdPrueba)

	_, _, _ = app.Ejecutar(nil, "cmd-prueba", "--verbose", "--output", "archivo.txt")

	assert.Contains(t, opcionesCapturadas, "--verbose")
	t.Log(parametrosCapturados)
	assert.Contains(t, parametrosCapturados["--output"], "archivo.txt")
}

func TestConsolaEscribeEnWriter(t *testing.T) {
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	con := consola.NuevaConsola(os.Stdin, w)
	con.EscribirLinea("Test123")
	con.Imprimir()
	w.Close()

	out, _ := io.ReadAll(r)
	assert.Contains(t, string(out), "Test123")
}

// TestSalidaConsola prueba la funcionalidad de salida en consola
func TestSalidaConsola(t *testing.T) {
	lector, escritor, err := os.Pipe()
	require.NoError(t, err)
	defer lector.Close()
	defer escritor.Close()
	con := consola.NuevaConsola(os.Stdin, escritor)

	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		con,
	)

	cmdPrueba := comando.NuevoComando(
		"cmd-prueba",
		"uso de prueba",
		[]string{},
		"Comando de Prueba",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			con.EscribirLinea("Hola Mundo")
			con.Imprimir()
			return nil, comando.EXITO, nil
		},
		[]string{},
	)

	app.RegistrarComando(cmdPrueba)

	// Ejecutamos el comando
	_, _, err = app.Ejecutar(con, "cmd-prueba")
	require.NoError(t, err)

	escritor.Close()
	salida, _ := io.ReadAll(lector)
	assert.Contains(t, string(salida), "Hola Mundo")
}

// TestComandosIntegrados prueba los comandos integrados como ayuda y salida
func TestComandosIntegrados(t *testing.T) {
	var buf bytes.Buffer
	io.Copy(&buf, os.Stdout)
	con := consola.NuevaConsola(os.Stdin, os.Stdout)

	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		con,
	)

	// Probar comando de ayuda
	_, _, _ = app.Ejecutar(con, "ayuda")
	b := make([]byte, 100)
	app.FSalida().Read(b)
	t.Log((b))
	//assert.True(t, strings.Contains(b, "Ayuda"))
	//assert.True(t, strings.Contains(b, "Comandos:"))

	// Limpiar buffer
	buf.Reset()

	// Probar comando de salida
	_, _, _ = app.Ejecutar(con, "chau")
	assert.True(t, app.DebeCerrar())
}

// TestManejoDeErrores prueba el manejo de errores en los comandos
func TestManejoDeErrores(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	cmdPrueba := comando.NuevoComando(
		"cmd-prueba",
		"uso de prueba",
		[]string{},
		"Comando de Prueba",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			return nil, comando.ERROR, os.ErrNotExist
		},
		[]string{},
	)

	app.RegistrarComando(cmdPrueba)

	_, codigo, err := app.Ejecutar(nil, "cmd-prueba")
	assert.Equal(t, comando.ERROR, codigo)
	assert.Equal(t, os.ErrNotExist, err)
}

// TestComandosAnidados prueba la funcionalidad de comandos anidados
func TestComandosAnidados(t *testing.T) {
	app := aplicacion.NuevaAplicacion(
		"app-prueba",
		"uso de prueba",
		"Descripción de Prueba",
		[]string{},
		consola.NuevaConsola(os.Stdin, os.Stdout),
	)

	padreEjecutado := false
	hijoEjecutado := false

	cmdPadre := comando.NuevoComando(
		"padre",
		"uso del padre",
		[]string{},
		"Comando Padre",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			padreEjecutado = true
			return nil, comando.EXITO, nil
		},
		[]string{},
	)

	cmdHijo := comando.NuevoComando(
		"hijo",
		"uso del hijo",
		[]string{},
		"Comando Hijo",
		func(con comando.Consola, opt comando.Opciones, params comando.Parametros, args ...any) (any, comando.CodigoError, error) {
			hijoEjecutado = true
			return nil, comando.EXITO, nil
		},
		[]string{},
	)

	app.RegistrarComando(cmdPadre)
	cmdPadre.RegistrarComando(cmdHijo)

	_, _, _ = app.Ejecutar(nil, "padre", "hijo")
	assert.False(t, padreEjecutado)
	assert.True(t, hijoEjecutado)
}
