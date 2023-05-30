This project is a POC that enables the end user to take a backup of a database instance in a cluster and later use that snapshot to recover the database instance when it restarts from scratch. `VolumeSnapshotContent` and `VolumeSnapshot` API resources are provided to create volume snapshots for users and administrators.

## Prerequisites
1. This project is in development mode run in a local kubernetes cluster. For this purpose, I have used `kind` application to generate local cluster.
2. This project also relies extensively on `VolumeSnapshot` kind to actually do the heavy lifting of taking backkups and restoration. 
3. Mysql database run as a POD and provisioned volume onto which data persists and intended to be backed up


## Theory
In order to take snapshots following resource needs to be created in cluster:
```
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: csi-do-test-snapshot
spec:
  source:
    persistentVolumeClaimName: csi-do-test-pvc
```

In order to restore the snapshots, following persistent volume claim need to be generated:
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: csi-do-test-pvc-restore
spec:
  dataSource:
    name: csi-do-test-snapshot
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

`VolumeSnapshotContent` is a snapshot taken from a volume in the cluster that has been provisioned by an administrator. It is a resource in the cluster just like a `PersistentVolume` is a cluster resource.

A `VolumeSnapshot` is a request for snapshot of a volume by a user. It is similar to a `PersistentVolumeClaim`.

> Note: Installation of the CRDs is the responsibility of the Kubernetes distribution.
> Without the required CRDs present, the creation of a VolumeSnapshotClass fails.

> Furthermore, the default storage provider in kind does not implement the CSI interface and thus is NOT capable of creating/handling volume snapshots. For that, you must first deploy a CSI driver. 

> You can use the following CSI snapshot CRDs and drivers from `https://github.com/kubernetes-csi/external-snapshotter/tree/release-3.0`

## Steps

- Deply the custom CRD to inform the cluster about the incoming CRs
- Create a cr to take user input, the controller will detect it and it will spin up a cr for snapshot
- To restore, we need create a pv where source has to be kept as snapshot, and then deployment has to modified and run again
- For restore we can do it in the same crd using boolean for now.
