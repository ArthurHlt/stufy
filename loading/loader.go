package loading

type LoadMsg struct {
	InitMsg   string
	FinishMsg string
	cb        func()
}

type Loader interface {
	Start(*LoadMsg)
	Stop(*LoadMsg)
}

var loadMessage *LoadMsg

var fLoader Loader

var hide bool

func SetLoader(loader Loader) {
	fLoader = loader
}

func Hide() {
	hide = true
}

func Show() {
	hide = false
}

func Start(init, finish string, cb func()) {
	if hide {
		return
	}
	loadMessage = &LoadMsg{
		InitMsg:   init,
		FinishMsg: finish,
		cb:        cb,
	}
	if fLoader != nil {
		fLoader.Start(loadMessage)
	}
}

func Stop() {
	if hide {
		return
	}
	if fLoader != nil {
		fLoader.Stop(loadMessage)
	}
}
