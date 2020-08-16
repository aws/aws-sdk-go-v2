package external

import (
	"os"
	"path"
	"testing"
)

// create config file if not exist

func CreateIfNotExist(filepath string) error {

	dir, _ := path.Split(filepath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func TestAddProfileCredentials(t *testing.T) {
	t.Run("TestAddProfileCredentials1", func(t *testing.T) {
		path := DefaultSharedCredentialsFilename()
		CreateIfNotExist(path)
		AccessIdTest := "AccessIdTest"
		SecretKeyTest := "SecretKeyTest"
		ProfileTest := "Addfirsttest"
		c := &Credentials{Profile: &ProfileTest, AccessKeyId: &AccessIdTest, SecretAccessKey: &SecretKeyTest}
		got, _ := AddProfileCredentials(c)
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}
func TestDeleteProfileCredentials(t *testing.T) {
	t.Run("TestDeleteprofile", func(t *testing.T) {
		got, _ := DeleteProfileCredentials("Addfirsttest")
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}

func TestAddProfileConfig(t *testing.T) {
	t.Run("TestAddProfileConfig", func(t *testing.T) {
		path := DefaultSharedConfigFilename()
		CreateIfNotExist(path)
		ProfiletTest := "firstconfig"
		RegionTest := "us-east-1"
		OutputTest := "json"
		c := &CredentialsConfig{Profile: &ProfiletTest, Region: &RegionTest, Output: &OutputTest}
		got, _ := AddProfileConfig(c)
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}

func TestDeleteProfileConfig(t *testing.T) {
	t.Run("TestDeleteprofile", func(t *testing.T) {
		got, _ := DeleteProfileConfig("firstconfig")
		want := true
		if got != want {
			t.Errorf("Got: %v - want: %v", got, want)
		}
	})

}
