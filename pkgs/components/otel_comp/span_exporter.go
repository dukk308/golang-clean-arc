package otel_comp

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type spanNamePrefixExporter struct {
	exporter trace.SpanExporter
}

func newSpanNamePrefixExporter(exporter trace.SpanExporter) trace.SpanExporter {
	return &spanNamePrefixExporter{exporter: exporter}
}

func (e *spanNamePrefixExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	modifiedSpans := make([]trace.ReadOnlySpan, len(spans))
	for i, span := range spans {
		modifiedSpans[i] = &spanWithPrefix{ReadOnlySpan: span}
	}
	return e.exporter.ExportSpans(ctx, modifiedSpans)
}

func (e *spanNamePrefixExporter) Shutdown(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}

type spanWithPrefix struct {
	trace.ReadOnlySpan
}

func (s *spanWithPrefix) Name() string {
	spanName := s.ReadOnlySpan.Name()
	return addPrefixToSpanName(spanName, s.ReadOnlySpan)
}

func addPrefixToSpanName(spanName string, s trace.ReadOnlySpan) string {
	if strings.HasPrefix(spanName, "GORM_") || strings.HasPrefix(spanName, "GRPC_") || strings.HasPrefix(spanName, "REDIS_") || strings.HasPrefix(spanName, "HTTP_") {
		return spanName
	}

	attrs := s.Attributes()
	
	var dbSystem string
	hasDbAttributes := false
	
	for _, attr := range attrs {
		key := string(attr.Key)
		value := attr.Value.AsString()
		
		if key == "db.system" {
			dbSystem = value
			hasDbAttributes = true
		}
		
		if key == "redis.command" {
			return "REDIS_" + spanName
		}
		
		if key == "db.statement" || key == "db.operation" || key == "db.sql.table" {
			hasDbAttributes = true
		}
	}
	
	if hasDbAttributes {
		if dbSystem == "redis" {
			return "REDIS_" + spanName
		}
		return "GORM_" + spanName
	}
	
	upperSpanName := strings.ToUpper(spanName)
	redisCommands := []string{"GET", "SET", "HGET", "HSET", "HMSET", "HDEL", "DEL", "EXISTS", "EXPIRE", "TTL", "INCR", "DECR", "LPUSH", "RPUSH", "LPOP", "RPOP", "LLEN", "SADD", "SREM", "SMEMBERS", "SCARD", "SPOP", "PUBLISH", "SUBSCRIBE"}
	for _, cmd := range redisCommands {
		if upperSpanName == cmd {
			return "REDIS_" + spanName
		}
	}
	
	spanKind := s.SpanKind()
	if spanKind == oteltrace.SpanKindClient {
		if strings.Contains(strings.ToLower(spanName), "select") || 
			strings.Contains(strings.ToLower(spanName), "insert") || 
			strings.Contains(strings.ToLower(spanName), "update") || 
			strings.Contains(strings.ToLower(spanName), "delete") ||
			strings.Contains(strings.ToLower(spanName), "create") ||
			strings.Contains(strings.ToLower(spanName), "drop") ||
			strings.Contains(strings.ToLower(spanName), "alter") {
			return "GORM_" + spanName
		}
	}
	
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
