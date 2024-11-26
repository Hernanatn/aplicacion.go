package teclado

const (
	NUL byte = 0b00000000 // ^@	\0	Null
	SOH byte = 0b00000001 // ^A		Inicio de Encabezado
	STX byte = 0b00000010 // ^B		Inicio de Texto
	ETX byte = 0b00000011 // ^C		Fin de Texto
	EOT byte = 0b00000100 // ^D		Fin de Transmisión
	ENQ byte = 0b00000101 // ^E		Consulta
	ACK byte = 0b00000110 // ^F		Reconocimiento
	BEL byte = 0b00000111 // ^G	\a	Bell
	BS  byte = 0b00001000 // ^H	\b	Retroceso
	HT  byte = 0b00001001 // ^I	\t	Tab Horizontal
	LF  byte = 0b00001010 // ^J	\n	Nueva Linea
	VT  byte = 0b00001011 // ^K	\v	Tab Vertical
	FF  byte = 0b00001100 // ^L	\f	Form Feed
	CR  byte = 0b00001101 // ^M	\r	Retorno de Carrete
	SO  byte = 0b00001110 // ^N		Sale de Shift
	SI  byte = 0b00001111 // ^O		Entra en Shift
	DLE byte = 0b00010000 // ^P		Escape de enlace de Data
	DC1 byte = 0b00010001 // ^Q		Control de Dispositivo 1 (XON)
	DC2 byte = 0b00010010 // ^R		Control de Dispositivo 2
	DC3 byte = 0b00010011 // ^S		Control de Dispositivo 3 (XOFF)
	DC4 byte = 0b00010100 // ^T		Control de Dispositivo 4
	NAK byte = 0b00010101 // ^U		Reconocimiento Negativo
	SYN byte = 0b00010110 // ^V		Pausa Sincrónica
	ETB byte = 0b00010111 // ^W		Bloque de Fin de Transmición
	CAN byte = 0b00011000 // ^X		Cancelar
	EM  byte = 0b00011001 // ^Y		Fin del Medio
	SUB byte = 0b00011010 // ^Z		Substitute
	ESC byte = 0b00011011 // ^[	\e	Escape
	FS  byte = 0b00011100 // ^\		Sepearador de Archivo
	GS  byte = 0b00011101 // ^]		Sepearador de Grupo
	RS  byte = 0b00011110 // ^^		Sepearador de Registro
	US  byte = 0b00011111 // ^_		Sepearador de Unidara
	DEL byte = 0b01111111 // ^?		Borrar

	CSI byte = 0b01011011 // [

	/*
		ESC[{line};{column}H	// moves cursor to line #, column #
		ESC[{line};{column}f	// moves cursor to line #, column #

		ESC[#G	// moves cursor to column #
		ESC[6n	// request cursor position (reports as ESC[#;#R)
		ESC M	// moves cursor one line up, scrolling if needed
		ESC 7	// save cursor position (DEC)
		ESC 8	// restores the cursor to the last saved position (DEC)
		ESC[s	// save cursor position (SCO)
		ESC[u	// restores the cursor to the last saved position (SCO)
	*/
)

const (
	CTRL_C byte = ETX // ^C		Fin de Texto
	ENTER  byte = CR  // ^M	\r	Retorno de Carrete
)

const ( // ASCII
	A       byte = 0b01000001
	B       byte = 0b01000010
	C       byte = 0b01000011
	D       byte = 0b01000100
	E       byte = 0b01000101
	F       byte = 0b01000110
	G       byte = 0b01000111
	H       byte = 0b01001000
	I       byte = 0b01001001
	J       byte = 0b01001010
	K       byte = 0b01001011
	L       byte = 0b01001100
	M       byte = 0b01001101
	N       byte = 0b01001110
	O       byte = 0b01001111
	P       byte = 0b01010000
	Q       byte = 0b01010001
	R       byte = 0b01010010
	S       byte = 0b01010011
	T       byte = 0b01010100
	U       byte = 0b01010101
	V       byte = 0b01010110
	W       byte = 0b01010111
	X       byte = 0b01011000
	Y       byte = 0b01011001
	Z       byte = 0b01011010
	a       byte = 0b01100001
	b       byte = 0b01100010
	c       byte = 0b01100011
	d       byte = 0b01100100
	e       byte = 0b01100101
	f       byte = 0b01100110
	g       byte = 0b01100111
	h       byte = 0b01101000
	i       byte = 0b01101001
	j       byte = 0b01101010
	k       byte = 0b01101011
	l       byte = 0b01101100
	m       byte = 0b01101101
	n       byte = 0b01101110
	o       byte = 0b01101111
	p       byte = 0b01110000
	q       byte = 0b01110001
	r       byte = 0b01110010
	s       byte = 0b01110011
	t       byte = 0b01110100
	u       byte = 0b01110101
	v       byte = 0b01110110
	w       byte = 0b01110111
	x       byte = 0b01111000
	y       byte = 0b01111001
	z       byte = 0b01111010
	ESPACIO byte = 0b00100000

	//Extendido
	CUADRADO = 0b10010110100000
)

var (
	CURSOR_CASA                []byte = []byte{ESC, CSI, H} // Debieran ser constantes. No mutar!
	CURSOR_PRINCIPIO_ANTERIOR  []byte = []byte{ESC, CSI, F} // Debieran ser constantes. No mutar!
	CURSOR_PRINCIPIO_SIGUIENTE []byte = []byte{ESC, CSI, E} // Debieran ser constantes. No mutar!
	FLECHA_ARRIBA              []byte = []byte{ESC, CSI, A} // Debieran ser constantes. No mutar!
	FLECHA_ABAJO               []byte = []byte{ESC, CSI, B} // Debieran ser constantes. No mutar!
	FLECHA_DERECHA             []byte = []byte{ESC, CSI, C} // Debieran ser constantes. No mutar!
	FLECHA_IZQUIERDA           []byte = []byte{ESC, CSI, D} // Debieran ser constantes. No mutar!
)
