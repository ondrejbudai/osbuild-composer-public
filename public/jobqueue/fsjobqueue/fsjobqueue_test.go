package fsjobqueue_test

import (
	"context"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondrejbudai/osbuild-composer-public/public/jobqueue/fsjobqueue"
	"github.com/ondrejbudai/osbuild-composer-public/public/jobqueue/jobqueuetest"
	"github.com/ondrejbudai/osbuild-composer-public/pkg/jobqueue"
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

func TestJobQueueBadJSON(t *testing.T) {
	dir := t.TempDir()

	// Write a purposfully invalid JSON file into the queue
	err := os.WriteFile(path.Join(dir, "/4f1cf5f8-525d-46b7-aef4-33c6a919c038.json"), []byte("{invalid json content"), 0600)
	require.Nil(t, err)

	q, err := fsjobqueue.New(dir)
	require.Nil(t, err)
	require.NotNil(t, q)
}

func sortUUIDs(entries []uuid.UUID) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].String() < entries[j].String()
	})
}

func TestAllRootJobIDs(t *testing.T) {
	dir := t.TempDir()
	q, err := fsjobqueue.New(dir)
	require.Nil(t, err)
	require.NotNil(t, q)

	var rootJobs []uuid.UUID

	// root with no dependencies
	jidRoot1, err := q.Enqueue("oneRoot", nil, nil, "OneRootJob")
	require.Nil(t, err)
	rootJobs = append(rootJobs, jidRoot1)

	// root with 2 dependencies
	jid1, err := q.Enqueue("twoDeps", nil, nil, "TwoDepJobs")
	require.Nil(t, err)
	jid2, err := q.Enqueue("twoDeps", nil, nil, "TwoDepJobs")
	require.Nil(t, err)
	jidRoot2, err := q.Enqueue("twoDeps", nil, []uuid.UUID{jid1, jid2}, "TwoDepJobs")
	require.Nil(t, err)
	rootJobs = append(rootJobs, jidRoot2)

	// root with 2 dependencies, one shared with the previous root
	jid3, err := q.Enqueue("sharedDeps", nil, nil, "SharedDepJobs")
	require.Nil(t, err)
	jidRoot3, err := q.Enqueue("sharedDeps", nil, []uuid.UUID{jid1, jid3}, "SharedDepJobs")
	require.Nil(t, err)
	rootJobs = append(rootJobs, jidRoot3)

	sortUUIDs(rootJobs)
	roots, err := q.AllRootJobIDs(context.TODO())
	require.Nil(t, err)
	require.Greater(t, len(roots), 0)
	sortUUIDs(roots)
	require.Equal(t, rootJobs, roots)
}

func TestDeleteJob(t *testing.T) {
	dir := t.TempDir()
	q, err := fsjobqueue.New(dir)
	require.Nil(t, err)
	require.NotNil(t, q)

	// root with no dependencies
	jidRoot1, err := q.Enqueue("oneRoot", nil, nil, "OneRootJob")
	require.Nil(t, err)

	err = q.DeleteJob(context.TODO(), jidRoot1)
	require.Nil(t, err)
	jobs, err := q.AllRootJobIDs(context.TODO())
	require.Nil(t, err)
	require.Equal(t, 0, len(jobs))

	// root with 2 dependencies
	jid1, err := q.Enqueue("twoDeps", nil, nil, "TwoDepJobs")
	require.Nil(t, err)
	jid2, err := q.Enqueue("twoDeps", nil, nil, "TwoDepJobs")
	require.Nil(t, err)
	jidRoot2, err := q.Enqueue("twoDeps", nil, []uuid.UUID{jid1, jid2}, "TwoDepJobs")
	require.Nil(t, err)

	// root with 2 dependencies, one shared with the previous root
	jid3, err := q.Enqueue("sharedDeps", nil, nil, "SharedDepJobs")
	require.Nil(t, err)
	jidRoot3, err := q.Enqueue("sharedDeps", nil, []uuid.UUID{jid1, jid3}, "SharedDepJobs")
	require.Nil(t, err)

	// This should only remove jidRoot2 and jid2, leaving jidRoot3, jid1, jid3
	err = q.DeleteJob(context.TODO(), jidRoot2)
	require.Nil(t, err)
	jobs, err = q.AllRootJobIDs(context.TODO())
	require.Nil(t, err)
	require.Equal(t, 1, len(jobs))
	assert.Equal(t, []uuid.UUID{jidRoot3}, jobs)

	// This should remove the rest
	err = q.DeleteJob(context.TODO(), jidRoot3)
	require.Nil(t, err)
	jobs, err = q.AllRootJobIDs(context.TODO())
	require.Nil(t, err)
	require.Equal(t, 0, len(jobs))

	// Make sure all the jobs are deleted
	allJobs := []uuid.UUID{jidRoot1, jidRoot2, jidRoot3, jid1, jid2, jid3}
	for _, jobId := range allJobs {
		jobType, _, _, _, err := q.Job(jobId)
		assert.Error(t, err, jobType)
	}

	// root with 2 jobs depending on another (simulates Koji jobs)
	kojiOSTree, err := q.Enqueue("ostree", nil, nil, "KojiJob")
	require.Nil(t, err)
	kojiDepsolve, err := q.Enqueue("depsolve", nil, nil, "KojiJob")
	require.Nil(t, err)
	kojiManifest, err := q.Enqueue("manifest", nil, []uuid.UUID{kojiOSTree, kojiDepsolve}, "KojiJob")
	require.Nil(t, err)
	kojiInit, err := q.Enqueue("init", nil, nil, "KojiJob")
	require.Nil(t, err)
	kojiRoot, err := q.Enqueue("final", nil, []uuid.UUID{kojiInit, kojiManifest, kojiDepsolve}, "KojiJob")
	require.Nil(t, err)

	// Delete the koji job
	err = q.DeleteJob(context.TODO(), kojiRoot)
	require.Nil(t, err)
	jobs, err = q.AllRootJobIDs(context.TODO())
	require.Nil(t, err)
	require.Equal(t, 0, len(jobs))

	// Make sure all the jobs are deleted
	kojiJobs := []uuid.UUID{kojiRoot, kojiInit, kojiOSTree, kojiDepsolve, kojiManifest}
	for _, jobId := range kojiJobs {
		jobType, _, _, _, err := q.Job(jobId)
		assert.Error(t, err, jobType)
	}
}
