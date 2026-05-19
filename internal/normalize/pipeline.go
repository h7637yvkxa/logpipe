package normalize

// Pipeline reads raw log lines from an input channel, normalizes each line
// using the provided Normalizer, and emits valid Entry values on the output channel.
// Lines that fail to parse are silently dropped (parse errors are sent to ErrCh).
type Pipeline struct {
	norm   *Normalizer
	In     <-chan string
	Out    chan *Entry
	ErrCh  chan error
}

// NewPipeline constructs a Pipeline for the given service and input channel.
func NewPipeline(service string, in <-chan string) *Pipeline {
	return &Pipeline{
		norm:  New(service),
		In:    in,
		Out:   make(chan *Entry, 64),
		ErrCh: make(chan error, 16),
	}
}

// Run starts processing lines from In until the channel is closed.
// It closes Out and ErrCh when done. Intended to be run in a goroutine.
func (p *Pipeline) Run() {
	defer close(p.Out)
	defer close(p.ErrCh)

	for line := range p.In {
		if line == "" {
			continue
		}
		entry, err := p.norm.Normalize(line)
		if err != nil {
			select {
			case p.ErrCh <- err:
			default:
			}
			continue
		}
		p.Out <- entry
	}
}
