package external

import (
	"fmt"
	"os"
	"testing"
	"path/filepath"

)

// create config file if not exist

func CreateIfNotExist(fpath string) bool {
	dir, _ := filepath.Split(fpath)
	var out bool
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	file, err := os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error when Reading")
	}
	file.Close()
	out = true
	return out
}

func TestAddProfileCredentials(t *testing.T) {
	t.Run("TestAddProfileCredentials1", func(t *testing.T) {
		path := DefaultSharedCredentialsFilename()
		CreateIfNotExist(path)
		AccessIdTest := "AccessIdTest"
		SecretKeyTest := "SecretKeyTest"
		ProfileTest := "Addfirsttest1"
		c := &Credentials{Profile: &ProfileTest, AccessKeyId: &AccessIdTest, SecretAccessKey: &SecretKeyTest}
		got, _ := AddProfileCredentials(c)
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})
	t.Run("TestDeleteProfileCredentials", func(t *testing.T) {
		got, _ := DeleteProfileCredentials("Addfirsttest1")
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}

// func TestDeleteProfileCredentials(t *testing.T) {
// 	t.Run("TestDeleteProfileCredentials", func(t *testing.T) {
// 		got, _ := DeleteProfileCredentials("Addfirsttest1")
// 		want := true
// 		if got != want {
// 			t.Errorf("Got: %v - want: %v", got, want)
// 		}
// 	})

// }

func TestAddProfileConfig(t *testing.T) {
	t.Run("TestAddProfileConfig", func(t *testing.T) {
		path := DefaultSharedConfigFilename()
		CreateIfNotExist(path)
		ProfiletTest := "firstconfig1"
		RegionTest := "us-east-1"
		OutputTest := "json"
		c := &CredentialsConfig{Profile: &ProfiletTest, Region: &RegionTest, Output: &OutputTest}
		got, _ := AddProfileConfig(c)
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})
	t.Run("TestDeleteprofileconfig", func(t *testing.T) {
		got, _ := DeleteProfileConfig("firstconfig1")
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}

// func TestDeleteProfileConfig(t *testing.T) {
// 	t.Run("TestDeleteprofileconfig", func(t *testing.T) {
// 		got, _ := DeleteProfileConfig("firstconfig1")
// 		want := true
// 		if got != want {
// 			t.Errorf("Got: %v - want: %v", got, want)
// 		}
// 	})

// }
