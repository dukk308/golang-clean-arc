package swagger_comp

import (
	"strings"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerComponent struct {
	config *SwaggerConfig
	logger logger.Logger
}

func NewSwaggerComponent(logger logger.Logger, config *SwaggerConfig) *SwaggerComponent {
	return &SwaggerComponent{
		config: config,
		logger: logger,
	}
}

func (s *SwaggerComponent) GetConfig() *SwaggerConfig {
	return s.config
}

func (s *SwaggerComponent) RegisterRoutes(router *gin.Engine) {
	if !s.config.Enabled {
		s.logger.Debug("Swagger is disabled, skipping registration")
		return
	}

	docPath := strings.TrimSuffix(s.config.Path, "/*any")
	if docPath == s.config.Path {
		docPath = strings.TrimSuffix(s.config.Path, "/")
	}
	docPath = docPath + "/doc.json"

	url := ginSwagger.URL(docPath)
	router.Any(s.config.Path, ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	s.logger.Infof("Swagger documentation registered at %s", s.config.Path)
}
