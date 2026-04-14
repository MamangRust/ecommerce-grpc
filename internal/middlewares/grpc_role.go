package middlewares

import (
	"context"
	"ecommerce/internal/service"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoleInterceptor struct {
	roleService service.RoleService
	rolesMap    map[string][]string // method -> allowed roles
}

func NewRoleInterceptor(roleService service.RoleService, rolesMap map[string][]string) *RoleInterceptor {
	return &RoleInterceptor{
		roleService: roleService,
		rolesMap:    rolesMap,
	}
}

func (interceptor *RoleInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		allowedRoles, ok := interceptor.rolesMap[info.FullMethod]
		if !ok || len(allowedRoles) == 0 {
			// If no roles specified, allowed for all authenticated users
			return handler(ctx, req)
		}

		userID, ok := ctx.Value("userID").(string)
		if !ok || userID == "" {
			return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
		}

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
		}

		roles, err := interceptor.roleService.FindByUserId(ctx, userIDInt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch user roles: %v", err)
		}

		for _, userRole := range roles {
			for _, allowedRole := range allowedRoles {
				if userRole.RoleName == allowedRole {
					return handler(ctx, req)
				}
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "you do not have permission to access this resource")
	}
}
