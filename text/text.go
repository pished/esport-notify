package text

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func SendText(message string) {
	log.Println("Creating sesison...")
	sess := session.Must(session.NewSession())
	log.Println("Session created...")

	svc := sns.New(sess)
	log.Println("Service created...")

	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(os.Getenv("TARGET_PHONE")),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(resp)
}
