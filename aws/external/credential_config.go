// Package external AWS local Credentials configuration in Go SDK, Can also able to add or remove with profilename
package external

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/ini"
)

// Credentials composite  the  profile,accesskey,secretkey,,sessiontoken fields
type Credentials struct {

	// profile name which we need, if not mentioned it will create in default profile
	Profile *string `type:"string"`

	// The identifier used for the temporary security credentials. For more information,
	// see Using Temporary Security Credentials to Request Access to AWS Resources
	// (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html)
	// in the AWS IAM User Guide.
	AccessKeyID *string `type:"string" required:"true"`

	// The key that is used to sign the request. For more information, see Using
	// Temporary Security Credentials to Request Access to AWS Resources (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html)
	// in the AWS IAM User Guide.
	SecretAccessKey *string `type:"string" required:"true" sensitive:"true"`

	// The token used for temporary credentials. For more information, see Using
	// Temporary Security Credentials to Request Access to AWS Resources (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html)
	// in the AWS IAM User Guide.
	SessionToken *string `type:"string" sensitive:"true"`
	// contains filtered or unexported fields
}

// CredentialsConfig composite  the  profile,region,output fields
type CredentialsConfig struct {
	Profile *string `type:"string"`
	Region  *string `type:"string" required:"true"`
	Output  *string `type:"string"`
}

// DeleteProfileCredentials method to delete a profile in .aws/credentials
// Need to pass the profile name as argument
// if success it will return true
func DeleteProfileCredentials(p string) (bool, error) {

	creds := []string{}
	var out bool

	path := DefaultSharedCredentialsFilename()

	// file backup
	BackupConfig(path)

	sec, err := ini.OpenFile(path)
	if err != nil {
		fmt.Println("Not able to Read the credential file")

	}
	Getvalues, isAvailable := sec.GetSection(p)

	if isAvailable == false {
		fmt.Println("profile is not available: ", p)

	} else {
		profile := fmt.Sprintf("[%v]", p)
		creds = append(creds, profile)
		if Getvalues.Has("aws_access_key_id") {
			accessID := fmt.Sprintf("aws_access_key_id = %v", Getvalues.String("aws_access_key_id"))
			creds = append(creds, accessID)

		}
		if Getvalues.Has("aws_secret_access_key") {
			secretKey := fmt.Sprintf("aws_secret_access_key = %v", Getvalues.String("aws_secret_access_key"))
			creds = append(creds, secretKey)

		}
		if Getvalues.Has("aws_session_token") {
			sessionToken := fmt.Sprintf("aws_session_token = %v", Getvalues.String("aws_session_token"))
			creds = append(creds, sessionToken)

		}

		for _, v := range creds {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			newContents := strings.ReplaceAll(string(read), v, "")
			err = ioutil.WriteFile(path, []byte(newContents), 0644)
			if err != nil {
				panic(err)
			}

		}
		out = true

	}

	return out, err
}

// AddProfileCredentials method to Add a profile in .aws/credentials , if success it will return true
// Need to pass the *Credentials  type  as argument , it composite  the  profile,accesskey,secretkey,,sessiontoken fields
// if success it will return true
func AddProfileCredentials(c *Credentials) (bool, error) {

	var out bool
	var profile string
	var accessID string
	var secretKey string
	var sessionToken string

	path := DefaultSharedCredentialsFilename()

	// file backup
	BackupConfig(path)

	if c.Profile == nil {
		defaultsprofile := "default"
		c.Profile = &defaultsprofile
	}

	sec, err := ini.OpenFile(path)
	if err != nil {
		fmt.Println("Not able to Read the credential file")

	}
	_, isAvailable := sec.GetSection(*c.Profile)

	if isAvailable == true {
		fmt.Println("profile is already confiured: ", *c.Profile)

	} else {

		profile = fmt.Sprintf("[%v]", *c.Profile)

		if c.AccessKeyID != nil {
			accessID = fmt.Sprintf("aws_access_key_id = %v", *c.AccessKeyID)

		}
		if c.SecretAccessKey != nil {
			secretKey = fmt.Sprintf("aws_secret_access_key = %v", *c.SecretAccessKey)

		}
		if c.SessionToken == nil {

			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			} else {
				newContents := fmt.Sprintf("%v\n%v\n%v", profile, accessID, secretKey)

				Contents := fmt.Sprintf("%v\n%v", string(read), newContents)
				err = ioutil.WriteFile(path, []byte(Contents), 0644)
				if err != nil {
					panic(err)
				}
				out = true
			}

		} else {
			sessionToken = fmt.Sprintf("aws_session_token = %v", *c.SessionToken)
			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			} else {
				newContents := fmt.Sprintf("%v\n%v\n%v\n%v", profile, accessID, secretKey, sessionToken)

				Contents := fmt.Sprintf("%v\n%v", string(read), newContents)
				err = ioutil.WriteFile(path, []byte(Contents), 0644)
				if err != nil {
					panic(err)
				}
				out = true
			}

		}

	}

	return out, err
}

// AddProfileConfig method to Add a profile in .aws/config , if success it will return true
// Need to pass the *CredentialsConfig  type  as argument , it composite  the  profile,region,output fields
// if success it will return true
func AddProfileConfig(c *CredentialsConfig) (bool, error) {

	var out bool
	var profile string
	var region string
	var output string

	path := DefaultSharedConfigFilename()

	// file backup
	BackupConfig(path)

	if c.Profile == nil {
		defaultsprofile := "default"
		c.Profile = &defaultsprofile
	} else {
		profile = fmt.Sprintf("profile %v", *c.Profile)
	}

	sec, err := ini.OpenFile(path)
	if err != nil {
		fmt.Println("Not able to Read the config file")

	}
	_, isAvailable := sec.GetSection(profile)
	if isAvailable == true {
		fmt.Println("profile is already confiured: ", *c.Profile)

	} else {
		profile = fmt.Sprintf("[profile %v]", *c.Profile)

		if c.Region != nil {
			region = fmt.Sprintf("region = %v", *c.Region)
		}

		if c.Output == nil {

			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			} else {
				newContents := fmt.Sprintf("%v\n%v", profile, region)

				Contents := fmt.Sprintf("%v\n%v", string(read), newContents)
				err = ioutil.WriteFile(path, []byte(Contents), 0644)
				if err != nil {
					panic(err)
				}
				out = true
			}

		} else {
			output = fmt.Sprintf("output = %v", *c.Output)
			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			} else {
				newContents := fmt.Sprintf("%v\n%v\n%v", profile, region, output)

				Contents := fmt.Sprintf("%v\n%v", string(read), newContents)
				err = ioutil.WriteFile(path, []byte(Contents), 0644)
				if err != nil {
					panic(err)
				}
				out = true
			}

		}

	}

	return out, err
}

// DeleteProfileConfig method to delete a profile in .aws/config
// Need to pass the profile name as argument
// if success it will return true
func DeleteProfileConfig(p string) (bool, error) {

	var profile string
	var region string
	var output string

	var out bool

	path := DefaultSharedConfigFilename()

	sec, err := ini.OpenFile(path)
	if err != nil {
		fmt.Println("Not able to Read the config file")

	}

	if p != "default" {
		p = fmt.Sprintf("profile %v", p)
	}

	_, isAvailable := sec.GetSection(p)

	if isAvailable == false {
		fmt.Println("profile is not available: ", p)

	} else {
		ProfileList := sec.List()
		e := os.Rename(path, path+".bkp")
		if e != nil {
			log.Fatal(e)
		}

		for _, values := range ProfileList {
			if values != p {
				f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
				}
				defer f.Close()

				Getval, _ := sec.GetSection(values)
				profile = fmt.Sprintf("[%v]", values)
				if _, err := f.WriteString(profile + "\n"); err != nil {
					log.Println(err)
				}

				if Getval.Has("region") {
					region = fmt.Sprintf("region = %v", Getval.String("region"))
					if _, err := f.WriteString(region + "\n"); err != nil {
						log.Println(err)
					}

				}
				if Getval.Has("output") {
					output = fmt.Sprintf("output = %v", Getval.String("output"))
					if _, err := f.WriteString(output + "\n\n"); err != nil {
						log.Println(err)
					}

				}
				out = true

			}
			f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			out = true

		}

	}

	return out, err
}

//BackupConfig it will take back when performing operation config.bkp credentials.bkp
func BackupConfig(path string) bool {
	var out bool
	readbkp, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("file not exist")
	}
	err = ioutil.WriteFile(path+".bkp", []byte(string(readbkp)), 0644)
	if err != nil {
		fmt.Println("Error when creating backupfile")
	} else {
		out = true
	}

	return out
}
