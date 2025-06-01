// Copyright (c) 2025 Antti Kivi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
