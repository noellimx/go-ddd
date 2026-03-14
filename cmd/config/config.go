// Package config default configuration consumed from env vars etc.
package config

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
)

// Cors ...
type Cors struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

// Server ...
type Server struct {
	Cors `json:"cors"`

	Address string `json:"address"`
}

// Database ...
type Database struct {
	ConnString string `json:"connstring"`
}

// System ...
type System struct {
	Hostname    string `json:"hostname"`
	ContainerID string `json:"container_id"`
}

// InstanceID returns instance id.
func (s *System) InstanceID() string {
	return fmt.Sprintf("%s-%s", s.Hostname, s.ContainerID)
}

// SuperUser ...
type SuperUser struct {
}

// Admin ...
type Admin struct {
	SuperUser SuperUser `json:"superuser"`
}

// App application level config.
type App struct {
	Database `json:"database"`

	Server `json:"server"`

	Admin `json:"admin"`

	System `json:"system"`

	ServiceName string `json:"service_name"`
}

const cIDLen = 10

// ReadDefaultConfig from env.
func ReadDefaultConfig(ctx context.Context) (App, error) {
	l := slog.Default()
	var c App
	c.Address = os.Getenv("LISTENING_ADDRESS")
	c.AllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	// database
	dbURL := os.Getenv("DATABASE_URL")
	c.ConnString = dbURL

	// system
	hostname, err := os.Hostname()
	if err != nil {
		l.ErrorContext(ctx, "os.Hostname() missing.", slog.String("err", err.Error()))
		hostname = "unknown-hostname"
	}

	c.Hostname = hostname

	// generate random container ID
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, cIDLen)
	for i := range b {
		num, nErr := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if nErr != nil {
			return App{}, nErr
		}
		b[i] = letters[num.Int64()]
	}
	c.ContainerID = string(b)

	c.ServiceName = os.Getenv("SERVICE_NAME")

	if c.ServiceName == "" {
		c.ServiceName = "unknown-service"
	}

	return c, nil
}
