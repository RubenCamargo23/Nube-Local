package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    jwtpkg "github.com/org/monitors-platform/pkg/jwt"
    "github.com/org/monitors-platform/pkg/response"
)

const ClaimsKey = "claims"

func JWTMiddleware(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            response.Fail(c, http.StatusUnauthorized, "token requerido")
            c.Abort()
            return
        }

        claims, err := jwtManager.Validate(strings.TrimPrefix(header, "Bearer "))
        if err != nil {
            response.Fail(c, http.StatusUnauthorized, "token inválido")
            c.Abort()
            return
        }

        c.Set(ClaimsKey, claims)
        c.Next()
    }
}

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet(ClaimsKey).(*jwtpkg.Claims)
        for _, required := range roles {
            for _, userRole := range claims.Roles {
                if userRole == required {
                    c.Next()
                    return
                }
            }
        }
        response.Fail(c, http.StatusForbidden, "acceso denegado")
        c.Abort()
    }
}
