package test

import (
	"cmd/safe"
	"testing"
	"time"
)

var safePath = "/home/bdw/Desktop/go.safe"
var password = "secr3tP@ssword"

var testSafe = &safe.Safe{time.Now().Unix(), "safe_test", nil}
var testCredential1 = safe.Credential{"github", "bndw", "p@ssw3rd", time.Now().Unix()}

func TestSaveSafe(t *testing.T) {
	err := testSafe.Save(safePath, password)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadSafe(t *testing.T) {
	_, err := safe.Load(safePath, password)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateCredential(t *testing.T) {
	mySafe, err := safe.Load(safePath, password)
	if err != nil {
		t.Error(err)
	}

	err = mySafe.AddCredential(testCredential1)
	if err != nil {
		t.Error(err)
	}

	err = mySafe.Save(safePath, password)
	if err != nil {
		t.Error(err)
	}
}

func TestCannotCreateDuplicateCredential(t *testing.T) {
	mySafe, err := safe.Load(safePath, password)
	if err != nil {
		t.Error(err)
	}

	err = mySafe.AddCredential(testCredential1)
	if err == nil {
		t.Error(err)
	}
}

func TestGetCredential(t *testing.T) {
	mySafe, err := safe.Load(safePath, password)
	if err != nil {
		t.Error(err)
	}

	_, err = mySafe.GetCredential(testCredential1.Alias)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveCredential(t *testing.T) {
	mySafe, err := safe.Load(safePath, password)
	if err != nil {
		t.Error(err)
	}

	err = mySafe.RemoveCredential(testCredential1.Alias)
	if err != nil {
		t.Error(err)
	}
}

func TestGeneratePassword(t *testing.T) {
	length := 50

	_, err := safe.GeneratePassword(length)
	if err != nil {
		t.Error(err)
	}
}
