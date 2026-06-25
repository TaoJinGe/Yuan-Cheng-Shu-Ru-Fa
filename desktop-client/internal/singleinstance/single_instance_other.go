//go:build !windows

package singleinstance

import "errors"

var ErrAlreadyRunning = errors.New("voice bridge client is already running")

type Lock struct{}

func Acquire(name string) (*Lock, error) {
	_ = name
	return &Lock{}, nil
}

func (l *Lock) Release() {
	_ = l
}

func ShowAlreadyRunning() {}
