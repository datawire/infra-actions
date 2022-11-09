package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	mock_aws "github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_InstancesAreReturnedWhenThereAreNoErrors(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEc2Client := mock_aws.NewMockAwsEc2ClientInterface(mockCtrl)

	ec2Client := &Ec2Client{
		Client: mockEc2Client,
	}

	t.Run("Should return no details when there are no EC2 instances running", func(t *testing.T) {
		expectedDetails := []*InstanceDetails{}

		mockEc2Client.
			EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(interface{}, interface{}, ...interface{}) (*ec2.DescribeInstancesOutput, error) {
				return &ec2.DescribeInstancesOutput{}, nil
			}).
			AnyTimes()

		instanceDetails, err := ec2Client.GetInstances([]types.Filter{})
		require.NoError(t, err)
		assert.Equal(t, expectedDetails, instanceDetails)
	})
}
