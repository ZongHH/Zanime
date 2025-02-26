package bootstrap

type Bootstrap struct {
	Container *Container
}

func NewBootstrap(configPath string) *Bootstrap {
	container := NewContainer(configPath)
	return &Bootstrap{
		Container: container,
	}
}

func (b *Bootstrap) Start() {
	b.Container.Interfaces.Start(b.Container.Config.GetGRPCAddr())
	b.Container.Consumers.Start()
}

func (b *Bootstrap) Stop() {
	b.Container.Bases.Close()
	b.Container.Consumers.Close()
	b.Container.Interfaces.Close()
}
