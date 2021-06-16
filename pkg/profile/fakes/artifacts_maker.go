// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/weaveworks/pctl/pkg/profile"
	"github.com/weaveworks/profiles/api/v1alpha1"
)

type FakeArtifactsMaker struct {
	MakeStub        func(v1alpha1.ProfileInstallation) error
	makeMutex       sync.RWMutex
	makeArgsForCall []struct {
		arg1 v1alpha1.ProfileInstallation
	}
	makeReturns struct {
		result1 error
	}
	makeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifactsMaker) Make(arg1 v1alpha1.ProfileInstallation) error {
	fake.makeMutex.Lock()
	ret, specificReturn := fake.makeReturnsOnCall[len(fake.makeArgsForCall)]
	fake.makeArgsForCall = append(fake.makeArgsForCall, struct {
		arg1 v1alpha1.ProfileInstallation
	}{arg1})
	stub := fake.MakeStub
	fakeReturns := fake.makeReturns
	fake.recordInvocation("Make", []interface{}{arg1})
	fake.makeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeArtifactsMaker) MakeCallCount() int {
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	return len(fake.makeArgsForCall)
}

func (fake *FakeArtifactsMaker) MakeCalls(stub func(v1alpha1.ProfileInstallation) error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = stub
}

func (fake *FakeArtifactsMaker) MakeArgsForCall(i int) v1alpha1.ProfileInstallation {
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	argsForCall := fake.makeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeArtifactsMaker) MakeReturns(result1 error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = nil
	fake.makeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactsMaker) MakeReturnsOnCall(i int, result1 error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = nil
	if fake.makeReturnsOnCall == nil {
		fake.makeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.makeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactsMaker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeArtifactsMaker) recordInvocation(key string, args []interface{}) {
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

var _ profile.ArtifactsMaker = new(FakeArtifactsMaker)