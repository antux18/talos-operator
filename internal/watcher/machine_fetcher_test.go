/*
Copyright 2025.

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

package watcher

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	talosv1alpha1 "github.com/alperencelik/talos-operator/api/v1alpha1"
)

const (
	testNamespace      = "default"
	testDesiredVersion = "v1.9.0"
)

func newTestFetcher(t *testing.T, tm *talosv1alpha1.TalosMachine) *TalosMachineFetcher {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := talosv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to build scheme: %v", err)
	}
	c := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&talosv1alpha1.TalosMachine{}).
		WithObjects(tm).
		Build()
	return &TalosMachineFetcher{
		Client:   c,
		Recorder: events.NewFakeRecorder(10),
	}
}

func testMachine(state string, annotations map[string]string) *talosv1alpha1.TalosMachine {
	return &talosv1alpha1.TalosMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "machine-a",
			Namespace:   testNamespace,
			Annotations: annotations,
		},
		Spec:   talosv1alpha1.TalosMachineSpec{Endpoint: "10.0.0.10", Version: testDesiredVersion},
		Status: talosv1alpha1.TalosMachineStatus{State: state},
	}
}

// TestUpdateResourceStatusMarksDrift covers the Available → Drifted transition
// and its guards: only Available machines are marked, and DryRun/Disable
// annotated machines are never written to.
func TestUpdateResourceStatusMarksDrift(t *testing.T) {
	key := types.NamespacedName{Namespace: testNamespace, Name: "machine-a"}
	drifted := MachineState{Version: testDesiredVersion, ConfigDiff: "some diff"}
	versionDrifted := MachineState{Version: "v1.8.0"}
	inSync := MachineState{Version: testDesiredVersion}

	tests := []struct {
		name      string
		machine   *talosv1alpha1.TalosMachine
		observed  MachineState
		wantState string
	}{
		{
			name:      "config diff marks Available machine as Drifted",
			machine:   testMachine(talosv1alpha1.StateAvailable, nil),
			observed:  drifted,
			wantState: talosv1alpha1.StateDrifted,
		},
		{
			name:      "version divergence marks Available machine as Drifted",
			machine:   testMachine(talosv1alpha1.StateAvailable, nil),
			observed:  versionDrifted,
			wantState: talosv1alpha1.StateDrifted,
		},
		{
			name:      "node in sync leaves state untouched",
			machine:   testMachine(talosv1alpha1.StateAvailable, nil),
			observed:  inSync,
			wantState: talosv1alpha1.StateAvailable,
		},
		{
			name:      "machine being converged is not stomped by a stale poll",
			machine:   testMachine(talosv1alpha1.StateInstalling, nil),
			observed:  drifted,
			wantState: talosv1alpha1.StateInstalling,
		},
		{
			name:      "already Drifted machine is not re-patched",
			machine:   testMachine(talosv1alpha1.StateDrifted, nil),
			observed:  drifted,
			wantState: talosv1alpha1.StateDrifted,
		},
		{
			name:      "DryRun machine is never written to",
			machine:   testMachine(talosv1alpha1.StateAvailable, map[string]string{talosv1alpha1.ReconcileModeAnnotation: "DryRun"}),
			observed:  drifted,
			wantState: talosv1alpha1.StateAvailable,
		},
		{
			name:      "Disabled machine is never written to",
			machine:   testMachine(talosv1alpha1.StateAvailable, map[string]string{talosv1alpha1.ReconcileModeAnnotation: "disable"}),
			observed:  drifted,
			wantState: talosv1alpha1.StateAvailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newTestFetcher(t, tt.machine)
			if err := f.UpdateResourceStatus(context.Background(), key, tt.observed); err != nil {
				t.Fatalf("UpdateResourceStatus returned error: %v", err)
			}
			var got talosv1alpha1.TalosMachine
			if err := f.Client.Get(context.Background(), key, &got); err != nil {
				t.Fatalf("failed to get machine: %v", err)
			}
			if got.Status.State != tt.wantState {
				t.Fatalf("expected state %q, got %q", tt.wantState, got.Status.State)
			}
		})
	}
}

// TestUpdateResourceStatusMissingMachine ensures a deleted machine does not
// surface an error on the poll path.
func TestUpdateResourceStatusMissingMachine(t *testing.T) {
	f := newTestFetcher(t, testMachine(talosv1alpha1.StateAvailable, nil))
	missing := types.NamespacedName{Namespace: testNamespace, Name: "gone"}
	if err := f.UpdateResourceStatus(context.Background(), missing, MachineState{ConfigDiff: "diff"}); err != nil {
		t.Fatalf("expected missing machine to be ignored, got: %v", err)
	}
}
