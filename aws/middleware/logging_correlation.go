package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/smithy-go/logging"
)

// CorrelatedLogger is a logging.ContextLogger that decorates the entries of
// the wrapped logger with the identifiers available on the context of an
// operation invocation:
//
//	[service=<service id> operation=<operation name> invocation-id=<id> attempt=<n>]
//
// The SDK records these identifiers on the operation context, so the same
// invocation-id ties a logged request to the response returned for it, and the
// attempt number separates the request and response of every retry of that
// invocation. This makes it possible to tell concurrent invocations of the
// same client apart when they log to a shared writer.
//
// Use NewCorrelatedLogger to create one, and configure it as the SDK logger,
// preferably through the config package so that clients constructed while the
// configuration is loaded pick it up too:
//
//	cfg, err := config.LoadDefaultConfig(ctx,
//		config.WithLogger(awsmiddleware.NewCorrelatedLogger(myLogger)))
//
// Only entries whose logger is resolved from the context of the invocation
// that produced them are decorated, i.e. loggers obtained with smithy's
// middleware.GetLogger(ctx) or logging.WithContext(ctx, logger) at the log
// call, which is how the SDK's request/response, retry, and signing logging
// obtains its logger. Entries logged by calling Logf on an unbound logger, such
// as an aws.Config.Logger the application holds and logs to directly, are
// passed through unmodified, as are entries for which no identifier is
// available, e.g. those logged before the request has been built.
//
// Which identifiers are available depends on how far the invocation has
// progressed and on which middleware the stack has: request/response entries
// carry all four, stacks without a retryer carry no attempt, and presigned
// request stacks carry service and operation only. Context values are not
// dropped when an invocation returns, so a logger bound to an invocation's
// context keeps reporting that invocation's identifiers after it returned,
// until the stack values are cleared.
//
// The prefix is placed on the first line of the entry only. Transport dumps
// such as the HTTP request/response logged by the SDK are multi-line, so the
// remaining lines carry no prefix: applications whose log pipeline splits
// entries by line, e.g. a file collector or agent, should use
// CorrelationPrefix with their own structured logging instead.
//
// A nil logger is treated as logging.Nop, and the zero value of this type
// discards entries.
type CorrelatedLogger struct {
	logger logging.Logger
	ctx    context.Context
}

// NewCorrelatedLogger returns a logger that prefixes the entries logged during
// an operation invocation with the invocation's correlation identifiers. A nil
// logger is treated as logging.Nop, and an already correlated logger is
// returned as is, so wrapping twice has no effect. The latter holds only when
// the logger passed in is the correlated one itself: a logger that wraps a
// correlated logger in turn will be decorated a second time.
func NewCorrelatedLogger(logger logging.Logger) *CorrelatedLogger {
	if correlated, ok := logger.(*CorrelatedLogger); ok {
		return correlated
	}
	if logger == nil {
		logger = logging.Nop{}
	}
	return &CorrelatedLogger{logger: logger}
}

// Logf passes the entry to the wrapped logger, prefixed with the correlation
// identifiers of the context the logger was bound to with WithContext. An
// unbound logger, and a bound logger without any identifier on its context,
// logs the entry unmodified.
//
// Implements logging.Logger.
func (l *CorrelatedLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	if l.ctx == nil {
		l.logger.Logf(classification, format, v...)
		return
	}

	// Forward the context to the wrapped logger too, so a logger that
	// implements logging.ContextLogger keeps deriving its own context aware
	// entries behind the correlation prefix.
	logger := logging.WithContext(l.ctx, l.logger)

	prefix := CorrelationPrefix(l.ctx)
	if len(prefix) == 0 {
		logger.Logf(classification, format, v...)
		return
	}

	// The prefix is escaped and prepended to the format string, instead of
	// being passed as the first argument of a new verb, so that an identifier
	// carrying a percent sign cannot introduce a verb of its own, and the
	// caller's arguments, including any that are referred to by index, keep
	// lining up with the verbs of its own format.
	logger.Logf(classification, escapePercent(prefix)+" "+format, v...)
}

// WithContext returns a logger that prefixes entries with the correlation
// identifiers held on the provided context. The returned logger always re-binds
// to ctx, so a logger bound to one invocation never stamps the identifiers of
// that invocation on the entries of a nested one.
//
// Implements logging.ContextLogger.
func (l *CorrelatedLogger) WithContext(ctx context.Context) logging.Logger {
	if l == nil {
		return logging.Nop{}
	}
	return &CorrelatedLogger{logger: l.logger, ctx: ctx}
}

// CorrelationPrefix returns the correlation identifiers held on ctx formatted
// as a log entry prefix, e.g.
//
//	[service=S3 operation=PutObject invocation-id=00000000-0000-4000-8000-000000000000 attempt=1]
//
// Identifiers that are not available on the context are omitted, and an empty
// string is returned when none of them are available. The returned prefix
// carries no trailing space.
//
// Use this to build a custom logging.ContextLogger, e.g. one that adds the
// identifiers as structured fields, or to add them to entries logged outside
// of the SDK's logging, including entries whose log output is split by line,
// where a prefix on the first line of a multi-line entry attributes only that
// line.
func CorrelationPrefix(ctx context.Context) string {
	var b strings.Builder

	if v := GetServiceID(ctx); len(v) != 0 {
		writeCorrelationField(&b, "service", v)
	}
	if v := GetOperationName(ctx); len(v) != 0 {
		writeCorrelationField(&b, "operation", v)
	}
	if v := GetInvocationID(ctx); len(v) != 0 {
		writeCorrelationField(&b, "invocation-id", v)
	}
	if v := GetAttemptNumber(ctx); v != 0 {
		writeCorrelationField(&b, "attempt", strconv.Itoa(v))
	}

	if b.Len() == 0 {
		return ""
	}

	b.WriteByte(']')
	return b.String()
}

func writeCorrelationField(b *strings.Builder, key, value string) {
	if b.Len() == 0 {
		b.WriteByte('[')
	} else {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(value)
}

func escapePercent(v string) string {
	if !strings.ContainsRune(v, '%') {
		return v
	}
	return strings.ReplaceAll(v, "%", "%%")
}
