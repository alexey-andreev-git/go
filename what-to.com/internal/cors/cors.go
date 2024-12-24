package cors

import (
	"strconv"
	"strings"
)

type (
	Cors struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		AllowCredentials bool
		MaxAge           int
	}
)

const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlMaxAge           = "Access-Control-Max-Age"
)

// NewCors returns a new Cors with empty values
func NewCors() *Cors {
	return &Cors{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{},
		AllowedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           600,
	}
}

// NewDefaultCors returns a new Cors with default values
// AllowedOrigins:   []string{"*"}
// AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
// AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
// AllowCredentials: true
// MaxAge:           600
func NewDefaultCors() *Cors {
	return &Cors{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		MaxAge:           600,
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (c *Cors) AddAllowedOrigins(origins ...string) {
	for _, origin := range origins {
		if !contains(c.AllowedOrigins, origin) {
			c.AllowedOrigins = append(c.AllowedOrigins, origin)
		}
	}
	// c.AllowedOrigins = append(c.AllowedOrigins, origins...)
}

func (c *Cors) SetAllowedOrigins(origins []string) {
	c.AllowedOrigins = origins
}

func (c *Cors) AddAllowedMethods(methods ...string) {
	for _, method := range methods {
		if !contains(c.AllowedMethods, method) {
			c.AllowedMethods = append(c.AllowedMethods, method)
		}
	}
	// c.AllowedMethods = append(c.AllowedMethods, methods...)
}

func (c *Cors) SetAllowedMethods(methods []string) {
	c.AllowedMethods = methods
}

func (c *Cors) AddAllowedHeaders(headers ...string) {
	for _, header := range headers {
		if !contains(c.AllowedHeaders, header) {
			c.AllowedHeaders = append(c.AllowedHeaders, header)
		}
	}
	// c.AllowedHeaders = append(c.AllowedHeaders, headers...)
}

func (c *Cors) SetAllowedHeaders(headers []string) {
	c.AllowedHeaders = headers
}

func (c *Cors) SetAllowCredentials(allow bool) {
	c.AllowCredentials = allow
}

func (c *Cors) SetMaxAge(age int) {
	c.MaxAge = age
}

func (c *Cors) GetHeaders() map[string]string {
	headers := make(map[string]string)
	headers[AccessControlAllowOrigin] = strings.Join(c.AllowedOrigins, ",")
	headers[AccessControlAllowMethods] = strings.Join(c.AllowedMethods, ",")
	headers[AccessControlAllowHeaders] = strings.Join(c.AllowedHeaders, ",")
	if c.AllowCredentials {
		headers[AccessControlAllowCredentials] = "true"
	} else {
		headers[AccessControlAllowCredentials] = "false"
	}
	headers[AccessControlMaxAge] = strconv.Itoa(c.MaxAge)
	return headers
}

func (c *Cors) GetAllowedOrigins() []string {
	return c.AllowedOrigins
}

func (c *Cors) GetAllowedMethods() []string {
	return c.AllowedMethods
}

func (c *Cors) GetAllowedHeaders() []string {
	return c.AllowedHeaders
}

func (c *Cors) GetAllowCredentials() bool {
	return c.AllowCredentials
}

func (c *Cors) GetMaxAge() int {
	return c.MaxAge
}
