package swagger_comp

import (
	"flag"
)

var (
	swaggerEnabledVal bool
	swaggerPathVal    string
	swaggerBasePathVal string
	swaggerTitleVal   string
	swaggerVersionVal string
	swaggerHostVal    string
)

var (
	SwaggerEnabled  = &swaggerEnabledVal
	SwaggerPath     = &swaggerPathVal
	SwaggerBasePath = &swaggerBasePathVal
	SwaggerTitle    = &swaggerTitleVal
	SwaggerVersion  = &swaggerVersionVal
	SwaggerHost     = &swaggerHostVal
)

func init() {
	if flag.Lookup("swagger-enabled") == nil {
		flag.BoolVar(&swaggerEnabledVal, "swagger-enabled", true, "Enable swagger documentation")
	}
	if flag.Lookup("swagger-path") == nil {
		flag.StringVar(&swaggerPathVal, "swagger-path", "/swagger/*any", "Swagger documentation path")
	}
	if flag.Lookup("swagger-base-path") == nil {
		flag.StringVar(&swaggerBasePathVal, "swagger-base-path", "/", "Swagger base path")
	}
	if flag.Lookup("swagger-title") == nil {
		flag.StringVar(&swaggerTitleVal, "swagger-title", "API Documentation", "Swagger API title")
	}
	if flag.Lookup("swagger-version") == nil {
		flag.StringVar(&swaggerVersionVal, "swagger-version", "1.0", "Swagger API version")
	}
	if flag.Lookup("swagger-host") == nil {
		flag.StringVar(&swaggerHostVal, "swagger-host", "", "Swagger host")
	}
}

func LoadSwaggerConfig() *SwaggerConfig {
	return &SwaggerConfig{
		Enabled:  *SwaggerEnabled,
		Path:     *SwaggerPath,
		BasePath: *SwaggerBasePath,
		Title:    *SwaggerTitle,
		Version:  *SwaggerVersion,
		Host:     *SwaggerHost,
	}
}
