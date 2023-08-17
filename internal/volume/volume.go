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

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

type Volume interface {
	GetDriver() string
	GetNode() string
	GetVolumeID() string
}

type csiVolume struct {
	clientset            *kubernetes.Clientset
	drivername           string
	node                 string
	persistentVolumeName string
	volumeID             string
	volumePath           string
}

func newCSIVolume(ctx context.Context, clientset *kubernetes.Clientset, va *storagev1.VolumeAttachment) Volume {
	return &csiVolume{
		clientset:            clientset,
		drivername:           va.Spec.Attacher,
		node:                 va.Spec.NodeName,
		persistentVolumeName: *va.Spec.Source.PersistentVolumeName,

		// TODO: volumeID could be an inline-volume too
		volumeID: pvToVolumeID(ctx, clientset, *va.Spec.Source.PersistentVolumeName),
	}
}

func pvToVolumeID(ctx context.Context, clientset *kubernetes.Clientset, pvName string) string {
	pv, err := clientset.CoreV1().PersistentVolumes().Get(ctx, pvName, metav1.GetOptions{})
	if err != nil {
		klog.Fatalf("failed to get PersistentVolume %s: %v", pvName, err)
	}

	// TODO: check pointers to reach VolumeHandle
	return pv.Spec.CSI.VolumeHandle
}

func (vol *csiVolume) GetDriver() string {
	return vol.drivername
}

func (vol *csiVolume) GetNode() string {
	return vol.node
}

func (vol *csiVolume) GetVolumeID() string {
	return vol.volumeID
}
