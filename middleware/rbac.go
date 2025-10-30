package middleware

import (
	"net/http"

	"health-store/models"

	"github.com/gin-gonic/gin"
)

// RequirePermission creates middleware that requires specific permissions
func RequirePermission(permission models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(string)
		if !models.HasPermission(role, permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Insufficient permissions",
				"required": string(permission),
				"role":     role,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission creates middleware that requires any of the specified permissions
func RequireAnyPermission(permissions ...models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(string)
		hasPermission := false

		for _, permission := range permissions {
			if models.HasPermission(role, permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Insufficient permissions",
				"required": "any of " + permissionsString(permissions),
				"role":     role,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that requires specific roles (legacy support)
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		currentRole := userRole.(string)
		authorized := false

		for _, role := range roles {
			if currentRole == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{
				"error":        "Insufficient role permissions",
				"required":     "any of " + rolesString(roles),
				"current_role": currentRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper function to convert slice to string
func permissionsString(permissions []models.Permission) string {
	result := ""
	for i, p := range permissions {
		if i > 0 {
			result += ", "
		}
		result += string(p)
	}
	return result
}

// Helper function to convert string slice to string
func rolesString(roles []string) string {
	result := ""
	for i, r := range roles {
		if i > 0 {
			result += ", "
		}
		result += r
	}
	return result
}
