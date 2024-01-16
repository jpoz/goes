package goes

type Mode int

const (
	ModeBuild Mode = iota
	ModeEmbedded
)

type Options struct {
	Mode   Mode
	Logger any
}
