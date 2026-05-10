package disjointset

import "errors"

const (
	// defaultRank is the initial rank assigned to every newly added node.
	// Union by rank uses this to decide which node becomes the parent.
	defaultRank = 1
)

// DSU is a generic Disjoint Set Union interface for any comparable type.
// It supports union by rank and path compression for near-constant time operations.
type DSU[T comparable] interface {
	// Union merges the components containing p and q.
	// Returns true if they were in different components, false if already connected.
	// Returns an error if either node is not registered.
	Union(p, q T) (bool, error)

	// Find reports whether p and q belong to the same component.
	// Returns an error if either nodse is not registered.
	Find(p, q T) (bool, error)

	// Root returns the representative node of the component containing p.
	// Returns an error if p is not registered.
	Root(p T) (T, error)

	// Add registers a new node. Returns true if added, false if already present.
	Add(p T) bool

	// Components returns the number of disjoint sets.
	Components() int
}

// dsu is the internal implementation of DSU.
type dsu[T comparable] struct {
	parent     map[T]T   // maps each node to its parent; root nodes point to themselves
	rank       map[T]int // tracks the rank of each node for union by rank
	components int       // number of disjoint sets
}

// NewDSU creates and returns an empty DSU.
func NewDSU[T comparable]() DSU[T] {
	return &dsu[T]{
		parent:     make(map[T]T),
		rank:       make(map[T]int),
		components: 0,
	}
}

// Union merges the components of p and q using union by rank.
// The root with the higher rank becomes the parent; ranks are summed on merge.
// Returns false without error if p and q are already in the same component.
func (d *dsu[T]) Union(p, q T) (bool, error) {
	rootP, err := d.Root(p)
	if err != nil {
		return false, err
	}
	rootQ, err := d.Root(q)
	if err != nil {
		return false, err
	}

	// already in the same component, nothing to do
	if rootP == rootQ {
		return false, nil
	}

	rankP, rankQ := d.rank[rootP], d.rank[rootQ]

	// higher rank becomes the parent to keep the tree shallow
	var parent, child T
	if rankP >= rankQ {
		parent, child = rootP, rootQ
	} else {
		parent, child = rootQ, rootP
	}

	d.parent[child] = parent
	d.rank[parent] = rankP + rankQ

	d.components--

	return true, nil
}

// Find reports whether p and q share the same root.
func (d *dsu[T]) Find(p, q T) (bool, error) {
	rootP, err := d.Root(p)
	if err != nil {
		return false, err
	}

	rootQ, err := d.Root(q)
	if err != nil {
		return false, err
	}

	return rootP == rootQ, nil
}

// Root walks up the parent chain to find the representative of p's component,
// then calls pathCompress to flatten the tree for future lookups.
func (d *dsu[T]) Root(p T) (T, error) {
	if _, ok := d.parent[p]; !ok {
		var zero T
		return zero, errors.New("element not found")
	}

	root := p

	// walk up until we reach a node that is its own parent
	for root != d.parent[root] {
		root = d.parent[root]
	}

	// flatten the path so future Root calls are faster
	d.pathCompress(p, root)

	return root, nil
}

// Add registers p as a new node in its own component.
// Returns false if p is already registered.
func (d *dsu[T]) Add(p T) bool {
	if _, ok := d.parent[p]; ok {
		return false
	}

	// new node is its own parent and gets the default rank
	d.parent[p] = p
	d.rank[p] = defaultRank
	d.components++

	return true
}

// Components returns the current number of disjoint sets.
func (d *dsu[T]) Components() int {
	return d.components
}

// pathCompress flattens the path from node to root by pointing
// every node along the way directly to the root.
func (d *dsu[T]) pathCompress(node, root T) {
	for node != d.parent[node] {
		parent := d.parent[node]
		d.parent[node] = root
		node = parent
	}
}
