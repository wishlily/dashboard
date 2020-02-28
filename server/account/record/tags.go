package record

import (
	"strings"

	"github.com/spf13/viper"
)

// Tag data type: Item Note extra info
type Tag int

const (
	// TagMember : member name
	TagMember Tag = iota
	// TagProj : project name
	TagProj
	// TagUnit : fund units
	TagUnit
	// TagNUV : account nuv
	TagNUV
	// TagDeadline : deadline time
	TagDeadline
	tagEnd
)

func (t Tag) key(section, key string) string {
	url := section + "." + key
	v := viper.GetString(url)
	if viper.IsSet(url) && len(v) != 0 {
		return v
	}
	return key
}

func (t Tag) String() string {
	const section = "csv.tags"
	switch t {
	case TagMember:
		return t.key(section, "member")
	case TagProj:
		return t.key(section, "proj")
	case TagUnit:
		return t.key(section, "unit")
	case TagNUV:
		return t.key(section, "nuv")
	case TagDeadline:
		return t.key(section, "deadline")
	}
	return t.key(section, "unkown")
}

// ParseTag : Parse string to Tag
func ParseTag(v string) Tag {
	var t Tag
	for t = TagMember; t < tagEnd; t++ {
		if strings.EqualFold(v, t.String()) {
			return t
		}
	}
	return tagEnd
}
