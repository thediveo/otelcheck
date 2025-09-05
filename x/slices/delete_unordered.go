// Copyright 2025 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package slices

// DeleteUnordered removes the element s[i] from s in place and not preserving
// the order of elements, returning the modified slice. DeleteUnordered zeroes
// the final slice element that has become unused and removes it from the
// returned slice.
//
// See also [Go Wiki: SliceTricks, Delete without preserving order].
//
// [Go Wiki: SliceTricks, Delete without preserving order]:
// https://go.dev/wiki/SliceTricks#delete-without-preserving-order
func DeleteUnordered[S []E, E any](s S, i int) S {
	s[i] = s[len(s)-1]
	var zero E
	s[len(s)-1] = zero
	return s[:len(s)-1]
}
