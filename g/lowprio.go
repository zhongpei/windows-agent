package g

// Copyright (C) 2018 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

import (
	"syscall"
	"github.com/pkg/errors"

)

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686219(v=vs.85).aspx
	aboveNormalPriorityClass   = 0x00008000
	belowNormalPriorityClass   = 0x00004000
	highPriorityClass          = 0x00000080
	idlePriorityClass          = 0x00000040
	normalPriorityClass        = 0x00000020
	processModeBackgroundBegin = 0x00100000
	processModeBackgroundEnd   = 0x00200000
	realtimePriorityClass      = 0x00000100
)

const (
	aboveNormalPriorityClassFlag = "aboveNormal"
	belowNormalPriorityClassFlag = "belowNormal"
	highPriorityClassFlag        = "high"
	idlePriorityClassFlag        = "idle"
	normalPriorityClassFlag      = "normal"
	realtimePriorityClassFlag    = "realtime"
)

// SetLowPriority lowers the process CPU scheduling priority, and possibly
// I/O priority depending on the platform and OS.
func setLowPriority(level uintptr) error {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	setPriorityClass := modkernel32.NewProc("SetPriorityClass")

	if err := setPriorityClass.Find(); err != nil {
		return errors.Wrap(err, "find proc")
	}

	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return errors.Wrap(err, "get process handler")
	}
	defer syscall.CloseHandle(handle)

	res, _, err := setPriorityClass.Call(uintptr(handle), level)
	if res != 0 {
		// "If the function succeeds, the return value is nonzero."
		return nil
	}
	return errors.Wrap(err, "set priority class") // wraps nil as nil
}

func initPrio(prio string) error {
	switch(prio) {
	case aboveNormalPriorityClassFlag:
		return setLowPriority(aboveNormalPriorityClass)
	case belowNormalPriorityClassFlag:
		return setLowPriority(belowNormalPriorityClass)
	case highPriorityClassFlag:
		return setLowPriority(highPriorityClass)
	case idlePriorityClassFlag:
		return setLowPriority(idlePriorityClass)
	case normalPriorityClassFlag:
		return setLowPriority(normalPriorityClass)
	case realtimePriorityClassFlag:
		return setLowPriority(realtimePriorityClass)

	}
	return errors.New("priority not found")
}
