// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/weaveworks/pctl/pkg/catalog"
	"github.com/weaveworks/profiles/api/v1alpha1"
)

type FakeCatalogManager struct {
	InstallStub        func(catalog.InstallConfig) error
	installMutex       sync.RWMutex
	installArgsForCall []struct {
		arg1 catalog.InstallConfig
	}
	installReturns struct {
		result1 error
	}
	installReturnsOnCall map[int]struct {
		result1 error
	}
	ShowStub        func(catalog.CatalogClient, string, string, string) (v1alpha1.ProfileCatalogEntry, error)
	showMutex       sync.RWMutex
	showArgsForCall []struct {
		arg1 catalog.CatalogClient
		arg2 string
		arg3 string
		arg4 string
	}
	showReturns struct {
		result1 v1alpha1.ProfileCatalogEntry
		result2 error
	}
	showReturnsOnCall map[int]struct {
		result1 v1alpha1.ProfileCatalogEntry
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCatalogManager) Install(arg1 catalog.InstallConfig) error {
	fake.installMutex.Lock()
	ret, specificReturn := fake.installReturnsOnCall[len(fake.installArgsForCall)]
	fake.installArgsForCall = append(fake.installArgsForCall, struct {
		arg1 catalog.InstallConfig
	}{arg1})
	stub := fake.InstallStub
	fakeReturns := fake.installReturns
	fake.recordInvocation("Install", []interface{}{arg1})
	fake.installMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCatalogManager) InstallCallCount() int {
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	return len(fake.installArgsForCall)
}

func (fake *FakeCatalogManager) InstallCalls(stub func(catalog.InstallConfig) error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = stub
}

func (fake *FakeCatalogManager) InstallArgsForCall(i int) catalog.InstallConfig {
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	argsForCall := fake.installArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCatalogManager) InstallReturns(result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	fake.installReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCatalogManager) InstallReturnsOnCall(i int, result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	if fake.installReturnsOnCall == nil {
		fake.installReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.installReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCatalogManager) Show(arg1 catalog.CatalogClient, arg2 string, arg3 string, arg4 string) (v1alpha1.ProfileCatalogEntry, error) {
	fake.showMutex.Lock()
	ret, specificReturn := fake.showReturnsOnCall[len(fake.showArgsForCall)]
	fake.showArgsForCall = append(fake.showArgsForCall, struct {
		arg1 catalog.CatalogClient
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.ShowStub
	fakeReturns := fake.showReturns
	fake.recordInvocation("Show", []interface{}{arg1, arg2, arg3, arg4})
	fake.showMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCatalogManager) ShowCallCount() int {
	fake.showMutex.RLock()
	defer fake.showMutex.RUnlock()
	return len(fake.showArgsForCall)
}

func (fake *FakeCatalogManager) ShowCalls(stub func(catalog.CatalogClient, string, string, string) (v1alpha1.ProfileCatalogEntry, error)) {
	fake.showMutex.Lock()
	defer fake.showMutex.Unlock()
	fake.ShowStub = stub
}

func (fake *FakeCatalogManager) ShowArgsForCall(i int) (catalog.CatalogClient, string, string, string) {
	fake.showMutex.RLock()
	defer fake.showMutex.RUnlock()
	argsForCall := fake.showArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeCatalogManager) ShowReturns(result1 v1alpha1.ProfileCatalogEntry, result2 error) {
	fake.showMutex.Lock()
	defer fake.showMutex.Unlock()
	fake.ShowStub = nil
	fake.showReturns = struct {
		result1 v1alpha1.ProfileCatalogEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeCatalogManager) ShowReturnsOnCall(i int, result1 v1alpha1.ProfileCatalogEntry, result2 error) {
	fake.showMutex.Lock()
	defer fake.showMutex.Unlock()
	fake.ShowStub = nil
	if fake.showReturnsOnCall == nil {
		fake.showReturnsOnCall = make(map[int]struct {
			result1 v1alpha1.ProfileCatalogEntry
			result2 error
		})
	}
	fake.showReturnsOnCall[i] = struct {
		result1 v1alpha1.ProfileCatalogEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeCatalogManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	fake.showMutex.RLock()
	defer fake.showMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCatalogManager) recordInvocation(key string, args []interface{}) {
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

var _ catalog.CatalogManager = new(FakeCatalogManager)
