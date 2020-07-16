package factory

import (
	"github.com/mrapry/go-lib/codebase/factory/dependency"
	"github.com/mrapry/go-lib/codebase/factory/types"
)

// ServiceFactory factory
type ServiceFactory interface {
	GetDependency() dependency.Dependency
	GetModules() []ModuleFactory
	Name() types.Service
}
