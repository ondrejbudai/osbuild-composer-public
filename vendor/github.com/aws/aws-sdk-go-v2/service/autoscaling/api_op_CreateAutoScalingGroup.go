// Code generated by smithy-go-codegen DO NOT EDIT.

package autoscaling

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

//	We strongly recommend using a launch template when calling this operation to
//
// ensure full functionality for Amazon EC2 Auto Scaling and Amazon EC2.
//
// Creates an Auto Scaling group with the specified name and attributes.
//
// If you exceed your maximum limit of Auto Scaling groups, the call fails. To
// query this limit, call the DescribeAccountLimitsAPI. For information about updating this limit, see [Quotas for Amazon EC2 Auto Scaling]
// in the Amazon EC2 Auto Scaling User Guide.
//
// If you're new to Amazon EC2 Auto Scaling, see the introductory tutorials in [Get started with Amazon EC2 Auto Scaling] in
// the Amazon EC2 Auto Scaling User Guide.
//
// Every Auto Scaling group has three size properties ( DesiredCapacity , MaxSize ,
// and MinSize ). Usually, you set these sizes based on a specific number of
// instances. However, if you configure a mixed instances policy that defines
// weights for the instance types, you must specify these sizes with the same units
// that you use for weighting instances.
//
// [Get started with Amazon EC2 Auto Scaling]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/get-started-with-ec2-auto-scaling.html
// [Quotas for Amazon EC2 Auto Scaling]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-quotas.html
func (c *Client) CreateAutoScalingGroup(ctx context.Context, params *CreateAutoScalingGroupInput, optFns ...func(*Options)) (*CreateAutoScalingGroupOutput, error) {
	if params == nil {
		params = &CreateAutoScalingGroupInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateAutoScalingGroup", params, optFns, c.addOperationCreateAutoScalingGroupMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateAutoScalingGroupOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateAutoScalingGroupInput struct {

	// The name of the Auto Scaling group. This name must be unique per Region per
	// account.
	//
	// The name can contain any ASCII character 33 to 126 including most punctuation
	// characters, digits, and upper and lowercased letters.
	//
	// You cannot use a colon (:) in the name.
	//
	// This member is required.
	AutoScalingGroupName *string

	// The maximum size of the group.
	//
	// With a mixed instances policy that uses instance weighting, Amazon EC2 Auto
	// Scaling may need to go above MaxSize to meet your capacity requirements. In
	// this event, Amazon EC2 Auto Scaling will never go above MaxSize by more than
	// your largest instance weight (weights that define how many units each instance
	// contributes to the desired capacity of the group).
	//
	// This member is required.
	MaxSize *int32

	// The minimum size of the group.
	//
	// This member is required.
	MinSize *int32

	// A list of Availability Zones where instances in the Auto Scaling group can be
	// created. Used for launching into the default VPC subnet in each Availability
	// Zone when not using the VPCZoneIdentifier property, or for attaching a network
	// interface when an existing network interface ID is specified in a launch
	// template.
	AvailabilityZones []string

	// Indicates whether Capacity Rebalancing is enabled. Otherwise, Capacity
	// Rebalancing is disabled. When you turn on Capacity Rebalancing, Amazon EC2 Auto
	// Scaling attempts to launch a Spot Instance whenever Amazon EC2 notifies that a
	// Spot Instance is at an elevated risk of interruption. After launching a new
	// instance, it then terminates an old instance. For more information, see [Use Capacity Rebalancing to handle Amazon EC2 Spot Interruptions]in the
	// in the Amazon EC2 Auto Scaling User Guide.
	//
	// [Use Capacity Rebalancing to handle Amazon EC2 Spot Interruptions]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-capacity-rebalancing.html
	CapacityRebalance *bool

	// Reserved.
	Context *string

	//  Only needed if you use simple scaling policies.
	//
	// The amount of time, in seconds, between one scaling activity ending and another
	// one starting due to simple scaling policies. For more information, see [Scaling cooldowns for Amazon EC2 Auto Scaling]in the
	// Amazon EC2 Auto Scaling User Guide.
	//
	// Default: 300 seconds
	//
	// [Scaling cooldowns for Amazon EC2 Auto Scaling]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-scaling-cooldowns.html
	DefaultCooldown *int32

	// The amount of time, in seconds, until a new instance is considered to have
	// finished initializing and resource consumption to become stable after it enters
	// the InService state.
	//
	// During an instance refresh, Amazon EC2 Auto Scaling waits for the warm-up
	// period after it replaces an instance before it moves on to replacing the next
	// instance. Amazon EC2 Auto Scaling also waits for the warm-up period before
	// aggregating the metrics for new instances with existing instances in the Amazon
	// CloudWatch metrics that are used for scaling, resulting in more reliable usage
	// data. For more information, see [Set the default instance warmup for an Auto Scaling group]in the Amazon EC2 Auto Scaling User Guide.
	//
	// To manage various warm-up settings at the group level, we recommend that you
	// set the default instance warmup, even if it is set to 0 seconds. To remove a
	// value that you previously set, include the property but specify -1 for the
	// value. However, we strongly recommend keeping the default instance warmup
	// enabled by specifying a value of 0 or other nominal value.
	//
	// Default: None
	//
	// [Set the default instance warmup for an Auto Scaling group]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-default-instance-warmup.html
	DefaultInstanceWarmup *int32

	// The desired capacity is the initial capacity of the Auto Scaling group at the
	// time of its creation and the capacity it attempts to maintain. It can scale
	// beyond this capacity if you configure auto scaling. This number must be greater
	// than or equal to the minimum size of the group and less than or equal to the
	// maximum size of the group. If you do not specify a desired capacity, the default
	// is the minimum size of the group.
	DesiredCapacity *int32

	// The unit of measurement for the value specified for desired capacity. Amazon
	// EC2 Auto Scaling supports DesiredCapacityType for attribute-based instance type
	// selection only. For more information, see [Create a mixed instances group using attribute-based instance type selection]in the Amazon EC2 Auto Scaling User
	// Guide.
	//
	// By default, Amazon EC2 Auto Scaling specifies units , which translates into
	// number of instances.
	//
	// Valid values: units | vcpu | memory-mib
	//
	// [Create a mixed instances group using attribute-based instance type selection]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/create-mixed-instances-group-attribute-based-instance-type-selection.html
	DesiredCapacityType *string

	// The amount of time, in seconds, that Amazon EC2 Auto Scaling waits before
	// checking the health status of an EC2 instance that has come into service and
	// marking it unhealthy due to a failed health check. This is useful if your
	// instances do not immediately pass their health checks after they enter the
	// InService state. For more information, see [Set the health check grace period for an Auto Scaling group] in the Amazon EC2 Auto Scaling User
	// Guide.
	//
	// Default: 0 seconds
	//
	// [Set the health check grace period for an Auto Scaling group]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/health-check-grace-period.html
	HealthCheckGracePeriod *int32

	// A comma-separated value string of one or more health check types.
	//
	// The valid values are EC2 , ELB , and VPC_LATTICE . EC2 is the default health
	// check and cannot be disabled. For more information, see [Health checks for instances in an Auto Scaling group]in the Amazon EC2 Auto
	// Scaling User Guide.
	//
	// Only specify EC2 if you must clear a value that was previously set.
	//
	// [Health checks for instances in an Auto Scaling group]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-health-checks.html
	HealthCheckType *string

	// The ID of the instance used to base the launch configuration on. If specified,
	// Amazon EC2 Auto Scaling uses the configuration values from the specified
	// instance to create a new launch configuration. To get the instance ID, use the
	// Amazon EC2 [DescribeInstances]API operation. For more information, see [Create an Auto Scaling group using parameters from an existing instance] in the Amazon EC2 Auto
	// Scaling User Guide.
	//
	// [Create an Auto Scaling group using parameters from an existing instance]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/create-asg-from-instance.html
	// [DescribeInstances]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
	InstanceId *string

	// An instance maintenance policy. For more information, see [Set instance maintenance policy] in the Amazon EC2
	// Auto Scaling User Guide.
	//
	// [Set instance maintenance policy]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-instance-maintenance-policy.html
	InstanceMaintenancePolicy *types.InstanceMaintenancePolicy

	// The name of the launch configuration to use to launch instances.
	//
	// Conditional: You must specify either a launch template ( LaunchTemplate or
	// MixedInstancesPolicy ) or a launch configuration ( LaunchConfigurationName or
	// InstanceId ).
	LaunchConfigurationName *string

	// Information used to specify the launch template and version to use to launch
	// instances.
	//
	// Conditional: You must specify either a launch template ( LaunchTemplate or
	// MixedInstancesPolicy ) or a launch configuration ( LaunchConfigurationName or
	// InstanceId ).
	//
	// The launch template that is specified must be configured for use with an Auto
	// Scaling group. For more information, see [Create a launch template for an Auto Scaling group]in the Amazon EC2 Auto Scaling User
	// Guide.
	//
	// [Create a launch template for an Auto Scaling group]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/create-launch-template.html
	LaunchTemplate *types.LaunchTemplateSpecification

	// One or more lifecycle hooks to add to the Auto Scaling group before instances
	// are launched.
	LifecycleHookSpecificationList []types.LifecycleHookSpecification

	// A list of Classic Load Balancers associated with this Auto Scaling group. For
	// Application Load Balancers, Network Load Balancers, and Gateway Load Balancers,
	// specify the TargetGroupARNs property instead.
	LoadBalancerNames []string

	// The maximum amount of time, in seconds, that an instance can be in service. The
	// default is null. If specified, the value must be either 0 or a number equal to
	// or greater than 86,400 seconds (1 day). For more information, see [Replace Auto Scaling instances based on maximum instance lifetime]in the Amazon
	// EC2 Auto Scaling User Guide.
	//
	// [Replace Auto Scaling instances based on maximum instance lifetime]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-max-instance-lifetime.html
	MaxInstanceLifetime *int32

	// The mixed instances policy. For more information, see [Auto Scaling groups with multiple instance types and purchase options] in the Amazon EC2 Auto
	// Scaling User Guide.
	//
	// [Auto Scaling groups with multiple instance types and purchase options]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-mixed-instances-groups.html
	MixedInstancesPolicy *types.MixedInstancesPolicy

	// Indicates whether newly launched instances are protected from termination by
	// Amazon EC2 Auto Scaling when scaling in. For more information about preventing
	// instances from terminating on scale in, see [Use instance scale-in protection]in the Amazon EC2 Auto Scaling User
	// Guide.
	//
	// [Use instance scale-in protection]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-instance-protection.html
	NewInstancesProtectedFromScaleIn *bool

	// The name of the placement group into which to launch your instances. For more
	// information, see [Placement groups]in the Amazon EC2 User Guide for Linux Instances.
	//
	// A cluster placement group is a logical grouping of instances within a single
	// Availability Zone. You cannot specify multiple Availability Zones and a cluster
	// placement group.
	//
	// [Placement groups]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html
	PlacementGroup *string

	// The Amazon Resource Name (ARN) of the service-linked role that the Auto Scaling
	// group uses to call other Amazon Web Services service on your behalf. By default,
	// Amazon EC2 Auto Scaling uses a service-linked role named
	// AWSServiceRoleForAutoScaling , which it creates if it does not exist. For more
	// information, see [Service-linked roles]in the Amazon EC2 Auto Scaling User Guide.
	//
	// [Service-linked roles]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/autoscaling-service-linked-role.html
	ServiceLinkedRoleARN *string

	// One or more tags. You can tag your Auto Scaling group and propagate the tags to
	// the Amazon EC2 instances it launches. Tags are not propagated to Amazon EBS
	// volumes. To add tags to Amazon EBS volumes, specify the tags in a launch
	// template but use caution. If the launch template specifies an instance tag with
	// a key that is also specified for the Auto Scaling group, Amazon EC2 Auto Scaling
	// overrides the value of that instance tag with the value specified by the Auto
	// Scaling group. For more information, see [Tag Auto Scaling groups and instances]in the Amazon EC2 Auto Scaling User
	// Guide.
	//
	// [Tag Auto Scaling groups and instances]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-tagging.html
	Tags []types.Tag

	// The Amazon Resource Names (ARN) of the Elastic Load Balancing target groups to
	// associate with the Auto Scaling group. Instances are registered as targets with
	// the target groups. The target groups receive incoming traffic and route requests
	// to one or more registered targets. For more information, see [Use Elastic Load Balancing to distribute traffic across the instances in your Auto Scaling group]in the Amazon EC2
	// Auto Scaling User Guide.
	//
	// [Use Elastic Load Balancing to distribute traffic across the instances in your Auto Scaling group]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/autoscaling-load-balancer.html
	TargetGroupARNs []string

	// A policy or a list of policies that are used to select the instance to
	// terminate. These policies are executed in the order that you list them. For more
	// information, see [Configure termination policies for Amazon EC2 Auto Scaling]in the Amazon EC2 Auto Scaling User Guide.
	//
	// Valid values: Default | AllocationStrategy | ClosestToNextInstanceHour |
	// NewestInstance | OldestInstance | OldestLaunchConfiguration |
	// OldestLaunchTemplate |
	// arn:aws:lambda:region:account-id:function:my-function:my-alias
	//
	// [Configure termination policies for Amazon EC2 Auto Scaling]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-termination-policies.html
	TerminationPolicies []string

	// The list of traffic sources to attach to this Auto Scaling group. You can use
	// any of the following as traffic sources for an Auto Scaling group: Classic Load
	// Balancer, Application Load Balancer, Gateway Load Balancer, Network Load
	// Balancer, and VPC Lattice.
	TrafficSources []types.TrafficSourceIdentifier

	// A comma-separated list of subnet IDs for a virtual private cloud (VPC) where
	// instances in the Auto Scaling group can be created. If you specify
	// VPCZoneIdentifier with AvailabilityZones , the subnets that you specify must
	// reside in those Availability Zones.
	VPCZoneIdentifier *string

	noSmithyDocumentSerde
}

type CreateAutoScalingGroupOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateAutoScalingGroupMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpCreateAutoScalingGroup{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpCreateAutoScalingGroup{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateAutoScalingGroup"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpCreateAutoScalingGroupValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateAutoScalingGroup(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opCreateAutoScalingGroup(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateAutoScalingGroup",
	}
}