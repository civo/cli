package kubernetes

import (
	"errors"
	"reflect"
	"testing"

	"github.com/civo/civogo"
)

func TestResolveNodePoolsToDelete(t *testing.T) {
	pools := []civogo.KubernetesPool{
		{ID: "aaaaaa-1111"},
		{ID: "bbbbbb-2222"},
		{ID: "cccccc-3333"},
	}

	tests := []struct {
		name      string
		requested []string
		want      []string
		wantErr   bool
	}{
		{
			name:      "full ID",
			requested: []string{"bbbbbb-2222"},
			want:      []string{"bbbbbb-2222"},
		},
		{
			name:      "prefix resolves to the full ID",
			requested: []string{"cccccc"},
			want:      []string{"cccccc-3333"},
		},
		{
			name:      "several pools",
			requested: []string{"aaaaaa", "cccccc"},
			want:      []string{"aaaaaa-1111", "cccccc-3333"},
		},
		{
			name:      "unknown ID",
			requested: []string{"zzzzzz"},
			wantErr:   true,
		},
		{
			name:      "one unknown ID among valid ones deletes nothing",
			requested: []string{"aaaaaa", "zzzzzz"},
			wantErr:   true,
		},
		{
			name:      "same pool twice",
			requested: []string{"aaaaaa", "aaaaaa-1111"},
			wantErr:   true,
		},
		{
			name:      "every pool",
			requested: []string{"aaaaaa", "bbbbbb", "cccccc"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveNodePoolsToDelete(pools, tt.requested)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveNodePoolsToDelete(%v) error = %v, wantErr %v", tt.requested, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveNodePoolsToDelete(%v) = %v, want %v", tt.requested, got, tt.want)
			}
		})
	}
}

type fakeNodePoolDeleter struct {
	calls  []string
	failOn string
}

func (f *fakeNodePoolDeleter) DeleteKubernetesClusterPool(id, poolID string) (*civogo.SimpleResponse, error) {
	f.calls = append(f.calls, id+"/"+poolID)
	if poolID == f.failOn {
		return nil, errors.New("boom")
	}
	return &civogo.SimpleResponse{Result: "success"}, nil
}

func TestDeleteNodePools(t *testing.T) {
	t.Run("deletes each pool through the pool endpoint", func(t *testing.T) {
		f := &fakeNodePoolDeleter{}
		deleted, err := deleteNodePools(f, "cluster-1", []string{"p1", "p2"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if want := []string{"p1", "p2"}; !reflect.DeepEqual(deleted, want) {
			t.Errorf("deleted = %v, want %v", deleted, want)
		}
		if want := []string{"cluster-1/p1", "cluster-1/p2"}; !reflect.DeepEqual(f.calls, want) {
			t.Errorf("calls = %v, want %v", f.calls, want)
		}
	})

	t.Run("stops at the first failure and reports what was deleted", func(t *testing.T) {
		f := &fakeNodePoolDeleter{failOn: "p2"}
		deleted, err := deleteNodePools(f, "cluster-1", []string{"p1", "p2", "p3"})
		if err == nil {
			t.Fatal("expected an error")
		}
		if want := []string{"p1"}; !reflect.DeepEqual(deleted, want) {
			t.Errorf("deleted = %v, want %v", deleted, want)
		}
		if len(f.calls) != 2 {
			t.Errorf("expected 2 calls, got %v", f.calls)
		}
	})
}
