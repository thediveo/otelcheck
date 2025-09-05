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

package chans

import (
	"context"
	"iter"
)

// All returns an iterator iterating over all elements from channel c, ending
// when either c has been closed and the last element read from c, or when the
// passed context is done.
//
// The passed context allows the iteration to be still bounded even if the
// passed channel never gets closed. Please note that Ginkgo v2 features
// [interruptible test specifications] that accept a [context.Context] argument.
//
// Please note that when calling All with a done context as well as at least one
// element available from the channel, there is no guarantee whether this
// element will be drained or not.
//
// [interruptible test specifications]: https://onsi.github.io/ginkgo/#spec-timeouts-and-interruptible-nodes
func All[E any](ctx context.Context, c <-chan E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			select {
			case e, ok := <-c:
				if !ok || !yield(e) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
