package principal

import "context"

var SystemPrincipal = Principal{
	Identity:    "system",
	Roles:       []string{"admin"},
	Permissions: []string{"*"},
}

type Principal struct {
	Identity    string
	Roles       []string
	Permissions []string
}

func New(identity string, roles []string, permissions []string) Principal {
	return Principal{
		Identity:    identity,
		Roles:       roles,
		Permissions: permissions,
	}
}

type ctxKey struct{}

func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

func WithSystemAsPrincipal(ctx context.Context) context.Context {
	return WithPrincipal(ctx, SystemPrincipal)
}

func GetPrincipal(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKey{}).(Principal)
	return p, ok
}

func MustGetPrincipal(ctx context.Context) Principal {
	p, ok := GetPrincipal(ctx)
	if !ok {
		panic("context: principal not found")
	}
	return p
}
