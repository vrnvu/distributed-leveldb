package store

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var tempStorageDir = "store_test"

func assertOpen(s *Store, t *testing.T, localID string) {
	if err := s.Open(true, localID); err != nil {
		t.Fatalf("failed to open store: %s", err)
	}
}

func assertSet(s *Store, t *testing.T, key string, value string) {
	if err := s.Set(key, value); err != nil {
		t.Fatalf("failed to set key: %s", err.Error())
	}
}

func assertGet(s *Store, t *testing.T, key string, value string) {
	result, err := s.Get(key)
	if err != nil {
		t.Fatalf("failed to get key: %s", err.Error())
	}
	if result != value {
		t.Fatalf("key has wrong value: %s", result)
	}
}

// Test_StoreOpen tests that the store can be opened.
func Test_StoreOpen(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)

	raftBind := "127.0.0.1:0"
	raftDir := tmpDir

	s, err := New(false, tmpDir, raftDir, raftBind)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer s.db.Close()

	s.RaftBind = "127.0.0.1:0"
	s.RaftDir = tmpDir
	if s == nil {
		t.Fatalf("failed to create store")
	}

	assertOpen(s, t, "node0")
}

// Test_StoreOpenSingleNode tests that a command can be applied to the log
func Test_StoreOpenSingleNode(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)

	raftBind := "127.0.0.1:0"
	raftDir := tmpDir

	s, err := New(false, tmpDir, raftDir, raftBind)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if s == nil {
		t.Fatalf("failed to create store")
	}

	assertOpen(s, t, "node0")

	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)

	assertSet(s, t, "foo", "bar")

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)
	assertGet(s, t, "foo", "bar")
}

// Test_StoreInMemOpenSingleNode tests that a command can be applied to the log
// stored in RAM.
func Test_StoreInMemOpenSingleNode(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)

	raftBind := "127.0.0.1:0"
	raftDir := tmpDir

	s, err := New(false, tmpDir, raftDir, raftBind)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if s == nil {
		t.Fatalf("failed to create store")
	}

	assertOpen(s, t, "node0")

	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)

	assertSet(s, t, "foo", "bar")

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)

	assertGet(s, t, "foo", "bar")
	if err := s.Delete("foo"); err != nil {
		t.Fatalf("failed to delete key: %s", err.Error())
	}

}

func Test_LeveldbRawAccess(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)

	raftBind := "127.0.0.1:0"
	raftDir := tmpDir

	s, err := New(false, tmpDir, raftDir, raftBind)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if s == nil {
		t.Fatalf("failed to create store")
	}

	assertOpen(s, t, "node0")

	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)

	assertSet(s, t, "foo", "bar")

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)

	assertGet(s, t, "foo", "bar")

	if err := s.Delete("foo"); err != nil {
		t.Fatalf("failed to delete key: %s", err.Error())
	}

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)
	_, err = s.Get("foo")
	if err != leveldb.ErrNotFound {
		t.Fatalf("failed this test should have failed after delete")
	}

}

func Test_LeveldbRaw(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)
	db, err := leveldb.OpenFile(tmpDir, nil)
	defer db.Close()
	if err != nil {
		t.Fatalf("error creating leveldb: %s", err.Error())
	}

	if err := db.Put([]byte("foo"), []byte("bar"), nil); err != nil {
		t.Fatalf("failed to set key: %s", err.Error())
	}
	value, err := db.Get([]byte("foo"), nil)

	actual := string(value)
	if err != nil {
		t.Fatalf("failed to get key: %s", err.Error())
	}
	if actual != "bar" {
		t.Fatalf("key has wrong value: %s", actual)
	}
}

// stored in RAM.
func Test_MultiplePut(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", tempStorageDir)
	defer os.RemoveAll(tmpDir)

	raftBind := "127.0.0.1:0"
	raftDir := tmpDir

	s, err := New(false, tmpDir, raftDir, raftBind)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if s == nil {
		t.Fatalf("failed to create store")
	}

	assertOpen(s, t, "node0")

	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)
	assertSet(s, t, "foo", "bar")
	assertSet(s, t, "aaa", "bar")
	assertSet(s, t, "bbb", "bar")
	assertSet(s, t, "ccccc", "bar")

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)

	assertGet(s, t, "foo", "bar")
	assertGet(s, t, "aaa", "bar")
	assertGet(s, t, "bbb", "bar")
	assertGet(s, t, "ccccc", "bar")
}
