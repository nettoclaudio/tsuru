// Copyright 2018 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import "context"

var (
	_ PlatformStorage = &MockPlatformStorage{}
	_ PlatformService = &MockPlatformService{}
)

// MockPlatformStorage implements PlatformStorage interface
type MockPlatformStorage struct {
	OnInsert      func(Platform) error
	OnFindByName  func(string) (*Platform, error)
	OnFindAll     func() ([]Platform, error)
	OnFindEnabled func() ([]Platform, error)
	OnUpdate      func(Platform) error
	OnDelete      func(Platform) error
}

func (m *MockPlatformStorage) Insert(ctx context.Context, p Platform) error {
	return m.OnInsert(p)
}

func (m *MockPlatformStorage) FindByName(ctx context.Context, name string) (*Platform, error) {
	return m.OnFindByName(name)
}

func (m *MockPlatformStorage) FindAll(ctx context.Context) ([]Platform, error) {
	return m.OnFindAll()
}

func (m *MockPlatformStorage) FindEnabled(ctx context.Context) ([]Platform, error) {
	return m.OnFindEnabled()
}

func (m *MockPlatformStorage) Update(ctx context.Context, p Platform) error {
	return m.OnUpdate(p)
}

func (m *MockPlatformStorage) Delete(ctx context.Context, p Platform) error {
	return m.OnDelete(p)
}

// MockPlatformService implements PlatformService interface
type MockPlatformService struct {
	OnCreate     func(PlatformOptions) error
	OnList       func(bool) ([]Platform, error)
	OnFindByName func(string) (*Platform, error)
	OnUpdate     func(PlatformOptions) error
	OnRemove     func(string) error
	OnRollback   func(PlatformOptions) error
}

func (m *MockPlatformService) Create(ctx context.Context, opts PlatformOptions) error {
	if m.OnCreate == nil {
		return nil
	}
	return m.OnCreate(opts)
}

func (m *MockPlatformService) List(ctx context.Context, enabledOnly bool) ([]Platform, error) {
	if m.OnList == nil {
		return nil, nil
	}
	return m.OnList(enabledOnly)
}

func (m *MockPlatformService) FindByName(ctx context.Context, name string) (*Platform, error) {
	if m.OnFindByName == nil {
		return &Platform{Name: name}, nil
	}
	return m.OnFindByName(name)
}

func (m *MockPlatformService) Update(ctx context.Context, opts PlatformOptions) error {
	if m.OnUpdate == nil {
		return nil
	}
	return m.OnUpdate(opts)
}

func (m *MockPlatformService) Remove(ctx context.Context, name string) error {
	if m.OnRemove == nil {
		return nil
	}
	return m.OnRemove(name)
}

func (m *MockPlatformService) Rollback(ctx context.Context, opts PlatformOptions) error {
	if m.OnRollback == nil {
		return nil
	}
	return m.OnRollback(opts)
}
