package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"radicalvpnd/logger"
	service "radicalvpnd/services"
	"radicalvpnd/settings"
	"radicalvpnd/version"
	"radicalvpnd/webapi"
	"radicalvpnd/wireguard"
	"time"

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
	r := gin.New()

	return &Protocol{secret: secret, engine: r}
}

func (p *Protocol) ensureAuth() bool {
	sett := settings.NewSettings()

	if sett.Session.Secret == "" {
		return false
	}

	return true
}

func (p *Protocol) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.Request.Header.Get("x-radical-daemon-secret")

		if secret != p.secret && c.Request.Method != "OPTIONS" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}

func (p *Protocol) CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, X-Requested-With, X-Radical-Daemon-Secret")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}

func (p *Protocol) LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// process the request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		log.Info(fmt.Sprintf("REQUEST -> %s %s %s %d", method, path, latency, status))
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

func (p *Protocol) getHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			MaxIdleConnsPerHost:   -1,
		},
	}
}

func (p *Protocol) proxyGetRequest(url string, c *gin.Context) {
	if !p.ensureAuth() {
		log.Info("401")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Cookie", settings.GetSessionCookie())

	resp, err := p.getHttpClient().Do(req)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	//parse body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		log.Info("not ok")
		c.AbortWithStatus(resp.StatusCode)
		c.Data(http.StatusOK, contentType, body)
	} else {
		c.Data(http.StatusOK, contentType, body)
	}
}

func (p *Protocol) LoadRoutes() {
	r := p.engine

	r.OPTIONS("/local/connect", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/version", func(c *gin.Context) {
		var outdated bool

		if version.IsRelease() {
			outdated = version.IsReleaseOutdated()
		} else {
			outdated = version.IsNightlyOutdated()
		}

		c.JSON(http.StatusOK, gin.H{
			"currentVersion": version.GetVersion(),
			"nightly":        version.IsNightly(),
			"release":        version.IsRelease(),
			"outdated":       outdated,
		})
	})

	r.GET("/local/connected", func(c *gin.Context) {
		if !p.ensureAuth() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		wg := wireguard.NewWireguard()

		c.JSON(http.StatusOK, gin.H{
			"connected": wg.IsConnected(),
		})
	})

	r.POST("/local/connect", func(c *gin.Context) {
		if !p.ensureAuth() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		internalBody := webapi.VpnConnect{}
		if err := c.ShouldBindJSON(&internalBody); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		wg := wireguard.NewWireguard()
		err := wg.Connect(internalBody.NodeLocation, internalBody.PrivacyFirewall)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.Status(http.StatusOK)
	})

	r.POST("/local/disconnect", func(c *gin.Context) {
		if !p.ensureAuth() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		wg := wireguard.NewWireguard()

		wg.Disconnect()

		c.Status(http.StatusOK)
	})

	r.GET("/server", func(c *gin.Context) {
		if !p.ensureAuth() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		req, err := http.NewRequest("GET", "https://radicalvpn.com/api/1.0/server", nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		req.Header.Set("Cookie", settings.GetSessionCookie())

		resp, err := p.getHttpClient().Do(req)

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

	r.GET("/", func(c *gin.Context) {
		p.proxyGetRequest("https://radicalvpn.com/api/1.0/auth", c)
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
			//parse body
			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusUnauthorized, "application/json; charset=utf-8", body)
		} else {
			sessionCookie := resp.Cookies()[0].Value

			sett := settings.NewSettings()
			sett.SetSession(sessionCookie)

			c.Status(http.StatusOK)
		}
	})

	r.GET("/privacy_firewall", func(c *gin.Context) {
		p.proxyGetRequest("https://radicalvpn.com/api/1.0/internal/privacy_firewall", c)
	})
}

func (p *Protocol) LoadMiddlewaares() {
	r := p.engine

	r.Use(gin.Recovery())
	r.Use(p.LogMiddleware())
	r.Use(p.AuthMiddleware())
	r.Use(p.CorsMiddleware())
}
