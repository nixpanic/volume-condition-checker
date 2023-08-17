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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

type VolumeLister interface {
	SetNodename(string)
	ListVolumes() ([]Volume, error)
}

type volumeAttachmentLister struct {
	clientset *kubernetes.Clientset

	// nodename is used to filter the attachments by node
	nodename string
}

func NewVolumeAttachmentLister() VolumeLister {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("could not get Kubernetes cluster configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("could not create a Kubernetes client: %v", err)
	}

	return &volumeAttachmentLister{
		clientset: clientset,
	}
}

func (val *volumeAttachmentLister) SetNodename(nodename string) {
	val.nodename = nodename
}

func (val *volumeAttachmentLister) ListVolumes() ([]Volume, error) {
	ctx := context.TODO()

	vas, err := val.clientset.StorageV1().VolumeAttachments().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Fatalf("could not get VolumeAttachments: %v", err)
	}

	vols := []Volume{}

	for _, va := range vas.Items {
		if val.nodename != "" && va.Spec.NodeName != val.nodename {
			continue
		}

		vols = append(vols, newCSIVolume(ctx, val.clientset, &va))
	}

	return vols, nil
}
