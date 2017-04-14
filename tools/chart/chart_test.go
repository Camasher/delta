package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestBuild_Tsunami(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/tsunami-gauge.xml")
	if err != nil {
		t.Fatalf("error: unable to load test amplitudes file: %v", err)
	}

	cfgs, err := loadConfig("./testdata/chart-tsunami.yaml")
	if err != nil {
		t.Fatalf("error: unable to load test tsunami config file: %v", err)
	}

	tsunami, err := buildPages(cfgs, "./testdata")
	if err != nil {
		t.Fatalf("error: unable to build test tsunami file: %v", err)
	}

	b2, err := encodePages(tsunami)
	if err != nil {
		t.Fatalf("error: unable to encode test tsunami file: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Error("**** tsunami xml mismatch ****")

		f1, err := ioutil.TempFile(os.TempDir(), "tmp")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f1.Name())
		f1.Write(b1)

		f2, err := ioutil.TempFile(os.TempDir(), "tmp")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f2.Name())
		f2.Write(b2)

		cmd := exec.Command("diff", "-c", f1.Name(), f2.Name())
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		err = cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		defer cmd.Wait()
		diff, err := ioutil.ReadAll(stdout)
		if err != nil {
			t.Fatal(err)
		}
		t.Error(string(diff))
	}
}
