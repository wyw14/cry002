package domain

import "strings"

func IdempotencyScope(workspaceID, actorID, operation, key string) string {
	parts := []string{workspaceID, actorID, operation, key}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, ":")
}
