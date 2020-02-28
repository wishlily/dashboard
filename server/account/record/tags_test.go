package record

import (
	"strings"
	"testing"
)

func TestTagKey(t *testing.T) {
	testCfgPre()
	testCfgClr(t)
	var tp Tag
	if v := tp.key("sec.v1", "test"); !strings.EqualFold(v, "test") {
		t.Fatalf("key should be `test`")
	}
	testCfgSet(t)
	if v := tp.key("csv.tags", "member"); !strings.EqualFold(v, "成员") {
		t.Fatalf("key should be `成员`: %v", v)
	}
}

func TestTagString(t *testing.T) {
	testCfgClr(t)
	for i, tc := range []struct {
		a Tag
		v string
	}{
		{TagMember, "member"},
		{TagProj, "proj"},
		{TagUnit, "unit"},
		{TagNUV, "nuv"},
		{TagDeadline, "deadline"},
		{99, "unkown"},
	} {
		if !strings.EqualFold(tc.a.String(), tc.v) {
			t.Fatalf("%d: should be %v: %v", i, tc.v, tc.a.String())
		}
	}
	testCfgSet(t)
	for i, tc := range []struct {
		a Tag
		v string
	}{
		{TagMember, "成员"},
		{TagProj, "项目"},
		{TagUnit, "份额"},
		{TagNUV, "净值"},
		{TagDeadline, "截至日期"},
		{99, "其他"},
	} {
		if !strings.EqualFold(tc.a.String(), tc.v) {
			t.Fatalf("%d: should be %v: %v", i, tc.v, tc.a.String())
		}
	}
}

func TestTagParse(t *testing.T) {
	testCfgClr(t)
	for i, tc := range []struct {
		a string
		v Tag
	}{
		{"member", TagMember},
		{"proj", TagProj},
		{"unit", TagUnit},
		{"nuv", TagNUV},
		{"unkown", tagEnd},
	} {
		if v := ParseTag(tc.a); tc.v != v {
			t.Fatalf("%d: should be %v: %v", i, tc.v, v)
		}
	}
	testCfgSet(t)
	for i, tc := range []struct {
		a string
		v Tag
	}{
		{"成员", TagMember},
		{"项目", TagProj},
		{"份额", TagUnit},
		{"净值", TagNUV},
		{"其他", tagEnd},
	} {
		if v := ParseTag(tc.a); tc.v != v {
			t.Fatalf("%d: should be %v: %v", i, tc.v, v)
		}
	}
}
