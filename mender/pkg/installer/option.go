package installer

type Option func(i *Installer)

func WithCommandRunner(c CommandRunner) Option {
	return func(i *Installer) {
		i.CommandRunner = c
	}
}
