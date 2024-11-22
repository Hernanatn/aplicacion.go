package color

type Color = string
type ColorFuente = Color
type ColorFondo = Color

const Resetear = Color("\033[0m")

const ResetearFuente = ColorFuente(Resetear)
const RojoFuente = ColorFuente("\033[31m")
const VerdeFuente = ColorFuente("\033[32m")
const AmarilloFuente = ColorFuente("\033[33m")
const AzulFuente = ColorFuente("\033[34m")
const MagentaFuente = ColorFuente("\033[35m")
const CyanFuente = ColorFuente("\033[36m")
const GrisFuente = ColorFuente("\033[38;2;70;70;70m")
const BlancoFuente = ColorFuente("\033[97m")

const ResetearFondo = ColorFondo(Resetear)
const NegroFondo = ColorFondo("\033[40m")
const GrisFondo = ColorFondo("\033[48;2;70;70;70m")
const RojoFondo = ColorFondo("\033[41m")
const VerdeFondo = ColorFondo("\033[42m")
const AmarilloFondo = ColorFondo("\033[43m")
const AzulFondo = ColorFondo("\033[44m")
const MagentaFondo = ColorFondo("\033[45m")
const CyanFondo = ColorFondo("\033[46m")
const BlancoFondo = ColorFondo("\033[47m")
