package protocol

import (
	"net/http"
	"radicalvpnd/logger"

	"github.com/gin-gonic/gin"
)

type Protocol struct {
	secret string
	engine *gin.Engine
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("prtctl")
}

func NewProtocol(secret string) *Protocol {
	gin.SetMode((gin.ReleaseMode))
	r := gin.Default()

	return &Protocol{secret: secret, engine: r}
}

func (p *Protocol) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.Request.Header.Get("X-RadicalDaemon-Secret")

		if secret != p.secret {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}

func (p *Protocol) Start() {
	log.Info("Starting Daemon Protocol..")

	p.LoadMiddlewaares()
	p.LoadRoutes()

	addr := "127.0.0.1:8095"
	log.Info("Daemon Protocol listening on ", addr)
	p.engine.Run(addr)
}

func (p *Protocol) LoadRoutes() {
	r := p.engine

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func (p *Protocol) LoadMiddlewaares() {
	r := p.engine

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(p.AuthMiddleware())
}
