package aws_s3

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToAwsS3(accessKey, secretKey, bucket, region, endpoint string, publisherID uint64, gameID uint64, email string, d []byte, size int64) (string, error) {

	originKey := fmt.Sprintf("%s-%s-%s", publisherID, gameID, email)
	data := []byte(originKey)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	})
	if err != nil {
		return "", err
	}

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(fmt.Sprintf("/profiles/%s", md5str)),
		ACL:                aws.String("public-read"),
		Body:               bytes.NewReader(d),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(d)),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf("https://%s.%s/profiles/%s", bucket, endpoint, md5str)
	return link, nil
}
