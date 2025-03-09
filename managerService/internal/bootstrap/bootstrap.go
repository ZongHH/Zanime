package bootstrap

type Bootstrap struct {
	Container *Container
}

func NewBootstrap(container *Container) *Bootstrap {
	return &Bootstrap{Container: container}
}

func (b *Bootstrap) Run() {
	b.Container.setAPIRouters(b.Container.Handlers)
	b.Container.Controller.Run()
}

func (b *Bootstrap) Shutdown() {
	b.Container.Controller.Shutdown()
}
