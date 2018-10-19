package message

import (
	"sort"
	"testing"
)

func TestUnread(t *testing.T) {
	cases := []struct {
		in, want Messages
		index    int
	}{
		{
			make(Messages, 0),
			make(Messages, 0),
			0,
		},
		{
			Messages{Message{Index: 0}, Message{Index: 1}},
			Messages{Message{Index: 1}},
			0,
		},
		{
			Messages{Message{Index: 0}, Message{Index: 1}},
			make(Messages, 0),
			1,
		},
		{
			Messages{Message{Index: 0}, Message{Index: 1}, Message{Index: 2}, Message{Index: 3}},
			Messages{Message{Index: 2}, Message{Index: 3}},
			1,
		},
	}

	for _, c := range cases {
		unread := c.in.unread(c.index)
		for i := range unread {
			if c.want[i].Index != unread[i].Index {
				t.Errorf("Expected '%v' unread for index '%d' to be '%v', got '%v'.", c.in, c.index, c.want, unread)
				break
			}
		}
	}
}

func TestSort(t *testing.T) {
	cases := []struct {
		in, want Messages
	}{
		{
			Messages{Message{Index: 3}, Message{Index: 2}, Message{Index: 1}},
			Messages{Message{Index: 1}, Message{Index: 2}, Message{Index: 3}},
		},
		{
			Messages{Message{Index: 3}, Message{Index: 1}, Message{Index: 2}},
			Messages{Message{Index: 1}, Message{Index: 2}, Message{Index: 3}},
		},
		{
			Messages{Message{Index: 1}, Message{Index: 2}, Message{Index: 3}},
			Messages{Message{Index: 1}, Message{Index: 2}, Message{Index: 3}},
		},
		{
			Messages{Message{Index: 2}, Message{Index: 1}, Message{Index: 1}},
			Messages{Message{Index: 1}, Message{Index: 1}, Message{Index: 2}},
		},
		{
			Messages{Message{Index: 1}, Message{Index: 1}, Message{Index: 1}},
			Messages{Message{Index: 1}, Message{Index: 1}, Message{Index: 1}},
		},
		{
			make(Messages, 0),
			make(Messages, 0),
		},
		{
			Messages{Message{Index: 1}},
			Messages{Message{Index: 1}},
		},
	}

	for _, c := range cases {
		sort.Sort(c.in)
		for i := range c.in {
			if c.in[i].Index != c.want[i].Index {
				t.Errorf("Expected messages '%v' to be sorted as '%v'.", c.in, c.want)
				break
			}
		}
	}
}
