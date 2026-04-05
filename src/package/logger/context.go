package logger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
)

type contextKey string // private type used to avoid key collisions in context

const (
	ctxKeyUserID   contextKey = "user_id"
	ctxKeyClientIP contextKey = "client_ip"
	ctxKeyToken    contextKey = "token_hint"
	ctxKeyTokenSHA contextKey = "token_hash"
)

func ContextWithUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func ContextWithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ctxKeyClientIP, ip)
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	hint := token
	if len(token) > 10 {
		hint = "..." + token[len(token)-10:]
	}
	ctx = context.WithValue(ctx, ctxKeyToken, hint)
	if token != "" {
		sum := sha256.Sum256([]byte(token))
		ctx = context.WithValue(ctx, ctxKeyTokenSHA, hex.EncodeToString(sum[:]))
	}
	return ctx
}

func userIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(ctxKeyUserID).(uint64)
	return v, ok
}

func clientIPFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyClientIP).(string); ok {
		return v
	}
	return ""
}

func tokenHintFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyToken).(string); ok {
		return v
	}
	return ""
}

func tokenHashFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTokenSHA).(string); ok {
		return v
	}
	return ""
}

// buildArgs constructs structured log arguments based on context values (userID, clientIP, tokenHint)
func buildArgs(ctx context.Context) []any {
	var args []any
	if userID, ok := userIDFromContext(ctx); ok {
		args = append(args, slog.Uint64("user_id", userID))
	}
	if ip := clientIPFromContext(ctx); ip != "" {
		args = append(args, slog.String("client_ip", ip))
	}
	if hint := tokenHintFromContext(ctx); hint != "" {
		args = append(args, slog.String("token_hint", hint))
	}
	if hash := tokenHashFromContext(ctx); hash != "" {
		args = append(args, slog.String("token_hash", hash))
	}
	return args
}

// ── Context-aware log functions ──

func DebugfWithCtx(ctx context.Context, format string, args ...any) {
	log().Debug(fmt.Sprintf(format, args...), buildArgs(ctx)...)
}

func InfofWithCtx(ctx context.Context, format string, args ...any) {
	log().Info(fmt.Sprintf(format, args...), buildArgs(ctx)...)
}

func WarnfWithCtx(ctx context.Context, format string, args ...any) {
	log().Warn(fmt.Sprintf(format, args...), buildArgs(ctx)...)
}

func ErrorfWithCtx(ctx context.Context, format string, args ...any) {
	log().Error(fmt.Sprintf(format, args...), buildArgs(ctx)...)
}
