package file

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

type OpenFlag uint16

const (
	FlagRead = OpenFlag(1 << iota)
	FlagWrite
	FlagAppend
	FlagCreate
	FlagExcl
	FlagTrunc
	FlagNoFollow
	FlagNoATime
	FlagPath
)
const (
	FlagReadWrite    = FlagRead | FlagWrite
	FlagWalkDefaults = FlagRead | FlagNoFollow | FlagNoATime
)

func (mask OpenFlag) HasReadWrite() bool {
	return mask&FlagReadWrite == FlagReadWrite
}

func (mask OpenFlag) HasRead() bool {
	return mask&FlagRead != 0
}

func (mask OpenFlag) HasWrite() bool {
	return mask&FlagRead != 0
}

func (mask OpenFlag) HasAppend() bool {
	return mask&FlagAppend != 0
}

func (mask OpenFlag) HasCreate() bool {
	return mask&FlagCreate != 0
}

func (mask OpenFlag) HasExcl() bool {
	return mask&FlagExcl != 0
}

func (mask OpenFlag) HasTrunc() bool {
	return mask&FlagTrunc != 0
}

func (mask OpenFlag) HasNoFollow() bool {
	return mask&FlagNoFollow != 0
}

func (mask OpenFlag) HasNoATime() bool {
	return mask&FlagNoATime != 0
}

func (mask OpenFlag) HasPath() bool {
	return mask&FlagPath != 0
}

func (mask OpenFlag) OSFlags() int {
	var flags int
	switch {
	case mask.HasReadWrite():
		flags |= os.O_RDWR
	case mask.HasRead():
		flags |= os.O_RDONLY
	case mask.HasWrite():
		flags |= os.O_WRONLY
	}
	if mask.HasAppend() {
		flags |= os.O_APPEND
	}
	if mask.HasCreate() {
		flags |= os.O_CREATE
	}
	if mask.HasExcl() {
		flags |= os.O_EXCL
	}
	if mask.HasTrunc() {
		flags |= os.O_TRUNC
	}
	if mask.HasNoFollow() {
		flags |= syscall.O_NOFOLLOW
	}
	if mask.HasNoATime() {
		flags |= syscall.O_NOATIME
	}
	if mask.HasPath() {
		flags |= unix.O_PATH
	}
	return flags
}
