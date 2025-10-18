package models

// Permission represents a specific permission
type Permission string

const (
	// User permissions
	PermissionCreateUser Permission = "user:create"
	PermissionReadUser   Permission = "user:read"
	PermissionUpdateUser Permission = "user:update"
	PermissionDeleteUser Permission = "user:delete"

	// Product permissions
	PermissionCreateProduct Permission = "product:create"
	PermissionReadProduct   Permission = "product:read"
	PermissionUpdateProduct Permission = "product:update"
	PermissionDeleteProduct Permission = "product:delete"

	// Category permissions
	PermissionCreateCategory Permission = "category:create"
	PermissionReadCategory   Permission = "category:read"
	PermissionUpdateCategory Permission = "category:update"
	PermissionDeleteCategory Permission = "category:delete"

	// Order permissions
	PermissionCreateOrder Permission = "order:create"
	PermissionReadOrder   Permission = "order:read"
	PermissionUpdateOrder Permission = "order:update"
	PermissionDeleteOrder Permission = "order:delete"

	// Cart permissions
	PermissionReadCart   Permission = "cart:read"
	PermissionUpdateCart Permission = "cart:update"

	// Feedback permissions
	PermissionCreateFeedback Permission = "feedback:create"
	PermissionReadFeedback   Permission = "feedback:read"

	// Report permissions
	PermissionReadReport Permission = "report:read"
)

// RolePermissions maps roles to their default permissions
var RolePermissions = map[string][]Permission{
	"admin": {
		// Admin has all permissions
		PermissionCreateUser, PermissionReadUser, PermissionUpdateUser, PermissionDeleteUser,
		PermissionCreateProduct, PermissionReadProduct, PermissionUpdateProduct, PermissionDeleteProduct,
		PermissionCreateCategory, PermissionReadCategory, PermissionUpdateCategory, PermissionDeleteCategory,
		PermissionCreateOrder, PermissionReadOrder, PermissionUpdateOrder, PermissionDeleteOrder,
		PermissionReadCart, PermissionUpdateCart,
		PermissionCreateFeedback, PermissionReadFeedback,
		PermissionReadReport,
	},
	"customer": {
		// Customer has limited permissions
		PermissionReadProduct, PermissionReadCategory,
		PermissionCreateOrder, PermissionReadOrder, PermissionUpdateOrder,
		PermissionReadCart, PermissionUpdateCart,
		PermissionCreateFeedback,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role string, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// GetRolePermissions returns all permissions for a role
func GetRolePermissions(role string) []Permission {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []Permission{}
	}

	return permissions
}
