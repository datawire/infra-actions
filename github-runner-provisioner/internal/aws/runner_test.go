package aws

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/aws_runners"
	mock_aws "github.com/datawire/infra-actions/github-runner-provisioner/internal/aws/mocks"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	mockCtrl      *gomock.Controller
	ec2Client     *Ec2Client
	mockEc2Client *mock_aws.MockAwsEc2ClientInterface
}

func setup(t *testing.T) *fixture {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	mockEc2Client := mock_aws.NewMockAwsEc2ClientInterface(mockCtrl)

	f := fixture{
		mockCtrl:      mockCtrl,
		ec2Client:     &Ec2Client{Client: mockEc2Client},
		mockEc2Client: mockEc2Client,
	}

	return &f
}

func Test_InstancesAreReturnedWhenThereAreNoErrors(t *testing.T) {
	t.Run("Should return no details when there are no EC2 instances running", func(t *testing.T) {
		f := setup(t)
		defer f.mockCtrl.Finish()

		expectedDetails := []*InstanceDetails{}
		noEc2Instances := &ec2.DescribeInstancesOutput{}

		f.mockEc2Client.
			EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(noEc2Instances, nil).
			AnyTimes()

		instanceDetails, err := f.ec2Client.GetInstances([]types.Filter{})
		require.NoError(t, err)
		assert.Equal(t, expectedDetails, instanceDetails)
	})

	t.Run("Should return the expected instances", func(t *testing.T) {
		f := setup(t)
		defer f.mockCtrl.Finish()

		expectedDetails := []*InstanceDetails{
			{
				LaunchTime:        utils.TimePtr(time.Unix(2000, 0)),
				InstanceId:        utils.StrPtr("ubuntu-arm64-instance"),
				ActionRunnerLabel: utils.StrPtr("ubuntu-arm64"),
			},
		}

		ec2Instances := &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							LaunchTime: utils.TimePtr(time.Unix(2000, 0)),
							InstanceId: utils.StrPtr("ubuntu-arm64-instance"),
							Tags: []types.Tag{
								{
									Key:   utils.StrPtr(aws_runners.LabelTag),
									Value: utils.StrPtr("ubuntu-arm64"),
								},
								{
									Key:   utils.StrPtr(aws_runners.NameTag),
									Value: utils.StrPtr(aws_runners.AppName),
								},
							},
						},
						{
							LaunchTime: utils.TimePtr(time.Unix(3000, 0)),
							InstanceId: utils.StrPtr("other-instance"),
							Tags:       []types.Tag{},
						},
					},
				},
			},
		}

		f.mockEc2Client.
			EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(ec2Instances, nil).
			AnyTimes()

		instanceDetails, err := f.ec2Client.GetInstances([]types.Filter{})
		require.NoError(t, err)
		assert.Equal(t, expectedDetails, instanceDetails)
	})

	t.Run("Should go over every page returned by the API", func(t *testing.T) {
		f := setup(t)
		defer f.mockCtrl.Finish()

		expectedDetails := []*InstanceDetails{
			{
				LaunchTime:        utils.TimePtr(time.Unix(2000, 0)),
				InstanceId:        utils.StrPtr("ubuntu-arm64-instance"),
				ActionRunnerLabel: utils.StrPtr("ubuntu-arm64"),
			},
		}

		page := &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							LaunchTime: utils.TimePtr(time.Unix(2000, 0)),
							InstanceId: utils.StrPtr("ubuntu-arm64-instance"),
							Tags: []types.Tag{
								{
									Key:   utils.StrPtr(aws_runners.LabelTag),
									Value: utils.StrPtr("ubuntu-arm64"),
								},
								{
									Key:   utils.StrPtr(aws_runners.NameTag),
									Value: utils.StrPtr(aws_runners.AppName),
								},
							},
						},
					},
				},
			},
		}

		f.mockEc2Client.
			EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ interface{}, params *ec2.DescribeInstancesInput, _ ...interface{}) (*ec2.DescribeInstancesOutput, error) {
				return page, nil
			}).
			AnyTimes()

		instanceDetails, err := f.ec2Client.GetInstances([]types.Filter{})
		require.NoError(t, err)
		assert.Equal(t, expectedDetails, instanceDetails)
	})
}

func Test_GetInstancesHandlesErrorGracefully(t *testing.T) {
	f := setup(t)
	defer f.mockCtrl.Finish()

	f.mockEc2Client.
		EXPECT().
		DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("Error calling AWS API")).
		AnyTimes()

	_, err := f.ec2Client.GetInstances([]types.Filter{})
	require.Error(t, err)
}
