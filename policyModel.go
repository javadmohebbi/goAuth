package goAuth

import "gorm.io/gorm"

type Policy struct {
	Section string `gorm:"type:nvarchar(100);not null"`

	// Permission Description
	// ----------------------------------------------
	//   #   Permission           rwud*      Binary
	// ----------------------------------------------
	//   0   none                 ----       0000
	//   1                        ---d       0001
	//   2                        --u-       0010
	//   3                        --ud       0011
	//   4                        -w--       0100
	//   5                        -w-d       0101
	//   6                        -wu-       0110
	//   7                        -wud       0111
	//   8                        r---       1000
	//   9                        r--d       1001
	//   10                       r-u-       1010
	//   11                       r-ud       1011
	//   12                       rw--       1100
	//   13                       rw-d       1101
	//   14                       rwu-       1110
	//   15                       rwud       1111
	// ----------------------------------------------
	// *rwdu => Read Write Update Delete
	// ----------------------------------------------
	Perm uint

	// Group ID
	GroupID uint

	gorm.Model
}
