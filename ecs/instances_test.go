package ecs

import (
	"fmt"
	"testing"


	if err != nil {
		fmt.Printf("Failed to describe Instance %s vnc url: %v \n", TestInstanceId, err)
	} else {
		fmt.Printf("VNC URL = %s \n", instanceVncUrl)
	}
}

func ExampleClient_StopInstance() {
	fmt.Printf("Stop Instance Example\n")

	client := NewTestClient()

	err := client.StopInstance(TestInstanceId, true)

	if err != nil {
		fmt.Printf("Failed to stop Instance %s vnc url: %v \n", TestInstanceId, err)
	}
}

func ExampleClient_DeleteInstance() {
	fmt.Printf("Delete Instance Example")

	client := NewTestClient()

	err := client.DeleteInstance(TestInstanceId)

	if err != nil {
		fmt.Printf("Failed to delete Instance %s vnc url: %v \n", TestInstanceId, err)
	}
}

func TestECSInstance(t *testing.T) {
	if TestQuick {
		return
	}
	client := NewTestClient()
	instance, err := client.DescribeInstanceAttribute(TestInstanceId)
	if err != nil {
		t.Fatalf("Failed to describe instance %s: %v", TestInstanceId, err)
	}
	t.Logf("Instance: %++v  %v", instance, err)
	err = client.StopInstance(TestInstanceId, true)
	if err != nil {
		t.Errorf("Failed to stop instance %s: %v", TestInstanceId, err)
	}
	status, err := client.WaitForInstance(TestInstanceId, defaultInstanceStrategy)
	if err != nil || status != Stopped {
		t.Errorf("Instance %s is failed to stop: %v", TestInstanceId, err)
	}
	t.Logf("Instance %s is stopped successfully.", TestInstanceId)
	err = client.StartInstance(TestInstanceId)
	if err != nil {
		t.Errorf("Failed to start instance %s: %v", TestInstanceId, err)
	}
	status, err = client.WaitForInstance(TestInstanceId, defaultInstanceStrategy)
	if err != nil  || status != Running {
		t.Errorf("Instance %s is failed to start: %v", TestInstanceId, err)
	}
	t.Logf("Instance %s is running successfully.", TestInstanceId)
	err = client.RebootInstance(TestInstanceId, true)
	if err != nil {
		t.Errorf("Failed to restart instance %s: %v", TestInstanceId, err)
	}
	status, err = client.WaitForInstance(TestInstanceId, defaultInstanceStrategy)
	if err != nil || status != Running {
		t.Errorf("Instance %s is failed to restart: %v", TestInstanceId, err)
	}
	t.Logf("Instance %s is running successfully.", TestInstanceId)
}

func TestECSInstanceCreationAndDeletion(t *testing.T) {

	if TestIAmRich == false { // Avoid payment
		return
	}

	client := NewTestClient()
	instance, err := client.DescribeInstanceAttribute(TestInstanceId)
	t.Logf("Instance: %++v  %v", instance, err)

	args := CreateInstanceArgs{
		RegionId:        instance.RegionId,
		ImageId:         instance.ImageId,
		InstanceType:    "ecs.t1.small",
		SecurityGroupId: instance.SecurityGroupIds.SecurityGroupId[0],
	}

	instanceId, err := client.CreateInstance(&args)
	if err != nil {
		t.Errorf("Failed to create instance from Image %s: %v", args.ImageId, err)
	}
	t.Logf("Instance %s is created successfully.", instanceId)

	instance, err = client.DescribeInstanceAttribute(instanceId)
	t.Logf("Instance: %++v  %v", instance, err)

	strategy := util.AttemptStrategy{
		Min:   5,
		Total: 60 * time.Second,
		Delay: 5 * time.Second,
	}

	status, err := client.WaitForInstance(instanceId, strategy)

	if err != nil || status != Stopped {
		t.Errorf("Instance %s is failed to create: %v", instanceId, err)
	}

	err = client.StartInstance(instanceId)
	if err != nil {
		t.Errorf("Failed to start instance %s: %v", instanceId, err)
	}

	status, err = client.WaitForInstance(instanceId, defaultInstanceStrategy)

	if err != nil || status != Running {
		t.Errorf("Instance %s is failed to running: %v", instanceId, err)
	}

	err = client.StopInstance(instanceId, true)
	if err != nil {
		t.Errorf("Failed to stop instance %s: %v", instanceId, err)
	}
	status, err = client.WaitForInstance(instanceId, defaultInstanceStrategy)

	if err != nil  || status != Stopped {
		t.Errorf("Instance %s is failed to stop: %v", instanceId, err)
	}
	t.Logf("Instance %s is stopped successfully.", instanceId)

	err = client.DeleteInstance(instanceId)

	if err != nil {
		t.Errorf("Failed to delete instance %s: %v", instanceId, err)
	}
	t.Logf("Instance %s is deleted successfully.", instanceId)
}
