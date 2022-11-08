package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

type Ec2Client struct {
	Client AwsEc2ClientInterface
}

func NewEc2Client() *Ec2Client {
	var err error

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error initializing the AWS client: %v", err)
	}

	return &Ec2Client{ec2.NewFromConfig(cfg)}
}
