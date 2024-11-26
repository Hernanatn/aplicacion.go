package multimenu

import (
	"errors"
	"slices"

	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
	"github.com/hernanatn/aplicacion.go/consola/teclado"
	"github.com/hernanatn/aplicacion.go/menu"
)

type Accion = menu.Accion
type Opcion = menu.Opcion

type MultiMenu struct {
	Opciones []*Opcion
	Consola  consola.Consola

	cursor        rune
	enfocada      int
	seleccionadas []int
	largo         int
	debeCerrar    bool
}

func NuevoMultiMenu(con consola.Consola, cur rune) *MultiMenu {
	return &MultiMenu{
		[]*Opcion{},
		con,
		cur,
		0,
		[]int{},
		0,
		false,
	}
}

func (m *MultiMenu) RegistrarOpcion(o *Opcion) *MultiMenu {
	m.Opciones = append(m.Opciones, o)
	m.largo++
	return m
}

func (m MultiMenu) borrarMenu() *MultiMenu {
	m.Consola.BorrarLinea()
	for i := 0; i < m.largo; i++ {
		m.Consola.EscribirBytes(teclado.CURSOR_PRINCIPIO_ANTERIOR)
		m.Consola.BorrarLinea()
	}
	return &m
}

func (m MultiMenu) imprimirOpcion(o *Opcion, enfocada bool) {

	var esEnfocada cadena.Cadena = " "
	var esSeleccionada cadena.Cadena = " "
	var formatos []cadena.Formateador = make([]cadena.Formateador, 0)
	var c cadena.Cadena

	if enfocada {
		esEnfocada = cadena.Cadena(string(m.cursor))
		formatos = append(formatos, cadena.Negrita)
	}
	if slices.Contains(m.seleccionadas, slices.Index(m.Opciones, o)) {
		esSeleccionada = cadena.Cadena("■")
		formatos = append(formatos, cadena.Negrita)
		formatos = append(formatos, cadena.Coloreador(color.CyanFondo))
	}

	c = (esEnfocada + esSeleccionada + "\t" + cadena.Cadena(o.Nombre)).Formatear(formatos...)
	m.Consola.ImprimirLinea(c)

}

func (m MultiMenu) renderizar() *MultiMenu {
	for i, o := range m.Opciones {
		m.imprimirOpcion(o, i == m.enfocada)
	}
	return &m
}

func (m MultiMenu) DebeCerrar() bool {
	return m.debeCerrar
}

func (m *MultiMenu) abrir() {
	m.debeCerrar = false
	m.enfocada = 0
	m.renderizar()
}
func (m *MultiMenu) Correr() ([]*Opcion, error) {
	m.abrir()
	var opcionesSeleccionadas []*Opcion
	var errores []error
	for !m.DebeCerrar() {
		var tecla []byte = make([]byte, 3)
		_, err := m.Consola.LeerTecla(&tecla)
		if err != nil {
			m.Consola.ImprimirLinea(cadena.Cadena(cadena.Error("menu.go / 93 > m.Consola.LeerTecla(&tecla)", err)))
			errores = append(errores, err)
		}
		switch tecla[0] {
		case teclado.CTRL_C:
			errores = append(errores, errors.New("programa cerrado por el usuario ^C"))
			m.debeCerrar = true

		case teclado.ESPACIO:
			i := slices.Index(m.seleccionadas, m.enfocada)
			if i >= 0 {
				m.seleccionadas = slices.Concat(m.seleccionadas[:i], m.seleccionadas[i+1:])
			} else {
				m.seleccionadas = append(m.seleccionadas, m.enfocada)
			}
			opcionesSeleccionadas = append(opcionesSeleccionadas, m.Opciones[m.enfocada])

		case teclado.ENTER:
			m.debeCerrar = true

		case teclado.ESC:
			switch tecla[1] {
			case teclado.CSI:
				switch tecla[2] {
				case teclado.FLECHA_ARRIBA[2]:
					if m.enfocada > 0 {
						m.enfocada--
					} else {
						m.enfocada = len(m.Opciones) - 1
					}
				case teclado.FLECHA_ABAJO[2]:
					if m.enfocada < len(m.Opciones)-1 {
						m.enfocada++
					} else {
						m.enfocada = 0
					}
				}
			default:

				m.borrarMenu()
				m.abrir()
				//m.debeCerrar = true
			}
		default:
		}
		m.borrarMenu()
		m.renderizar()
	}

	for _, i := range m.seleccionadas {
		opcionesSeleccionadas = append(opcionesSeleccionadas, m.Opciones[i])
	}
	return opcionesSeleccionadas, errors.Join(errores...)
}
