// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package envstore provides a concise interface for handling app's environment
// variables over backend storages. In addiction, it provides some
// implementations for that interface.
package envstore

import "github.com/tsuru/tsuru/app/bind"

// EnvStorer defines the commom way to handle the app's environment variables.
type EnvStorer interface {
	Get(...string) (map[string]bind.EnvVar, error)
	Set(...bind.EnvVar) error
	Unset(...string) error
}

// MongoDB is a struct that implements a EnvStorer interface.
type MongoDB struct {
}

// Get attempts to fetch.
func (m *MongoDB) Get(names ...string) (map[string]bind.EnvVar, error) {
	return nil, nil
}

// Set records a bunch of environment variables into MongoDB.
func (m *MongoDB) Set(envs ...bind.EnvVar) error {
	return nil
}

// Unset removes a envimenroment variables. Returns nil when the unset
// was successfully, otherwise returns a error with ca.
func (m *MongoDB) Unset(names ...string) error {
	return nil
}
