package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// Security returns security headers middleware
func Security() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		HSTSMaxAge:                31536000,
		HSTSExcludeSubdomains:     false,
		ContentSecurityPolicy:     "default-src 'self'",
		CSPReportOnly:             false,
		HSTSPreloadEnabled:        false,
		ReferrerPolicy:            "no-referrer",
		PermissionPolicy:          "",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		XPermittedCrossDomain:     "none",
	})
}
