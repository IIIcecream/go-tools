package storage

import "testing"

func TestSavePassword(t *testing.T) {
	fs, e := New()
	if e != nil {
		t.Errorf("create file storage obj failed: %s", e)
	}
	password := "test1"
	fs.SavePassword("test1", password)
	entry, e := fs.FindPassword("test1")
	if e != nil {
		t.Errorf("find password failed: %s", e)
	}
	if entry.Password != password {
		t.Errorf("pasword expected: %s, get: %s", password, entry.Password)
	}
}
