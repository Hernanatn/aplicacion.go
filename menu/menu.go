package menu

import (
	"errors"

	"github.com/hernanatn/aplicacion.go/comando"
	"github.com/hernanatn/aplicacion.go/consola"
	"github.com/hernanatn/aplicacion.go/consola/cadena"
	"github.com/hernanatn/aplicacion.go/consola/color"
	"github.com/hernanatn/aplicacion.go/consola/teclado"
)

type Accion = comando.Accion

type Opcion struct {
	Nombre string
	Accion Accion
}

type Menu struct {
	Opciones []*Opcion
	Consola  consola.Consola

	cursor       rune
	seleccionada int
	largo        int
	debeCerrar   bool
}

func NuevoMenu(con consola.Consola, cur rune) *Menu {
	return &Menu{
		[]*Opcion{},
		con,
		cur,
		0,
		0,
		false,
	}
}

func (m *Menu) RegistrarOpcion(o *Opcion) *Menu {
	m.Opciones = append(m.Opciones, o)
	m.largo++
	return m
}

func (m Menu) borrarMenu() *Menu {
	m.Consola.BorrarLinea()
	for i := 0; i < m.largo; i++ {
		m.Consola.EscribirBytes(teclado.CURSOR_PRINCIPIO_ANTERIOR)
		m.Consola.BorrarLinea()
	}
	return &m
}

func (m Menu) imprimirOpcion(o *Opcion, seleccionada bool) {

	var c cadena.Cadena
	if seleccionada {
		c = cadena.Cadena(string(m.cursor) + "\t" + o.Nombre).Negrita().Colorear(color.CyanFondo)
	} else {
		c = cadena.Cadena(" \t" + o.Nombre)
	}

	m.Consola.ImprimirLinea(c)

}

func (m Menu) renderizar() *Menu {
	for i, o := range m.Opciones {
		m.imprimirOpcion(o, i == m.seleccionada)
	}
	return &m
}

func (m Menu) DebeCerrar() bool {
	return m.debeCerrar
}

func (m *Menu) abrir() {
	m.debeCerrar = false
	m.seleccionada = 0
	m.renderizar()
}
func (m *Menu) Correr() (*Opcion, error) {
	m.abrir()
	var opcionRelevante *Opcion
	var errores []error
	for !m.DebeCerrar() {
		var tecla []byte = make([]byte, 3)
		_, err := m.Consola.LeerTecla(&tecla)
		if err != nil {
			m.Consola.ImprimirError("menu.go / 93 > m.Consola.LeerTecla(&tecla)", err)
			errores = append(errores, err)
		}
		switch tecla[0] {
		case teclado.CTRL_C:
			errores = append(errores, errors.New("programa cerrado por el usuario ^C"))
			m.debeCerrar = true

		case teclado.ENTER:
			opcionRelevante = m.Opciones[m.seleccionada]
			m.debeCerrar = true

		case teclado.ESC:
			switch tecla[1] {
			case teclado.CSI:
				switch tecla[2] {
				case teclado.FLECHA_ARRIBA[2]:
					if m.seleccionada > 0 {
						m.seleccionada--
					} else {
						m.seleccionada = len(m.Opciones) - 1
					}
				case teclado.FLECHA_ABAJO[2]:
					if m.seleccionada < len(m.Opciones)-1 {
						m.seleccionada++
					} else {
						m.seleccionada = 0
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
	return opcionRelevante, errors.Join(errores...)
}
