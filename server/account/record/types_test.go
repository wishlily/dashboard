package record

import (
	"strings"
	"testing"
)

func TestTypeKey(t *testing.T) {
	testCfgPre()
	testCfgClr(t)
	var tp Type
	if v := tp.key("sec.v1", "test"); !strings.EqualFold(v, "test") {
		t.Fatalf("key should be `test`")
	}
	testCfgSet(t)
	if v := tp.key("csv.types", "I"); !strings.EqualFold(v, "收入") {
		t.Fatalf("key should be `收入`: %v", v)
	}
}

func TestTypeString(t *testing.T) {
	testCfgClr(t)
	for i, tc := range []struct {
		a Type
		v string
	}{
		{TypeI, "I"},
		{TypeO, "O"},
		{TypeR, "R"},
		{TypeB, "B"},
		{TypeL, "L"},
		{TypeX, "X"},
		{99, "U"},
	} {
		if !strings.EqualFold(tc.a.String(), tc.v) {
			t.Fatalf("%d: should be %v: %v", i, tc.v, tc.a.String())
		}
	}
	testCfgSet(t)
	for i, tc := range []struct {
		a Type
		v string
	}{
		{TypeI, "收入"},
		{TypeO, "支出"},
		{TypeR, "转账"},
		{TypeB, "借入"},
		{TypeL, "借出"},
		{TypeX, "修正"},
		{99, "unkownn"},
	} {
		if !strings.EqualFold(tc.a.String(), tc.v) {
			t.Fatalf("%d: should be %v: %v", i, tc.v, tc.a.String())
		}
	}
}

func TestTypeParse(t *testing.T) {
	testCfgClr(t)
	for i, tc := range []struct {
		a string
		v Type
	}{
		{"I", TypeI},
		{"O", TypeO},
		{"R", TypeR},
		{"B", TypeB},
		{"L", TypeL},
		{"X", TypeX},
		{"U", typeEnd},
	} {
		if v := ParseType(tc.a); tc.v != v {
			t.Fatalf("%d: should be %d: %d", i, tc.v, v)
		}
	}
	testCfgSet(t)
	for i, tc := range []struct {
		a string
		v Type
	}{
		{"收入", TypeI},
		{"支出", TypeO},
		{"转账", TypeR},
		{"借入", TypeB},
		{"借出", TypeL},
		{"修正", TypeX},
		{"unkownn", typeEnd},
		{"balabala", typeEnd},
	} {
		if v := ParseType(tc.a); tc.v != v {
			t.Fatalf("%d: should be %d: %d", i, tc.v, v)
		}
	}
}
