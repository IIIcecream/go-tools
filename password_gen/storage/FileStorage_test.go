package storage

import (
	"os"
	"testing"
)

func TestSavePassword(t *testing.T) {
	fs, e := New()
	if e != nil {
		t.Errorf("create file storage obj failed: %s", e)
	}
	fs.storagePath = "./test.enc"
	os.RemoveAll(fs.storagePath)

	password := "test1"
	fs.SavePassword("test1", password)
	entries, e := fs.FindPassword("test1")
	if e != nil {
		t.Errorf("find password failed: %s", e)
	}
	if len(entries) != 1 {
		t.Errorf("find entries size expected: %d, get: %d", len(entries), 1)
	}
	if entries[0].Password != password {
		t.Errorf("pasword expected: %s, get: %s", password, entries[0].Password)
	}
}
