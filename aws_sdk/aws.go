package aws_sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/polo44/goutilities"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
)

type AwsDefaultConfig struct {
	UserName        string `json:"UserName" yaml:"UserName"`
	AccessKeyID     string `json:"AccessKeyID" yaml:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey" yaml:"SecretAccessKey"`
	AwsEndPoint     string `json:"AwsEndPoint" yaml:"AwsEndPoint"`
}

type AwsS3 struct {
	AwsDefaultConfig `yaml:",inline" json:",inline"`
	AwsBucket        string `json:"AwsBucket" yaml:"AwsBucket"`
	PrivateEndpoint  string `json:"PrivateEndpoint" yaml:"PrivateEndpoint"`
	PublicEndpoint   string `json:"PublicEndpoint" yaml:"PublicEndpoint"`
	AWSEncryptionKey string `json:"AWSEncryptionKey" yaml:"AWSEncryptionKey"`
}

// AWSConfig holds info about AWS Config
type AWSConfig struct {
	S3  AwsS3            `json:"S3" yaml:"S3"`
	SES AwsDefaultConfig `json:"SES" yaml:"SES"`
}

// S3Ctx holds the default context for s3 operations
var AwsCfg *AWSConfig

// Init inits the repo
func Init(fileName string) error {

	// ─── WE LOAD AWS ────────────────────────────────────────────────────────────────
	AwsCfg = &AWSConfig{}

	reader, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("can't open aws config file. err: %s", err.Error())
	}

	err = yaml.NewDecoder(reader).Decode(&AwsCfg)
	if err != nil {
		return fmt.Errorf("can't load aws project settings. err: %s", err.Error())
	}

	if strings.TrimSpace(AwsCfg.S3.AccessKeyID) == "" {
		return errors.New(goutilities.ErrorsRender("Can't init aws", errors.New("AWS AccessKeyID can't be blank")))
	}
	if strings.TrimSpace(AwsCfg.S3.SecretAccessKey) == "" {
		return errors.New(goutilities.ErrorsRender("Can't init aws", errors.New("AWS SecretAccessKey can't be blank")))
	}
	if strings.TrimSpace(AwsCfg.S3.AwsBucket) == "" {
		return errors.New(goutilities.ErrorsRender("Can't init aws", errors.New("AwsBucket can't be blank")))
	}

	// ─── WE SETUP AWS S3 ────────────────────────────────────────────────────────────────
	awsS3Cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: AwsCfg.S3.AccessKeyID, SecretAccessKey: AwsCfg.S3.SecretAccessKey,
			},
		},
	), config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: AwsCfg.S3.AwsEndPoint}, nil
		})))

	if err != nil {
		return errors.New(goutilities.ErrorsRender("Can't init aws s3", err))
	}
	S3Sess = &S3Session{client: s3.NewFromConfig(awsS3Cfg)}

	return nil
}

// WriterAtTOResponseWriter allows us to wrap the WriteAt interface require in aws.ObjectDownload an write directly inside an ResponseWriter
type WriterAtTOResponseWriter struct {
	w *http.ResponseWriter
}

// WriteAt Override the WriteAt to write directly inside the ResponseWriter
func (writer WriterAtTOResponseWriter) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return http.ResponseWriter(*writer.w).Write(p)
}
