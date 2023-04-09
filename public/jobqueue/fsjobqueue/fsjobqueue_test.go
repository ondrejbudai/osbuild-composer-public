package fsjobqueue_test

import (
	"testing"

	"github.com/ondrejbudai/osbuild-composer-public/pkg/jobqueue"
	"github.com/stretchr/testify/require"

	"github.com/ondrejbudai/osbuild-composer-public/public/jobqueue/fsjobqueue"
	"github.com/ondrejbudai/osbuild-composer-public/public/jobqueue/jobqueuetest"
)

func TestJobQueueInterface(t *testing.T) {
	jobqueuetest.TestJobQueue(t, func() (jobqueue.JobQueue, func(), error) {
		dir := t.TempDir()
		q, err := fsjobqueue.New(dir)
		if err != nil {
			return nil, nil, err
		}
		stop := func() {
		}
		return q, stop, nil
	})
}

func TestNonExistant(t *testing.T) {
	q, err := fsjobqueue.New("/non-existant-directory")
	require.Error(t, err)
	require.Nil(t, q)
}
