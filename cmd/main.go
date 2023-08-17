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

package main

import (
	"flag"

	"k8s.io/klog"

	"github.com/csi-addons/volume-condition-checker/internal/volume"
)

func main() {
	var nodename string

	flag.StringVar(&nodename, "nodename", "", "name of the kubernetes node")
	klog.InitFlags(nil)
	flag.Parse()

	lister := volume.NewVolumeAttachmentLister()
	if nodename != "" {
		lister.SetNodename(nodename)
	}
	vols, err := lister.ListVolumes()
	if err != nil {
		klog.Fatalf("failed to get VolumeAttachments: %v", err)
	}

	for _, v := range vols {
		drivername := v.GetDriver()
		klog.Infof("volume has driver %q: %v\n", drivername, v)

		drv, err := volume.FindDriver(drivername)
		if err != nil {
			klog.Errorf("could not find driver %q: %v", drivername, err)
			continue
		}

		if !drv.SupportsVolumeCondition() {
			continue
		}

		healthy, err := drv.IsHealthy(v)
		if err != nil {
			klog.Errorf("failed to check if %q is healthy: %v", v.GetVolumeID(), err)
			continue
		}
		klog.Infof("volume-handle %q is healthy: %v\n", v.GetVolumeID(), healthy)
	}
}
