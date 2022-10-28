package main

import (
	"context"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

var ec2Client *ec2.Client

const AMI_MACOS_12_6_ARM64 = "ami-01b8fcd5770ceb9c1"

func init() {
	var err error

	ec2Client, err = newAwsClient()
	if err != nil {
		log.Fatalf("Error initializinf AWS client: %v", err)
	}
}

func newAwsClient() (*ec2.Client, error) {
	cfg, err := config2.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	return client, nil
}
