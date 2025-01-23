package autoscaling

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

// manually tests the most complex jmespath expression in the entire SDK
func TestGroupInServiceStateRetryable(t *testing.T) {
	for name, tt := range map[string]struct {
		Output *DescribeAutoScalingGroupsOutput
		Expect bool
	}{
		// terminal cases: there are no groups that have NOT spun up yet,
		// indicated by # of InService < MinSize
		"empty output": {
			Output: &DescribeAutoScalingGroupsOutput{},
			Expect: false,
		},
		"empty group list": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{},
			},
			Expect: false,
		},
		"empty instance list": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{},
					},
				},
			},
			Expect: false,
		},
		"1 !InService": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateQuarantined},
						},
					},
				},
			},
			Expect: false,
		},
		"1 InService, MinSize nil": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateInService},
						},
					},
				},
			},
			Expect: false,
		},
		"1 InService, MinSize 0": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(0),
					},
				},
			},
			Expect: false,
		},
		"1 InService, MinSize 1": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(1),
					},
				},
			},
			Expect: false,
		},
		// retry cases: at least one group is spinning up
		"0 InService(nil slice), MinSize 2": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						MinSize: aws.Int32(2),
					},
				},
			},
			Expect: true,
		},
		"0 InService(empty slice), MinSize 2": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{},
						MinSize:   aws.Int32(2),
					},
				},
			},
			Expect: true,
		},
		"1 Pending, MinSize 2": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStatePending},
						},
						MinSize: aws.Int32(2),
					},
				},
			},
			Expect: true,
		},
		"1 Pending, 1 InService, MinSize 2": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStatePending},
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(2),
					},
				},
			},
			Expect: true,
		},
		// one group is done but another is spinning up
		// it's unclear if the service would ever return something like this,
		// but that's how the waiter is written
		"(1/2) + (2/2)": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStatePending},
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(2),
					},
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateInService},
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(2),
					},
				},
			},
			Expect: true,
		},
		// and then a terminal case where it spun up
		"2 InService, MinSize 2": {
			Output: &DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []types.AutoScalingGroup{
					{
						Instances: []types.Instance{
							{LifecycleState: types.LifecycleStateInService},
							{LifecycleState: types.LifecycleStateInService},
						},
						MinSize: aws.Int32(2),
					},
				},
			},
			Expect: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			fmt.Println(name)
			retryable, err := groupInServiceStateRetryable(context.Background(), nil, tt.Output, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.Expect != retryable {
				t.Errorf("%s: expected retryable=%v, got %v", name, tt.Expect, retryable)
			}
		})
	}
}
