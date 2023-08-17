# Volume Condition Checker

The Volume Condition Checker for Kubernetes uses the [Container Storage
Interface Specification's `NodeGetVolumeStats` operation][nodegetvolumestats]
to detect if a PersistentVolume has an _abnormal_ condition. CSI drivers can
return the condition of a volume in the `NodeVolumeStatsResponse` message.

## Abnormal Volume Condition reporting

... _to be done_

## Potential Consumers of Abnormal Volume Condition check results

- [Rook](https://rook.io) is a Kubernetes Operator that is able to [Network
  Fence][rook_fencing] a workernode where a Ceph volume is unhealthy.

- [Node Problem Detector][k8s_npd] provides a generic interface for reporting
  problems on a node. A project like [medik8s](https://medik8s.io/) can remedy
  node problems once they are reported.

## Dependencies

The `NodeGetVolumeStats` operation in the current CSI Specification (v1.8.0)
defines the `VolumeCondition` as an _alpha_ feature. Very few CSI-drivers seem
to implement the volume condition at the moment. Drivers are required to
implement the feature, and expose `VOLUME_CONDITION` as a
`NodeServiceCapability`, otherwise the Volume Condition Checker will not be
able to check the condition of the volume.

[nodegetvolumestats]: https://github.com/container-storage-interface/spec/blob/master/spec.md#nodegetvolumestats
[rook_fencing]: https://rook.github.io/docs/rook/v1.12/Storage-Configuration/Block-Storage-RBD/block-storage/#handling-node-loss
[k8s_npd]: https://github.com/kubernetes/node-problem-detector/
