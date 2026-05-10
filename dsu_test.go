package disjointset_test

import (
	"testing"

	"github.com/manavrajvanshi/disjointset"
)

func TestNewDisjointset(t *testing.T) {
	d := disjointset.NewDSU[int]()
	if d == nil {
		t.Error("expected non-nil DSU")
	}
	if d.Components() != 0 {
		t.Errorf("expected 0 components, got %d", d.Components())
	}
}

func TestDSU(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		tests := []struct {
			name           string
			adds           []int
			wantReturns    []bool
			wantComponents int
		}{
			{
				name:           "single node",
				adds:           []int{1},
				wantReturns:    []bool{true},
				wantComponents: 1,
			},
			{
				name:           "duplicate node returns false",
				adds:           []int{1, 1},
				wantReturns:    []bool{true, false},
				wantComponents: 1,
			},
			{
				name:           "multiple nodes unique",
				adds:           []int{1, 2, 3},
				wantReturns:    []bool{true, true, true},
				wantComponents: 3,
			},
			{
				name:           "multiple nodes with repeats",
				adds:           []int{1, 2, 3, 4, 3, 2, 1},
				wantReturns:    []bool{true, true, true, true, false, false, false},
				wantComponents: 4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := disjointset.NewDSU[int]()
				for i, n := range tt.adds {
					got := d.Add(n)
					if got != tt.wantReturns[i] {
						t.Errorf("Add(%d) = %v, want %v", n, got, tt.wantReturns[i])
					}
				}
				if d.Components() != tt.wantComponents {
					t.Errorf("got %d components, want %d", d.Components(), tt.wantComponents)
				}
			})
		}
	})

	t.Run("Root", func(t *testing.T) {
		tests := []struct {
			name      string
			adds      []int
			unions    [][2]int
			node      int
			wantRoot  int
			wantError bool
		}{
			{
				name:      "unregistered node",
				adds:      []int{},
				node:      1,
				wantError: true,
			},
			{
				name:     "root of one of many nodes",
				adds:     []int{1, 2, 3, 4, 5},
				node:     5,
				wantRoot: 5,
			},
			{
				name:     "root after union higher rank becomes root",
				adds:     []int{1, 2, 3, 4, 5},
				unions:   [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
				node:     5,
				wantRoot: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := disjointset.NewDSU[int]()
				for _, n := range tt.adds {
					d.Add(n)
				}
				for _, u := range tt.unions {
					d.Union(u[0], u[1])
				}
				root, err := d.Root(tt.node)
				if (err != nil) != tt.wantError {
					t.Errorf("got error %v, wantError %v", err, tt.wantError)
				}
				if err == nil && root != tt.wantRoot {
					t.Errorf("got root %d, want %d", root, tt.wantRoot)
				}
			})
		}
	})

	t.Run("Find", func(t *testing.T) {
		tests := []struct {
			name      string
			adds      []int
			unions    [][2]int
			p, q      int
			wantFound bool
			wantError bool
		}{
			{
				name:      "unregistered node",
				adds:      []int{1},
				p:         1,
				q:         2,
				wantError: true,
			},
			{
				name:      "both unregistered",
				adds:      []int{},
				p:         1,
				q:         2,
				wantError: true,
			},
			{
				name:      "same node",
				adds:      []int{1},
				p:         1,
				q:         1,
				wantFound: true,
			},
			{
				name:      "different components",
				adds:      []int{1, 2},
				p:         1,
				q:         2,
				wantFound: false,
			},
			{
				name:      "same component after union",
				adds:      []int{1, 2},
				unions:    [][2]int{{1, 2}},
				p:         1,
				q:         2,
				wantFound: true,
			},
			{
				name:      "transitive union",
				adds:      []int{1, 2, 3},
				unions:    [][2]int{{1, 2}, {2, 3}},
				p:         1,
				q:         3,
				wantFound: true,
			},
			{
				name:      "not connected in partial graph",
				adds:      []int{1, 2, 3, 4},
				unions:    [][2]int{{1, 2}, {3, 4}},
				p:         1,
				q:         3,
				wantFound: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := disjointset.NewDSU[int]()
				for _, n := range tt.adds {
					d.Add(n)
				}
				for _, u := range tt.unions {
					d.Union(u[0], u[1])
				}
				found, err := d.Find(tt.p, tt.q)
				if (err != nil) != tt.wantError {
					t.Errorf("got error %v, wantError %v", err, tt.wantError)
				}
				if err == nil && found != tt.wantFound {
					t.Errorf("got found %v, want %v", found, tt.wantFound)
				}
			})
		}
	})

	t.Run("Union", func(t *testing.T) {
		tests := []struct {
			name           string
			adds           []int
			unions         [][2]int
			p, q           int
			wantMerged     bool
			wantError      bool
			wantComponents int
		}{
			{
				name:      "unregistered node",
				adds:      []int{1},
				p:         1,
				q:         2,
				wantError: true,
			},
			{
				name:      "both unregistered",
				adds:      []int{},
				p:         1,
				q:         2,
				wantError: true,
			},
			{
				name:           "first union",
				adds:           []int{1, 2},
				p:              1,
				q:              2,
				wantMerged:     true,
				wantComponents: 1,
			},
			{
				name:           "already connected",
				adds:           []int{1, 2},
				unions:         [][2]int{{1, 2}},
				p:              1,
				q:              2,
				wantMerged:     false,
				wantComponents: 1,
			},
			{
				name:           "three nodes two unions",
				adds:           []int{1, 2, 3},
				unions:         [][2]int{{1, 2}},
				p:              2,
				q:              3,
				wantMerged:     true,
				wantComponents: 1,
			},
			{
				name:           "union same node with itself",
				adds:           []int{1},
				p:              1,
				q:              1,
				wantMerged:     false,
				wantComponents: 1,
			},
			{
				name:           "transitive already connected",
				adds:           []int{1, 2, 3},
				unions:         [][2]int{{1, 2}, {2, 3}},
				p:              1,
				q:              3,
				wantMerged:     false,
				wantComponents: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := disjointset.NewDSU[int]()
				for _, n := range tt.adds {
					d.Add(n)
				}
				for _, u := range tt.unions {
					d.Union(u[0], u[1])
				}
				merged, err := d.Union(tt.p, tt.q)
				if (err != nil) != tt.wantError {
					t.Errorf("got error %v, wantError %v", err, tt.wantError)
				}
				if err == nil {
					if merged != tt.wantMerged {
						t.Errorf("got merged %v, want %v", merged, tt.wantMerged)
					}
					if d.Components() != tt.wantComponents {
						t.Errorf("got %d components, want %d", d.Components(), tt.wantComponents)
					}
				}
			})
		}
	})

	t.Run("Components", func(t *testing.T) {
		tests := []struct {
			name           string
			adds           []int
			unions         [][2]int
			wantComponents int
		}{
			{
				name:           "empty DSU",
				adds:           []int{},
				wantComponents: 0,
			},
			{
				name:           "one node",
				adds:           []int{1},
				wantComponents: 1,
			},
			{
				name:           "three nodes no unions",
				adds:           []int{1, 2, 3},
				wantComponents: 3,
			},
			{
				name:           "three nodes one union",
				adds:           []int{1, 2, 3},
				unions:         [][2]int{{1, 2}},
				wantComponents: 2,
			},
			{
				name:           "three nodes all unioned",
				adds:           []int{1, 2, 3},
				unions:         [][2]int{{1, 2}, {2, 3}},
				wantComponents: 1,
			},
			{
				name:           "duplicate union doesnt change components",
				adds:           []int{1, 2},
				unions:         [][2]int{{1, 2}, {1, 2}},
				wantComponents: 1,
			},
			{
				name:           "two separate components",
				adds:           []int{1, 2, 3, 4},
				unions:         [][2]int{{1, 2}, {3, 4}},
				wantComponents: 2,
			},
			{
				name:           "five nodes all unioned",
				adds:           []int{1, 2, 3, 4, 5},
				unions:         [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
				wantComponents: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := disjointset.NewDSU[int]()
				for _, n := range tt.adds {
					d.Add(n)
				}
				for _, u := range tt.unions {
					d.Union(u[0], u[1])
				}
				if d.Components() != tt.wantComponents {
					t.Errorf("got %d components, want %d", d.Components(), tt.wantComponents)
				}
			})
		}
	})
}
