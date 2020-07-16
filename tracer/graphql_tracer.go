package tracer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/trace"
	"github.com/mrapry/go-lib/golibhelper"
)

// GraphQLTracer struct
type GraphQLTracer struct{}

// TraceQuery method
func (GraphQLTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, varTypes map[string]*introspection.Type) (context.Context, trace.TraceQueryFinishFunc) {
	trace := StartTrace(ctx, fmt.Sprintf("GraphQL-Root:%s", operationName))

	tags := trace.Tags()
	tags["graphql.query"] = queryString
	tags["graphql.operationName"] = operationName
	if len(variables) != 0 {
		tags["graphql.variables"] = variables
	}

	return trace.Context(), func(errs []*gqlerrors.QueryError) {
		defer trace.Finish()
		if len(errs) > 0 {
			msg := errs[0].Error()
			if len(errs) > 1 {
				msg += fmt.Sprintf(" (and %d more errors)", len(errs)-1)
			}
			trace.SetError(errors.New(msg))
		}
	}
}

// TraceField method
func (GraphQLTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, trace.TraceFieldFinishFunc) {
	start := time.Now()
	return ctx, func(err *gqlerrors.QueryError) {
		end := time.Now()
		if !trivial && typeName != "Query" {
			statusColor := golibhelper.Green
			status := " OK  "
			if err != nil {
				statusColor = golibhelper.Red
				status = "ERROR"
			}

			arg, _ := json.Marshal(args)
			fmt.Fprintf(os.Stdout, "%s[GRAPHQL]%s => %s %10s %s | %v | %s %s %s | %13v | %s %s %s | %s\n",
				golibhelper.White, golibhelper.Reset,
				golibhelper.Blue, typeName, golibhelper.Reset,
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, status, golibhelper.Reset,
				end.Sub(start),
				golibhelper.Magenta, label, golibhelper.Reset,
				arg,
			)
		}
	}
}
