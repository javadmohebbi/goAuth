package goAuth

import (
	"strconv"
)

// an struct for goAuth module
// all of the major functionalities of this
// package has impelemented in this struct
type GoAuth struct {
	Policies []GoAuthPolicy `json:"policies"`
}

// struct to store user policies
type GoAuthPolicy struct {
	Section string `json:"section"`
	UGO     UGO    `json:"ugo"`
}

// UGO
type UGO uint

// it will get string int
func (u UGO) binaryString() string {
	// string bytes
	s := strconv.FormatInt(int64(u), 2)
	switch len(s) {
	case 0:
		s = "0000"
	case 1:
		s = "000" + s
	case 2:
		s = "00" + s
	case 3:
		s = "0" + s
	case 4:
		// nothing to do
		break
	default:
		// problem
		s = "0000"
	}

	return s
}

// Get boolean r,w,u,d
func (u UGO) Bools() (bool, bool, bool, bool) {
	bs := u.binaryString()

	if len(bs) != 4 {
		return false, false, false, false
	}

	rr, ww, uu, dd := false, false, false, false

	if bs[0] == '1' {
		rr = true
	}
	if bs[1] == '1' {
		ww = true
	}
	if bs[2] == '1' {
		uu = true
	}
	if bs[3] == '1' {
		dd = true
	}

	return rr, ww, uu, dd
}

// Get boolean r,w,u,d
func (u UGO) String() string {
	bs := u.binaryString()

	if len(bs) != 4 {
		return "----"
	}

	rr, ww, uu, dd := "-", "-", "-", "-"

	if bs[0] == '1' {
		rr = "r"
	}
	if bs[1] == '1' {
		ww = "w"
	}
	if bs[2] == '1' {
		uu = "u"
	}
	if bs[3] == '1' {
		dd = "d"
	}

	return rr + ww + uu + dd
}
