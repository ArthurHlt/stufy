package loading

import (
	"github.com/briandowns/spinner"
	"time"
)

type SpinCliLoader struct {
	spinner *spinner.Spinner
}

func (s SpinCliLoader) Start(ldMsg *LoadMsg) {
	s.spinner.Prefix = ldMsg.InitMsg
	s.spinner.FinalMSG = ldMsg.FinishMsg
	s.spinner.Start()
}

func (s SpinCliLoader) Stop(_ *LoadMsg) {
	s.spinner.Stop()
}

func NewSpinCliLoader() *SpinCliLoader {
	return &SpinCliLoader{
		spinner: spinner.New(spinner.CharSets[35], 100*time.Millisecond),
	}
}
