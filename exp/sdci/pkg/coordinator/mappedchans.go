package coordinator

import (
	"sync"

	"cirello.io/exp/sdci/pkg/models"
)

// TODO: remove when using proper queues.
type mappedChans struct {
	m sync.Map // map of models.Build.RepoFullName to chan *models.Build
}

func (mc *mappedChans) ch(repoFullName string) chan *models.Build {
	ch := make(chan *models.Build, 10)
	foundCh, _ := mc.m.LoadOrStore(repoFullName, ch)
	return foundCh.(chan *models.Build)
}
