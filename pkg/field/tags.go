// Copyright 2017 Sergey Safonov
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package field

// Tags is slice of string that preserv uniqueness of values.
type Tags struct {
	tags []string
}

func NewTags(tags ...string) *Tags {
	t := &Tags{}
	t.Set(tags...)

	return t
}

// Value returns stored tags
func (t *Tags) Value() []string {
	return t.tags
}

// Set sets new tags.
func (t *Tags) Set(tags ...string) {
	newTags := []string{}
	uniqTags := make(map[string]bool)
	for _, tag := range tags {
		if _, ok := uniqTags[tag]; ok {
			continue
		}

		uniqTags[tag] = true
		newTags = append(newTags, tag)
	}

	t.tags = newTags
}

// Add adds given tags. This method keeps uniqueness of tags.
func (t *Tags) Add(tags ...string) {
	tagFound := make(map[string]bool, len(t.tags))
	for _, tag := range t.tags {
		tagFound[tag] = true
	}

	for _, tag := range tags {
		if _, ok := tagFound[tag]; ok {
			continue
		}
		t.tags = append(t.tags, tag)
	}
}

// Remove removes given tags. Does nothing if tag is not found.
func (t *Tags) Remove(tags ...string) {
	tagIndex := make(map[string]int, len(t.tags))
	for i, tag := range t.tags {
		tagIndex[tag] = i
	}

	for _, tag := range tags {
		if i, ok := tagIndex[tag]; ok {
			t.tags = append(t.tags[:i], t.tags[i+1:]...)
		}
	}
}
