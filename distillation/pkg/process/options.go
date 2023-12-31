package process

type Option func(p *Process)

func WithClock(c Clock) Option {
	return func(p *Process) {
		p.clock = c
	}
}
