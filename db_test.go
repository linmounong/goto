package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func newTmpfileName() string {
	tmpfile, err := ioutil.TempFile("", "gototest")
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}

// Opens database before test and cleans up after it.
func dbTestClosure(f func()) func() {
	return func() {
		*dbFilename = newTmpfileName()
		defer os.Remove(*dbFilename)
		DbOpen()
		defer DbClose()
		f()
	}
}

func TestDb(t *testing.T) {
	dbTestClosure(func() {
		if _, err := DbUpdateOrCreate("user", "a", "url1"); err != nil {
			t.Fatal(err)
		}
		if info := DbFind("a"); info == nil || info.Url != "url1" {
			t.Error("can't find a")
		}

		if _, err := DbUpdateOrCreate("user", "a", "url2"); err != nil {
			t.Fatal(err)
		}
		if info := DbFind("a"); info == nil || info.Url != "url2" {
			t.Error("can't find a")
		}

		if _, err := DbUpdateOrCreate("user2", "a", "url3"); err == nil {
			t.Fatal("user2 should not modify a")
		}
		DbIncr("a")
		DbIncr("a")
		DbIncr("a")
		if info := DbFind("a"); info == nil || info.UseCount != 3 {
			t.Error("can't find a")
		}
		DbRemove("a")
		if info := DbFind("a"); info != nil {
			t.Error("found deleted a")
		}

		DbUpdateOrCreate("user2", "ba", "url1")
		DbUpdateOrCreate("user2", "bb", "url2")
		DbUpdateOrCreate("user2", "bc", "url3")
		if results := DbSearch("b", "", "", 1, 10); len(results) != 3 {
			t.Error("expecting 3 results", results)
		}
		if results := DbSearch("b", "", "", 1, 2); len(results) != 2 {
			t.Error("expecting 2 results", results)
		}
		if results := DbSearch("b", "", "", 2, 2); len(results) != 1 {
			t.Error("expecting 1 result", results)
		}
	})()
}
