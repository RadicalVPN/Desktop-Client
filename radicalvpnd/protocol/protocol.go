package protocol

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"radicalvpnd/logger"
	service "radicalvpnd/services"
	"radicalvpnd/settings"
	"radicalvpnd/webapi"
	"radicalvpnd/wireguard"

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

func (p *Protocol) Start(startedPortChannel chan<- string) {
	log.Info("Starting Daemon Protocol..")

	p.LoadMiddlewaares()
	p.LoadRoutes()

	var listeningPort string
	envSecret, envSecretPresent := os.LookupEnv("RADICALVPND_PORT")
	if envSecretPresent {
		listeningPort = ":" + envSecret
	} else {
		listeningPort = ":0"
	}

	listener, _ := net.Listen("tcp", listeningPort)
	_, port, _ := net.SplitHostPort(listener.Addr().String())

	startedPortChannel <- port

	log.Info("Daemon Protocol listening on 127.0.0.1:", port)
	http.Serve(listener, p.engine)
}

func (p *Protocol) LoadRoutes() {
	r := p.engine

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/local/connected", func(c *gin.Context) {
		wg := wireguard.NewWireguard()

		c.JSON(http.StatusOK, gin.H{
			"connected": wg.IsConnected(),
		})
	})

	r.POST("/local/connect", func(c *gin.Context) {
		wg := wireguard.NewWireguard()

		err := wg.Connect()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.Status(http.StatusOK)
	})

	r.POST("/local/disconnect", func(c *gin.Context) {
		wg := wireguard.NewWireguard()

		wg.Disconnect()

		c.Status(http.StatusOK)
	})

	r.GET("/server", func(c *gin.Context) {
		req, err := http.NewRequest("GET", "https://radicalvpn.com/api/1.0/server", nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		req.Header.Set("Cookie", settings.GetSessionCookie())

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		servers := []webapi.Server{}

		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		servers = service.PingServers(servers)

		c.JSON(http.StatusOK, servers)
	})

	r.GET("/me", func(c *gin.Context) {

	})

	r.POST("/login", func(c *gin.Context) {
		internalBody := webapi.Signin{}
		internalBody.RememberMe = true

		if err := c.ShouldBindJSON(&internalBody); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		body, _ := json.Marshal(internalBody)
		requestBody := bytes.NewBuffer(body)

		resp, err := http.Post("https://radicalvpn.com/api/1.0/auth", "application/json", requestBody)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			sessionCookie := resp.Cookies()[0].Value

			sett := settings.NewSettings()
			sett.SetSession(sessionCookie)

			c.Status(http.StatusOK)
		}
	})
}

func (p *Protocol) LoadMiddlewaares() {
	r := p.engine

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(p.AuthMiddleware())
}
