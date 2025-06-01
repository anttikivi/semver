package semver

// Versions attaches the methods of [sort.Interface] to a version slice, sorting
// in increasing order.
type Versions []*Version

// Len is the number of elements in Versions.
func (x Versions) Len() int {
	return len(x)
}

// Less reports whether the element with index i must sort before the element
// with index j.
//
// If both Less(i, j) and Less(j, i) are false, then the elements at index i and
// j are considered equal. Sort may place equal elements in any order in
// the final result, while Stable preserves the original input order of equal
// elements.
//
// Less describes a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true
//     as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be
//     false as well.
func (x Versions) Less(i, j int) bool {
	return Compare(x[i], x[j]) < 0
}

// Swap swaps the elements with indexes i and j.
func (x Versions) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
