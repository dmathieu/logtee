package handler

import (
	"testing"
	"time"

	"github.com/dmathieu/logtee/logtee"
	"github.com/dmathieu/logtee/metrics"
	"github.com/stretchr/testify/assert"
)

func TestHandlerHandle(t *testing.T) {
	handler := Handler{}
	log := logtee.LogLine{
		Time: time.Now(),
		Path: "/",
	}

	err := handler.Handle(log)
	assert.NoError(t, err)

	snap, err := metrics.RegistryProvider().Collect(time.Now().Add(time.Second * -1))
	assert.NoError(t, err)

	assert.Equal(t, int64(1), snap.FetchCounter("requests.total"))
	assert.Equal(t, int64(1), snap.FetchCounter("requests.section.home"))
}
