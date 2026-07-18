package controller

import (
	talosv1alpha1 "github.com/alperencelik/talos-operator/api/v1alpha1"
)

const (
	TalosPlatformKey = "PLATFORM"
	// TalosModeContainer is the mode for Talos running in a container
	TalosModeContainer = "container"
	// TalosModeMetal is the mode for Talos running on bare metal
	TalosModeMetal = "metal"

	// MachineType
	TalosMachineTypeControlPlane = "controlplane"
	TalosMachineTypeWorker       = "worker"

	// Reconcile Modes — aliases of the api/v1alpha1 constants, kept so the
	// controller package keeps its historical names.

	// ReconcileModeAnnotation is the annotation key for the reconcile mode
	ReconcileModeAnnotation = talosv1alpha1.ReconcileModeAnnotation
	// ReconcileMode is the mode of the reconciliation, it could be Reconcile, Disable or DryRun
	ReconcileModeNormal  = talosv1alpha1.ReconcileModeNormal
	ReconcileModeDisable = talosv1alpha1.ReconcileModeDisable
	// ReconcileModeDryRun runs the reconciliation without performing any mutating operations.
	// Kubernetes writes are validated via server-side dry-run; Talos API and file-system
	// operations are skipped and reported as "Would do X" events.
	ReconcileModeDryRun = talosv1alpha1.ReconcileModeDryRun

	// ReconcileModeImport is the mode for importing existing Talos resources
	ReconcileModeImport = talosv1alpha1.ReconcileModeImport

	// For tests
	DefaultNamespace = "default"

	// Deleting is the reason used in conditions when a resource is being deleted
	ConditionReasonDeleting = "Deleting"

	// EventReasonDryRun is the reason used in events emitted while reconciling in DryRun mode
	EventReasonDryRun = "DryRun"

	// AppLabelKey is the standard pod label used to select pods backing a Talos
	// control plane StatefulSet/Service.
	AppLabelKey = "app"

	// Field index keys for owner-ref lookups
	IndexControlPlaneRefName = "spec.controlPlaneRef.name"
	IndexWorkerRefName       = "spec.workerRef.name"

	// DeletionPolicyReset is the deletion policy that triggers a Talos reset
	DeletionPolicyReset = "reset"

	// PXE boot stack

	// PXE boot stack enabled value
	PxeBootStackEnabled = "true"

	// proc related paths
	ProcPath        = "/proc"
	ProcCmdlineFile = "cmdline"
	DnsmasqCmdline  = "/sbin/tini\u0000--\u0000/usr/bin/dnsmasq.sh\u0000"

	// dnsmasq configuration path
	DnsmasqConfigPath = "/etc/dnsmasq.d/dnsmasq.conf"
	// Default dnsmasq configuration that disables DNS
	DefaultDnsmasqConfig = "port=0"
	// TFTP files
	TftpDir          = "/var/lib/tftp"
	IpxeEfiX8664File = "ipxe-efi-x86_64.efi"
	IpxeEfiArm64File = "ipxe-efi-arm64.efi"
	IpxeEfiX8664Arch = "x86_64-efi"
	IpxeEfiArm64Arch = "arm64-efi"
	IpxeDownloadFile = "ipxe.efi"

	// Matchbox configuration directory mount point in the talos-operator container
	MatchboxConfigPath = "/var/lib/matchbox"
	// Matchbox configuration subdirectories
	MatchboxAssetsDir   = "assets"
	MatchboxGroupsDir   = "groups"
	MatchboxProfilesDir = "profiles"
)
