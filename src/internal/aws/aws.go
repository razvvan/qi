package aws

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"gopkg.in/ini.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	keyAwsAccessKeyID     = "aws_access_key_id"
	keyAwsSecretAccessKey = "aws_secret_access_key"
	keyAwsMfaDevice       = "aws_mfa_device"
	keyAwsSessionToken    = "aws_session_token"
)

func GenerateNewSessionCredentials(envName, profilePrefix string) error {
	sectionName := profilePrefix + envName
	section, err := loadLongTermSection(sectionName)
	if err != nil {
		return err
	}

	awsMFADevice, err := section.GetKey(keyAwsMfaDevice)
	if err != nil {
		return err
	}

	err = updateMFAProfile(envName, awsMFADevice.String(), profilePrefix)
	if err != nil {
		return err
	}

	return nil
}

func updateMFAProfile(envName, awsMFADevice, profilePrefix string) error {
	sectionName := profilePrefix + envName + "-mfa"

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profilePrefix + envName,
	})
	if err != nil {
		return err
	}

	mfaCode := ""
	err = survey.AskOne(&survey.Input{Message: "MFA Code"}, &mfaCode)
	if err != nil {
		return err
	}

	svcMFA := sts.New(sess)
	sessionTokenOutput, err := svcMFA.GetSessionToken(
		&sts.GetSessionTokenInput{TokenCode: &mfaCode, SerialNumber: &awsMFADevice},
	)
	if err != nil {
		return err
	}

	_, err = saveCredentials(
		sectionName,
		*sessionTokenOutput.Credentials.AccessKeyId,
		*sessionTokenOutput.Credentials.SecretAccessKey,
		*sessionTokenOutput.Credentials.SessionToken,
		"",
	)
	if err != nil {
		return err
	}

	return nil
}

func loadLongTermSection(sectionName string) (*ini.Section, error) {
	_, section, err := loadSectionFromFile(sectionName)

	if err != nil && strings.Contains(err.Error(), "does not exist") {
		section, err = populateSection(sectionName)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return section, nil
}

func populateSection(sectionName string) (*ini.Section, error) {
	accessKeyID, secretAccessKey, mfaDevice, err := askForCredentials()
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")},
	})
	if err != nil {
		return nil, err
	}

	svc := sts.New(sess)
	_, err = svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	section, err := saveCredentials(sectionName, accessKeyID, secretAccessKey, "", mfaDevice)
	if err != nil {
		return nil, err
	}

	return section, nil
}
