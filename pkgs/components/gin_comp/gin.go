package gin_comp

import "github.com/gin-gonic/gin"

type GinComponent struct {
	router *gin.RouterGroup
	engine *gin.Engine
}

func NewGinComponent() *GinComponent {
	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	router := engine.Group("/api")

	return &GinComponent{
		router: router,
		engine: engine,
	}
}

func (g *GinComponent) Router() *gin.RouterGroup {
	return g.router
}

func (g *GinComponent) Engine() *gin.Engine {
	return g.engine
}
