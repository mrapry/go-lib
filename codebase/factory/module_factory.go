package factory

import (
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/codebase/interfaces"
)

// ModuleFactory factory
type ModuleFactory interface {
	RestHandler() interfaces.EchoRestHandler
	GRPCHandler() interfaces.GRPCHandler
	GraphQLHandler() interfaces.GraphQLHandler
	WorkerHandler(workerType types.Worker) interfaces.WorkerHandler
	Name() types.Module
}
