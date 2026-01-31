package tracer

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type spanNamePrefixProcessor struct {
	next trace.SpanProcessor
}

func newSpanNamePrefixProcessor(next trace.SpanProcessor) trace.SpanProcessor {
	return &spanNamePrefixProcessor{next: next}
}

func (p *spanNamePrefixProcessor) OnStart(parent context.Context, s trace.ReadWriteSpan) {
	spanName := s.Name()
	prefixedName := p.addPrefixFromName(spanName, s)
	if prefixedName != spanName {
		s.SetName(prefixedName)
	}
	p.next.OnStart(parent, s)
}

func (p *spanNamePrefixProcessor) OnEnd(s trace.ReadOnlySpan) {
	p.next.OnEnd(s)
}

func (p *spanNamePrefixProcessor) Shutdown(ctx context.Context) error {
	return p.next.Shutdown(ctx)
}

func (p *spanNamePrefixProcessor) ForceFlush(ctx context.Context) error {
	return p.next.ForceFlush(ctx)
}

func (p *spanNamePrefixProcessor) addPrefixFromName(spanName string, s trace.ReadWriteSpan) string {
	if strings.HasPrefix(spanName, "GORM_") || strings.HasPrefix(spanName, "GRPC_") || strings.HasPrefix(spanName, "REDIS_") || strings.HasPrefix(spanName, "HTTP_") {
		return spanName
	}

	upperSpanName := strings.ToUpper(spanName)
	redisCommands := []string{"GET", "SET", "HGET", "HSET", "HMSET", "HDEL", "DEL", "EXISTS", "EXPIRE", "TTL", "INCR", "DECR", "LPUSH", "RPUSH", "LPOP", "RPOP", "LLEN", "SADD", "SREM", "SMEMBERS", "SCARD", "SPOP", "PUBLISH", "SUBSCRIBE"}
	for _, cmd := range redisCommands {
		if upperSpanName == cmd {
			return "REDIS_" + spanName
		}
	}
	
	spanKind := s.SpanKind()
	if spanKind == oteltrace.SpanKindClient || spanKind == oteltrace.SpanKindServer {
		if strings.Contains(spanName, "/") {
			if strings.HasPrefix(spanName, "/") {
				return "GRPC_" + spanName
			}
			if strings.Contains(spanName, "pb.") {
				return "GRPC_" + spanName
			}
			if strings.Contains(spanName, "Service/") {
				return "GRPC_" + spanName
			}
		}
	}
	
	return spanName
}
