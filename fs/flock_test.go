package fs

import (
	"io/ioutil"
	"testing"
)

func TestFLocker(t *testing.T) {

	ioutil.WriteFile(".lock", []byte("FLocker"), 0644)

	FLocker, err := NewFLocker(".lock")

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {

		if err := FLocker.Lock(); err != nil {
			t.Fatal(err)
		}

		if err := FLocker.Unlock(); err != nil {
			t.Fatal(err)
		}
	}

}
