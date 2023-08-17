/*
Copyright 2023 The Kubernetes-CSI-Addons Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

type Driver interface {
	SupportsVolumeCondition() bool
	IsHealthy(Volume) (*bool, error)
}

type csiDriver struct {
	name string

	identityClient csi.IdentityClient
	nodeClient     csi.NodeClient

	node string

	supportVolumeCondition *bool
}

func FindDriver(name string) (Driver, error) {
	endpoint := GetPlatform().GetCSISocket(name)
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to endpoint %s: %w", endpoint, err)
	}

	drv := &csiDriver{
		name:           name,
		identityClient: csi.NewIdentityClient(conn),
		nodeClient:     csi.NewNodeClient(conn),
	}

	// TODO: verify if drv has the "name" identity

	return drv, nil
}

// TODO: verify that the driver provides NodeServiceCapability_RPC_VOLUME_CONDITION
func (drv *csiDriver) SupportsVolumeCondition() bool {
	if drv.supportVolumeCondition != nil {
		return *drv.supportVolumeCondition
	}

	yes := true
	no := false

	res, err := drv.nodeClient.NodeGetCapabilities(context.TODO(), &csi.NodeGetCapabilitiesRequest{})
	if err != nil {
		klog.Errorf("failed to get capabilities of driver %q: %v", drv.name, err)
		return false
	}

	for _, capability := range res.GetCapabilities() {
		if capability.GetRpc().GetType() == csi.NodeServiceCapability_RPC_VOLUME_CONDITION {
			drv.supportVolumeCondition = &yes
			return true
		}
	}

	klog.Infof("driver %q does not support VOLUME_CONDITION", drv.name)
	drv.supportVolumeCondition = &no
	return false
}

func (drv *csiDriver) IsHealthy(v Volume) (*bool, error) {
	if !drv.SupportsVolumeCondition() {
		return nil, fmt.Errorf("driver %q does not support VOLUME_CONDITION", drv.name)
	}

	volumePath := GetPlatform().GetStagingPath(drv.name, v.GetVolumeID())
	if volumePath == "" {
		// sttaging path not found, use publish path
		volumePath = GetPlatform().GetPublishPath(drv.name, v.GetVolumeID())
	}

	req := &csi.NodeGetVolumeStatsRequest{
		VolumeId:   v.GetVolumeID(),
		VolumePath: volumePath,
	}

	res, err := drv.nodeClient.NodeGetVolumeStats(context.TODO(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call NodeGetVolumeStats: %w", err)
	}

	if res.GetVolumeCondition() == nil {
		return nil, fmt.Errorf("VolumeCondition unknown")
	}

	healthy := !res.GetVolumeCondition().Abnormal

	return &healthy, nil
}
