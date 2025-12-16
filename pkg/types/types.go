package types

// SnapshotDiskInfo contains VM moref, snapshot moref, disk paths, and compute resource path for inspection
// This is used by both vm_service (to retrieve the info) and inspection (to use it)
// Supports multiple disks - DiskPaths and BaseDiskPaths are arrays
type SnapshotDiskInfo struct {
	VMMoref             string
	SnapshotMoref       string
	DiskPaths           []string // Current disk paths (may include snapshot deltas)
	BaseDiskPaths       []string // Base disk paths (without snapshot deltas)
	ComputeResourcePath string   // Path to compute resource (host/cluster) for vpx:// URL (e.g., "/Datacenter/Cluster/host.example.com")
}
