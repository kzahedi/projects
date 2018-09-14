package main

import "fmt"

// Tweet contains the tree
type Tweet struct {
	ID            int
	ParentID      int
	Name          string
	TwitterHandle string
	Text          string
	Date          string
	Replies       int
	Retweets      int
	Likes         int
	Link          string
	Type          string
	Mentions      []string
	Children      []Tweet
}

func makeTweet() Tweet {
	return Tweet{ID: -1,
		Name:          "",
		TwitterHandle: "",
		Type:          "",
		Text:          "",
		Date:          "",
		Link:          "",
		Replies:       0,
		Retweets:      0,
		Likes:         0,
		Children:      nil,
		ParentID:      -1,
		Mentions:      make([]string, 0, 0),
	}
}

func (t Tweet) String() string {
	s := ""
	s = fmt.Sprintf("%sName: %s\n", s, t.Name)
	s = fmt.Sprintf("%sTwitter Handle: %s\n", s, t.TwitterHandle)
	s = fmt.Sprintf("%sText: \"%s\"\n", s, t.Text)
	s = fmt.Sprintf("%sReplies: %d\n", s, t.Replies)
	s = fmt.Sprintf("%sLikes: %d\n", s, t.Likes)
	s = fmt.Sprintf("%sRetweets: %d\n", s, t.Retweets)
	s = fmt.Sprintf("%sParent: %d\n", s, t.ParentID)
	s = fmt.Sprintf("%sID: %d\n", s, t.ID)
	s = fmt.Sprintf("%sMentions:\n", s)
	for _, v := range t.Mentions {
		s = fmt.Sprintf("%s  %s\n", s, v)
	}
	for _, t := range t.Children {
		s = fmt.Sprintf("%s%s", s, t)
	}
	return s
}

// ValidTwitterHandle checks if this tweet is corrupted
func (t Tweet) ValidTwitterHandle() bool {
	return t.TwitterHandle != ""
}

// Valid return true if the tree is valid
func (t Tweet) Valid() bool {
	if t.ValidTwitterHandle() == false {
		return false
	}
	for _, c := range t.Children {
		if c.Valid() == false {
			return false
		}
	}
	return true
}
