package semver

// Versions attaches the methods of [sort.Interface] to []*Version, sorting in
// increasing order.
type Versions []*Version

// Len is the number of elements in Versions.
func (x Versions) Len() int {
	return len(x)
}

// Swap swaps the elements with indexes i and j.
func (x Versions) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
