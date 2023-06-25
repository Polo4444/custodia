package aws_sdk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Session struct {
	client *s3.Client
}

var S3Sess *S3Session

//
// ───────────────────────────────────────────────────── END GLOBAL VARIABLES ─────

// FilePathJoin join path on replace path separator \ by /
func FilePathJoin(elem ...string) string {
	return strings.ReplaceAll(filepath.Join(elem...), "\\", "/")
}

// ObjectUploadPublicWithReader upload a public object to a specified bucket with a reader
func (s *S3Session) ObjectUploadPublicWithReader(objectName string, dataToUpload io.Reader) (string, error) {

	// We init an uploader
	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
		Body:   dataToUpload,
		ACL:    types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return "", err
	}

	return strings.Replace(result.Location, AwsCfg.S3.PrivateEndpoint, AwsCfg.S3.PublicEndpoint, 1), nil
}

// ObjectUploadPrivateWithReader upload a private object to a specified bucket with a reader
func (s *S3Session) ObjectUploadPrivateWithReader(objectName string, dataToUpload io.Reader) (string, error) {

	// We init an uploader
	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
		Body:   dataToUpload,
		ACL:    types.ObjectCannedACLPrivate,
	})

	if err != nil {
		return "", err
	}

	return strings.Replace(result.Location, AwsCfg.S3.PrivateEndpoint, AwsCfg.S3.PublicEndpoint, 1), nil
}

// ObjectDownloadIntoResponseWriter download a object directly inside a response writer
func (s *S3Session) ObjectDownloadIntoResponseWriter(objectName string, w *http.ResponseWriter) error {

	// We init the downloader
	objectDownloaded, err := s.ObjectDownload(objectName)
	if err != nil {
		return err
	}

	http.ResponseWriter(*w).Header().Set("Content-Length", fmt.Sprintf("%d", len(objectDownloaded)))
	_, err = io.Copy(http.ResponseWriter(*w), bytes.NewReader(objectDownloaded))
	if err != nil {
		return err
	}

	return nil
}

// ObjectDetails get details about an object
func (s *S3Session) ObjectDetails(objectName string) (*s3.HeadObjectOutput, error) {

	return s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
	})
}

// ObjectDownload download a object an return a slice of byte
func (s *S3Session) ObjectDownload(objectName string) ([]byte, error) {

	// Get object info
	headObject, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, err
	}

	buf := make([]byte, int(headObject.ContentLength))
	// wrap with aws.WriteAtBuffer
	w := manager.NewWriteAtBuffer(buf)

	// We init the downloader e
	downloader := manager.NewDownloader(s.client)
	_, err = downloader.Download(context.Background(), w, &s3.GetObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// IsObjectExist check if an object/file exits
func (s *S3Session) IsObjectExist(objectName string) bool {

	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
	})

	return err == nil
}

// DeleteObject delete an object
func (s *S3Session) DeleteObject(objectName string) error {

	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(AwsCfg.S3.AwsBucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Session) RenameObject(objectName string, newName string) error {

	_, err := s.client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(AwsCfg.S3.AwsBucket),
		CopySource: aws.String(FilePathJoin(AwsCfg.S3.AwsBucket, objectName)),
		Key:        aws.String(newName),
		ACL:        types.ObjectCannedACLBucketOwnerFullControl,
	})

	if err != nil {
		return err
	}

	return nil
}
