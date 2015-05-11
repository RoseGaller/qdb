// Copyright 2015 Reborndb Org. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package store

import (
	"os"
	"testing"
	"time"

	"github.com/reborndb/go/log"
	"github.com/reborndb/go/testing/assert"
	"github.com/reborndb/qdb/pkg/engine/rocksdb"
)

var (
	testbl *Binlog
)

func reinit() {
	if testbl != nil {
		testbl.Close()
		testbl = nil
	}
	const path = "/tmp/test_qdb/binlog/testdb-rocksdb"
	if err := os.RemoveAll(path); err != nil {
		log.PanicErrorf(err, "remove '%s' failed", path)
	} else {
		conf := rocksdb.NewDefaultConfig()
		if testdb, err := rocksdb.Open(path, conf, false); err != nil {
			log.PanicError(err, "open rocksdb failed")
		} else {
			testbl = New(testdb)
		}
	}
}

func init() {
	reinit()
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func checkerror(t *testing.T, err error, exp bool) {
	if err != nil || !exp {
		reinit()
	}
	assert.ErrorIsNil(t, err)
	assert.Must(t, exp)
}

func checkcompact(t *testing.T) {
	err := testbl.CompactAll()
	checkerror(t, err, true)
}

func checkempty(t *testing.T) {
	it := testbl.getIterator()
	it.SeekToFirst()
	empty, err := !it.Valid(), it.Error()
	testbl.putIterator(it)
	checkerror(t, err, empty)
}

func sleepms(n int) {
	time.Sleep(time.Millisecond * time.Duration(n))
}
