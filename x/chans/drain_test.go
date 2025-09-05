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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("channels", func() {

	It("drains nothing from a closed channel", func(ctx context.Context) {
		ch := make(chan int, 10)
		close(ch)
		Expect(All(ctx, ch)).To(BeEmpty())
	})

	It("drains nothing with a done context", func(ctx context.Context) {
		ch := make(chan int, 10)
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		Expect(All(ctx, ch)).To(BeEmpty())
	})

	It("drains completely", func(ctx context.Context) {
		ch := make(chan int, 10)
		go func() {
			for i := range 5 {
				ch <- i
			}
			close(ch)
		}()
		Expect(All(ctx, ch)).To(ConsistOf(0, 1, 2, 3, 4))
	})

})
