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

package field_test

import (
	"reflect"
	"testing"

	"github.com/spoof/go-grafana/pkg/field"
)

func TestTags_Set(t *testing.T) {
	ts := []struct {
		initial   []string
		tagsToSet []string
		expected  []string
	}{
		{[]string{"tag1"}, []string{"tag2", "tag3"}, []string{"tag2", "tag3"}},
		{[]string{"tag1"}, []string{}, []string{}},
	}

	for _, tt := range ts {
		tags := field.NewTags(tt.initial...)
		tags.Set(tt.tagsToSet...)
		if !reflect.DeepEqual(tags.Value(), tt.expected) {
			t.Errorf("Tags{tags: %v}.Set(%v): expected %v, got %v", tt.initial, tt.tagsToSet, tt.expected, tags.Value())
		}
	}
}
func TestTags_Add(t *testing.T) {
	ts := []struct {
		initial   []string
		tagsToAdd []string
		expected  []string
	}{
		{[]string{"tag1"}, []string{"tag2", "tag3"}, []string{"tag1", "tag2", "tag3"}},
		{[]string{"tag1"}, []string{"tag1", "tag2"}, []string{"tag1", "tag2"}},
	}

	for _, tt := range ts {
		tags := field.NewTags(tt.initial...)
		tags.Add(tt.tagsToAdd...)
		if !reflect.DeepEqual(tags.Value(), tt.expected) {
			t.Errorf("Tags{tags: %v}.Add(%v): expected %v, got %v", tt.initial, tt.tagsToAdd, tt.expected, tags.Value())
		}
	}
}

func TestTags_Remove(t *testing.T) {
	ts := []struct {
		initial      []string
		tagsToRemove []string
		expected     []string
	}{
		{[]string{"tag1"}, []string{"tag1"}, []string{}},
		{[]string{"tag1", "tag2"}, []string{"tag2"}, []string{"tag1"}},
		{[]string{"tag1", "tag2"}, []string{"tag3"}, []string{"tag1", "tag2"}},
	}

	for _, tt := range ts {
		tags := field.NewTags(tt.initial...)
		tags.Remove(tt.tagsToRemove...)
		if !reflect.DeepEqual(tags.Value(), tt.expected) {
			t.Errorf("Tags{tags: %v}.Remove(%v): expected %v, got %v", tt.initial, tt.tagsToRemove, tt.expected, tags.Value())
		}
	}
}
