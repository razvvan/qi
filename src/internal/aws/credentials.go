package aws

import (
	"os/user"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/ini.v1"
)

func awsCredentialsFile() string {
	user, _ := user.Current()

	return user.HomeDir + "/.aws/credentials"
}

func askForCredentials() (accessKeyID string, secretAccessKey string, mfaDevice string, err error) {
	err = survey.AskOne(&survey.Password{Message: keyAwsAccessKeyID}, &accessKeyID)
	if err != nil {
		return "", "", "", err
	}

	err = survey.AskOne(&survey.Password{Message: keyAwsSecretAccessKey}, &secretAccessKey)
	if err != nil {
		return "", "", "", err
	}

	err = survey.AskOne(&survey.Input{Message: keyAwsMfaDevice}, &mfaDevice)
	if err != nil {
		return "", "", "", err
	}

	return accessKeyID, secretAccessKey, mfaDevice, nil
}

func loadSectionFromFile(sectionName string) (*ini.File, *ini.Section, error) {
	cfg, err := ini.Load(awsCredentialsFile())
	if err != nil {
		return nil, nil, err
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return cfg, nil, err
	}

	return cfg, section, nil
}

func saveCredentials(
	sectionName string,
	accessKeyID string,
	secretAccessKey string,
	sessionToken string,
	mfaDevice string,
) (*ini.Section, error) {
	cfg, section, err := loadSectionFromFile(sectionName)
	if err != nil && strings.Contains(err.Error(), "does not exist") {
		section, err = cfg.NewSection(sectionName)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	section.Key(keyAwsAccessKeyID).SetValue(accessKeyID)
	section.Key(keyAwsSecretAccessKey).SetValue(secretAccessKey)

	if mfaDevice != "" {
		section.Key(keyAwsMfaDevice).SetValue(mfaDevice)
	}

	if sessionToken != "" {
		section.Key(keyAwsSessionToken).SetValue(sessionToken)
	}

	err = cfg.SaveTo(awsCredentialsFile())
	if err != nil {
		return nil, err
	}

	return section, nil
}
