// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/weaveworks/pctl/pkg/git"
)

type FakeGit struct {
	AddStub        func() error
	addMutex       sync.RWMutex
	addArgsForCall []struct {
	}
	addReturns struct {
		result1 error
	}
	addReturnsOnCall map[int]struct {
		result1 error
	}
	CheckoutStub        func(string) error
	checkoutMutex       sync.RWMutex
	checkoutArgsForCall []struct {
		arg1 string
	}
	checkoutReturns struct {
		result1 error
	}
	checkoutReturnsOnCall map[int]struct {
		result1 error
	}
	CloneStub        func(string, string, string) error
	cloneMutex       sync.RWMutex
	cloneArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	cloneReturns struct {
		result1 error
	}
	cloneReturnsOnCall map[int]struct {
		result1 error
	}
	CommitStub        func() error
	commitMutex       sync.RWMutex
	commitArgsForCall []struct {
	}
	commitReturns struct {
		result1 error
	}
	commitReturnsOnCall map[int]struct {
		result1 error
	}
	CreateBranchStub        func() error
	createBranchMutex       sync.RWMutex
	createBranchArgsForCall []struct {
	}
	createBranchReturns struct {
		result1 error
	}
	createBranchReturnsOnCall map[int]struct {
		result1 error
	}
	CreateNewBranchStub        func(string) error
	createNewBranchMutex       sync.RWMutex
	createNewBranchArgsForCall []struct {
		arg1 string
	}
	createNewBranchReturns struct {
		result1 error
	}
	createNewBranchReturnsOnCall map[int]struct {
		result1 error
	}
	GetDirectoryStub        func() string
	getDirectoryMutex       sync.RWMutex
	getDirectoryArgsForCall []struct {
	}
	getDirectoryReturns struct {
		result1 string
	}
	getDirectoryReturnsOnCall map[int]struct {
		result1 string
	}
	HasChangesStub        func() (bool, error)
	hasChangesMutex       sync.RWMutex
	hasChangesArgsForCall []struct {
	}
	hasChangesReturns struct {
		result1 bool
		result2 error
	}
	hasChangesReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	InitStub        func() error
	initMutex       sync.RWMutex
	initArgsForCall []struct {
	}
	initReturns struct {
		result1 error
	}
	initReturnsOnCall map[int]struct {
		result1 error
	}
	IsRepositoryStub        func() error
	isRepositoryMutex       sync.RWMutex
	isRepositoryArgsForCall []struct {
	}
	isRepositoryReturns struct {
		result1 error
	}
	isRepositoryReturnsOnCall map[int]struct {
		result1 error
	}
	MergeStub        func(string) (bool, error)
	mergeMutex       sync.RWMutex
	mergeArgsForCall []struct {
		arg1 string
	}
	mergeReturns struct {
		result1 bool
		result2 error
	}
	mergeReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	PushStub        func() error
	pushMutex       sync.RWMutex
	pushArgsForCall []struct {
	}
	pushReturns struct {
		result1 error
	}
	pushReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGit) Add() error {
	fake.addMutex.Lock()
	ret, specificReturn := fake.addReturnsOnCall[len(fake.addArgsForCall)]
	fake.addArgsForCall = append(fake.addArgsForCall, struct {
	}{})
	stub := fake.AddStub
	fakeReturns := fake.addReturns
	fake.recordInvocation("Add", []interface{}{})
	fake.addMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) AddCallCount() int {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return len(fake.addArgsForCall)
}

func (fake *FakeGit) AddCalls(stub func() error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = stub
}

func (fake *FakeGit) AddReturns(result1 error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = nil
	fake.addReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) AddReturnsOnCall(i int, result1 error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = nil
	if fake.addReturnsOnCall == nil {
		fake.addReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Checkout(arg1 string) error {
	fake.checkoutMutex.Lock()
	ret, specificReturn := fake.checkoutReturnsOnCall[len(fake.checkoutArgsForCall)]
	fake.checkoutArgsForCall = append(fake.checkoutArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.CheckoutStub
	fakeReturns := fake.checkoutReturns
	fake.recordInvocation("Checkout", []interface{}{arg1})
	fake.checkoutMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) CheckoutCallCount() int {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	return len(fake.checkoutArgsForCall)
}

func (fake *FakeGit) CheckoutCalls(stub func(string) error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = stub
}

func (fake *FakeGit) CheckoutArgsForCall(i int) string {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	argsForCall := fake.checkoutArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) CheckoutReturns(result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	fake.checkoutReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CheckoutReturnsOnCall(i int, result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	if fake.checkoutReturnsOnCall == nil {
		fake.checkoutReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkoutReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Clone(arg1 string, arg2 string, arg3 string) error {
	fake.cloneMutex.Lock()
	ret, specificReturn := fake.cloneReturnsOnCall[len(fake.cloneArgsForCall)]
	fake.cloneArgsForCall = append(fake.cloneArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.CloneStub
	fakeReturns := fake.cloneReturns
	fake.recordInvocation("Clone", []interface{}{arg1, arg2, arg3})
	fake.cloneMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) CloneCallCount() int {
	fake.cloneMutex.RLock()
	defer fake.cloneMutex.RUnlock()
	return len(fake.cloneArgsForCall)
}

func (fake *FakeGit) CloneCalls(stub func(string, string, string) error) {
	fake.cloneMutex.Lock()
	defer fake.cloneMutex.Unlock()
	fake.CloneStub = stub
}

func (fake *FakeGit) CloneArgsForCall(i int) (string, string, string) {
	fake.cloneMutex.RLock()
	defer fake.cloneMutex.RUnlock()
	argsForCall := fake.cloneArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeGit) CloneReturns(result1 error) {
	fake.cloneMutex.Lock()
	defer fake.cloneMutex.Unlock()
	fake.CloneStub = nil
	fake.cloneReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CloneReturnsOnCall(i int, result1 error) {
	fake.cloneMutex.Lock()
	defer fake.cloneMutex.Unlock()
	fake.CloneStub = nil
	if fake.cloneReturnsOnCall == nil {
		fake.cloneReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.cloneReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Commit() error {
	fake.commitMutex.Lock()
	ret, specificReturn := fake.commitReturnsOnCall[len(fake.commitArgsForCall)]
	fake.commitArgsForCall = append(fake.commitArgsForCall, struct {
	}{})
	stub := fake.CommitStub
	fakeReturns := fake.commitReturns
	fake.recordInvocation("Commit", []interface{}{})
	fake.commitMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) CommitCallCount() int {
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	return len(fake.commitArgsForCall)
}

func (fake *FakeGit) CommitCalls(stub func() error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
	fake.CommitStub = stub
}

func (fake *FakeGit) CommitReturns(result1 error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
	fake.CommitStub = nil
	fake.commitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CommitReturnsOnCall(i int, result1 error) {
	fake.commitMutex.Lock()
	defer fake.commitMutex.Unlock()
	fake.CommitStub = nil
	if fake.commitReturnsOnCall == nil {
		fake.commitReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CreateBranch() error {
	fake.createBranchMutex.Lock()
	ret, specificReturn := fake.createBranchReturnsOnCall[len(fake.createBranchArgsForCall)]
	fake.createBranchArgsForCall = append(fake.createBranchArgsForCall, struct {
	}{})
	stub := fake.CreateBranchStub
	fakeReturns := fake.createBranchReturns
	fake.recordInvocation("CreateBranch", []interface{}{})
	fake.createBranchMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) CreateBranchCallCount() int {
	fake.createBranchMutex.RLock()
	defer fake.createBranchMutex.RUnlock()
	return len(fake.createBranchArgsForCall)
}

func (fake *FakeGit) CreateBranchCalls(stub func() error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = stub
}

func (fake *FakeGit) CreateBranchReturns(result1 error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = nil
	fake.createBranchReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CreateBranchReturnsOnCall(i int, result1 error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = nil
	if fake.createBranchReturnsOnCall == nil {
		fake.createBranchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createBranchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CreateNewBranch(arg1 string) error {
	fake.createNewBranchMutex.Lock()
	ret, specificReturn := fake.createNewBranchReturnsOnCall[len(fake.createNewBranchArgsForCall)]
	fake.createNewBranchArgsForCall = append(fake.createNewBranchArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.CreateNewBranchStub
	fakeReturns := fake.createNewBranchReturns
	fake.recordInvocation("CreateNewBranch", []interface{}{arg1})
	fake.createNewBranchMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) CreateNewBranchCallCount() int {
	fake.createNewBranchMutex.RLock()
	defer fake.createNewBranchMutex.RUnlock()
	return len(fake.createNewBranchArgsForCall)
}

func (fake *FakeGit) CreateNewBranchCalls(stub func(string) error) {
	fake.createNewBranchMutex.Lock()
	defer fake.createNewBranchMutex.Unlock()
	fake.CreateNewBranchStub = stub
}

func (fake *FakeGit) CreateNewBranchArgsForCall(i int) string {
	fake.createNewBranchMutex.RLock()
	defer fake.createNewBranchMutex.RUnlock()
	argsForCall := fake.createNewBranchArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) CreateNewBranchReturns(result1 error) {
	fake.createNewBranchMutex.Lock()
	defer fake.createNewBranchMutex.Unlock()
	fake.CreateNewBranchStub = nil
	fake.createNewBranchReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CreateNewBranchReturnsOnCall(i int, result1 error) {
	fake.createNewBranchMutex.Lock()
	defer fake.createNewBranchMutex.Unlock()
	fake.CreateNewBranchStub = nil
	if fake.createNewBranchReturnsOnCall == nil {
		fake.createNewBranchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createNewBranchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) GetDirectory() string {
	fake.getDirectoryMutex.Lock()
	ret, specificReturn := fake.getDirectoryReturnsOnCall[len(fake.getDirectoryArgsForCall)]
	fake.getDirectoryArgsForCall = append(fake.getDirectoryArgsForCall, struct {
	}{})
	stub := fake.GetDirectoryStub
	fakeReturns := fake.getDirectoryReturns
	fake.recordInvocation("GetDirectory", []interface{}{})
	fake.getDirectoryMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) GetDirectoryCallCount() int {
	fake.getDirectoryMutex.RLock()
	defer fake.getDirectoryMutex.RUnlock()
	return len(fake.getDirectoryArgsForCall)
}

func (fake *FakeGit) GetDirectoryCalls(stub func() string) {
	fake.getDirectoryMutex.Lock()
	defer fake.getDirectoryMutex.Unlock()
	fake.GetDirectoryStub = stub
}

func (fake *FakeGit) GetDirectoryReturns(result1 string) {
	fake.getDirectoryMutex.Lock()
	defer fake.getDirectoryMutex.Unlock()
	fake.GetDirectoryStub = nil
	fake.getDirectoryReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeGit) GetDirectoryReturnsOnCall(i int, result1 string) {
	fake.getDirectoryMutex.Lock()
	defer fake.getDirectoryMutex.Unlock()
	fake.GetDirectoryStub = nil
	if fake.getDirectoryReturnsOnCall == nil {
		fake.getDirectoryReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getDirectoryReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeGit) HasChanges() (bool, error) {
	fake.hasChangesMutex.Lock()
	ret, specificReturn := fake.hasChangesReturnsOnCall[len(fake.hasChangesArgsForCall)]
	fake.hasChangesArgsForCall = append(fake.hasChangesArgsForCall, struct {
	}{})
	stub := fake.HasChangesStub
	fakeReturns := fake.hasChangesReturns
	fake.recordInvocation("HasChanges", []interface{}{})
	fake.hasChangesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGit) HasChangesCallCount() int {
	fake.hasChangesMutex.RLock()
	defer fake.hasChangesMutex.RUnlock()
	return len(fake.hasChangesArgsForCall)
}

func (fake *FakeGit) HasChangesCalls(stub func() (bool, error)) {
	fake.hasChangesMutex.Lock()
	defer fake.hasChangesMutex.Unlock()
	fake.HasChangesStub = stub
}

func (fake *FakeGit) HasChangesReturns(result1 bool, result2 error) {
	fake.hasChangesMutex.Lock()
	defer fake.hasChangesMutex.Unlock()
	fake.HasChangesStub = nil
	fake.hasChangesReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) HasChangesReturnsOnCall(i int, result1 bool, result2 error) {
	fake.hasChangesMutex.Lock()
	defer fake.hasChangesMutex.Unlock()
	fake.HasChangesStub = nil
	if fake.hasChangesReturnsOnCall == nil {
		fake.hasChangesReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.hasChangesReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) Init() error {
	fake.initMutex.Lock()
	ret, specificReturn := fake.initReturnsOnCall[len(fake.initArgsForCall)]
	fake.initArgsForCall = append(fake.initArgsForCall, struct {
	}{})
	stub := fake.InitStub
	fakeReturns := fake.initReturns
	fake.recordInvocation("Init", []interface{}{})
	fake.initMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) InitCallCount() int {
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	return len(fake.initArgsForCall)
}

func (fake *FakeGit) InitCalls(stub func() error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = stub
}

func (fake *FakeGit) InitReturns(result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	fake.initReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) InitReturnsOnCall(i int, result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	if fake.initReturnsOnCall == nil {
		fake.initReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) IsRepository() error {
	fake.isRepositoryMutex.Lock()
	ret, specificReturn := fake.isRepositoryReturnsOnCall[len(fake.isRepositoryArgsForCall)]
	fake.isRepositoryArgsForCall = append(fake.isRepositoryArgsForCall, struct {
	}{})
	stub := fake.IsRepositoryStub
	fakeReturns := fake.isRepositoryReturns
	fake.recordInvocation("IsRepository", []interface{}{})
	fake.isRepositoryMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) IsRepositoryCallCount() int {
	fake.isRepositoryMutex.RLock()
	defer fake.isRepositoryMutex.RUnlock()
	return len(fake.isRepositoryArgsForCall)
}

func (fake *FakeGit) IsRepositoryCalls(stub func() error) {
	fake.isRepositoryMutex.Lock()
	defer fake.isRepositoryMutex.Unlock()
	fake.IsRepositoryStub = stub
}

func (fake *FakeGit) IsRepositoryReturns(result1 error) {
	fake.isRepositoryMutex.Lock()
	defer fake.isRepositoryMutex.Unlock()
	fake.IsRepositoryStub = nil
	fake.isRepositoryReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) IsRepositoryReturnsOnCall(i int, result1 error) {
	fake.isRepositoryMutex.Lock()
	defer fake.isRepositoryMutex.Unlock()
	fake.IsRepositoryStub = nil
	if fake.isRepositoryReturnsOnCall == nil {
		fake.isRepositoryReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.isRepositoryReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Merge(arg1 string) (bool, error) {
	fake.mergeMutex.Lock()
	ret, specificReturn := fake.mergeReturnsOnCall[len(fake.mergeArgsForCall)]
	fake.mergeArgsForCall = append(fake.mergeArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.MergeStub
	fakeReturns := fake.mergeReturns
	fake.recordInvocation("Merge", []interface{}{arg1})
	fake.mergeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGit) MergeCallCount() int {
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	return len(fake.mergeArgsForCall)
}

func (fake *FakeGit) MergeCalls(stub func(string) (bool, error)) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = stub
}

func (fake *FakeGit) MergeArgsForCall(i int) string {
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	argsForCall := fake.mergeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) MergeReturns(result1 bool, result2 error) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = nil
	fake.mergeReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) MergeReturnsOnCall(i int, result1 bool, result2 error) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = nil
	if fake.mergeReturnsOnCall == nil {
		fake.mergeReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.mergeReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) Push() error {
	fake.pushMutex.Lock()
	ret, specificReturn := fake.pushReturnsOnCall[len(fake.pushArgsForCall)]
	fake.pushArgsForCall = append(fake.pushArgsForCall, struct {
	}{})
	stub := fake.PushStub
	fakeReturns := fake.pushReturns
	fake.recordInvocation("Push", []interface{}{})
	fake.pushMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGit) PushCallCount() int {
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	return len(fake.pushArgsForCall)
}

func (fake *FakeGit) PushCalls(stub func() error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = stub
}

func (fake *FakeGit) PushReturns(result1 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	fake.pushReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) PushReturnsOnCall(i int, result1 error) {
	fake.pushMutex.Lock()
	defer fake.pushMutex.Unlock()
	fake.PushStub = nil
	if fake.pushReturnsOnCall == nil {
		fake.pushReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pushReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	fake.cloneMutex.RLock()
	defer fake.cloneMutex.RUnlock()
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	fake.createBranchMutex.RLock()
	defer fake.createBranchMutex.RUnlock()
	fake.createNewBranchMutex.RLock()
	defer fake.createNewBranchMutex.RUnlock()
	fake.getDirectoryMutex.RLock()
	defer fake.getDirectoryMutex.RUnlock()
	fake.hasChangesMutex.RLock()
	defer fake.hasChangesMutex.RUnlock()
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	fake.isRepositoryMutex.RLock()
	defer fake.isRepositoryMutex.RUnlock()
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	fake.pushMutex.RLock()
	defer fake.pushMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGit) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ git.Git = new(FakeGit)
