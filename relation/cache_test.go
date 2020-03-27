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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
	"testing"
)

type Object struct {
	ID string
	SubObjects []*Object
}

func ObjectKey(obj interface{}) (string ,error) {
	o := obj.(*Object)
	return o.ID, nil
}

func ObjectRefers(obj interface{}) ([]string, error) {
	o := obj.(*Object)
	list := make([]string, 0, len(o.SubObjects))
	for _, sub := range o.SubObjects {
		list = append(list, sub.ID)
	}
	return list, nil
}

var objectCache = &cache {
	NewThreadSafeMap(ObjectRefers),
	ObjectKey,
}

func TestRelationCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Relation Cache Suite")
}

var _ = Describe("Unit test", func() {
	c := NewCache(ObjectKey, ObjectRefers)
	subObj1 := &Object{
		ID:         "sub_object1",
		SubObjects: nil,
	}
	obj1 := &Object{
		ID:         "object1",
		SubObjects: []*Object{
			subObj1,
		},
	}
	BeforeEach(func() {
		By("Add object with refers", func() {
			Expect(c.Add(obj1)).ShouldNot(HaveOccurred())
			referenced, err := c.ReferencedKeys(subObj1.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(referenced).Should(Equal([]string{obj1.ID}))
		})
		By("Add sub object", func() {
			Expect(c.Add(subObj1)).ShouldNot(HaveOccurred())
			referenced, err := c.ReferencedKeys(subObj1.ID)
			Expect(referenced).Should(Equal([]string{obj1.ID}))
			refers, err := c.ReferKeys(obj1.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(refers).Should(Equal([]string{subObj1.ID}))
		})
	})
	AfterEach(func() {
		c = NewCache(ObjectKey, ObjectRefers)
	})


	It("Get object", func() {
		item, exists, err := c.Get(obj1)
		item2, _, _ := c.GetByKey(obj1.ID)
		Expect(item).Should(Equal(item2))
		Expect(exists).Should(BeTrue())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(item).Should(Equal(obj1))

	})

	It("List object", func() {
		list := c.List()
		Expect(list).Should(Or(Equal([]interface{}{obj1, subObj1}), Equal([]interface{}{subObj1, obj1})))
		listKeys := c.ListKeys()
		Expect(listKeys).Should(Or(Equal([]string{obj1.ID, subObj1.ID}), Equal([]string{subObj1.ID, obj1.ID})))
	})

	It("Delete sub object", func() {
		Expect(c.Delete(subObj1)).ShouldNot(HaveOccurred())
		_, err := c.ReferencedKeys(subObj1.ID)
		Expect(err).Should(HaveOccurred())
	})

	It("Delete root object", func() {
		Expect(c.Delete(obj1)).ShouldNot(HaveOccurred())
		_, err := c.ReferKeys(obj1.ID)
		Expect(err).Should(HaveOccurred())
		listKeys, err := c.ReferencedKeys(subObj1.ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(listKeys)).Should(Equal(0))
	})

	It("Update sub object", func() {
		subObj := *subObj1
		Expect(c.Add(&subObj)).ShouldNot(HaveOccurred())
		references, err := c.ReferencedKeys(subObj.ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(references).Should(Equal([]string{obj1.ID}))
	})

	It("Get reference", func() {
		references, err := c.Referenced(subObj1)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(references).Should(Equal([]interface{}{obj1}))
		references, err = c.Referenced(obj1)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(references).Should(BeEmpty())
	})

	It("Multiple reference", func() {
		obj2 := &Object{ID: "obj2", SubObjects: []*Object{subObj1}}
		Expect(c.Add(obj2)).ShouldNot(HaveOccurred())
		references, err := c.ReferencedKeys(subObj1.ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(references).Should(Or(Equal([]string{obj1.ID, obj2.ID}), Equal([]string{obj2.ID, obj1.ID})))
	})

	It("Get refers", func() {
		keys, err := c.ReferKeys(obj1.ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(keys).Should(Equal([]string{subObj1.ID}))
	})

	It("Random add and delete", func() {
		m := make(map[string]bool)
		c.Delete(obj1)
		for i := 1; i < 100; i++ {
			i := rand.Int()
			id := strconv.Itoa(i)
			m[id] = true
			obj := &Object{
				ID:         id,
				SubObjects: []*Object{subObj1},
			}
			c.Add(obj)
		}
		referenced, _ := c.ReferencedKeys(subObj1.ID)
		Expect(len(referenced)).Should(Equal(len(m)))
		for id := range m {
			obj := &Object{
				ID:         id,
				SubObjects: nil,
			}
			c.Delete(obj)
			_, exists, _ := c.Get(obj)
			Expect(exists).Should(BeFalse())
		}
		referenced, _ = c.ReferencedKeys(subObj1.ID)
		Expect(len(referenced)).Should(Equal(0))
	})
})



