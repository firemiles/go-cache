/*
 * Copyright (c) 2020 firemiles(miles.dev@outlook.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package relation

import (
	"strconv"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache add and delete performance test", func() {
	var c Cache
	BeforeEach(func() {
		c = NewCache(ObjectKey, ObjectRefers)
	})

	Measure("Add 100000 object one goroutine", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < 100000; i++ {
				c.Add(&Object{ID: strconv.Itoa(i)})
			}
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.3), "add 100000 object should't take too long")
	}, 10)
	Measure("Add 100000 object with dependence one goroutine", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < 100000; i++ {
				c.Add(&Object{ID: strconv.Itoa(i), SubObjects: []*Object{&Object{ID: "a"}}})
			}
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.6), "add 100000 object should't take too long")
	}, 10)

	Measure("Delete 100000 object one goroutine", func(b Benchmarker) {
		for i := 0; i < 100000; i++ {
			c.Add(&Object{ID: strconv.Itoa(i)})
		}
		runtime := b.Time("runtime", func() {
			for i := 0; i < 100000; i++ {
				c.Delete(&Object{ID: strconv.Itoa(i)})
			}
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.2),
			"delete 100000 object should't take too long")
	}, 10)

	Measure("Delete 100000 object with dependence one goroutine", func(b Benchmarker) {
		for i := 0; i < 100000; i++ {
			c.Add(&Object{ID: strconv.Itoa(i), SubObjects: []*Object{&Object{ID: "a"}}})
		}
		runtime := b.Time("runtime", func() {
			for i := 0; i < 100000; i++ {
				c.Delete(&Object{ID: strconv.Itoa(i)})
			}
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.4),
			"delete 100000 object should't take too long")
	}, 10)

	Measure("Add 100000 object 4 goroutines", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			var wait sync.WaitGroup
			for i := 0; i < 4; i++ {
				wait.Add(1)
				go func(i int) {
					start, end := i*25000, (i+1)*25000
					for j := start; j < end; j++ {
						c.Add(&Object{ID: strconv.Itoa(j)})
					}
					wait.Done()
				}(i)
			}
			wait.Wait()
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.3), "add 100000 object 4 goroutines  should't take too long")
	}, 10)

	Measure("sync.Map add 100000 object 4 goroutines", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			var wait sync.WaitGroup
			m := sync.Map{}
			for i := 0; i < 4; i++ {
				wait.Add(1)
				go func(i int) {
					start, end := i*25000, (i+1)*25000
					for j := start; j < end; j++ {
						m.Store(strconv.Itoa(j), &Object{})
					}
					wait.Done()
				}(i)
			}
			wait.Wait()
		})
		Expect(runtime.Seconds()).Should(BeNumerically("<", 0.3), "add 100000 object 4 goroutines  should't take too long")
	}, 10)
})
