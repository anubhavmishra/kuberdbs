package main

import (
	"fmt"

	"log"

	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

// Server that serves kuberdbs web service
type Server struct {
	port      int
	version   string
	engine    *gin.Engine
	redisConn redis.Conn
}

type ServerConfig struct {
	port      int
	redisAddr string
}

func NewServer(config *ServerConfig) *Server {
	// create a new engine
	engine := gin.New()

	// redis connection
	redisConn, err := redis.Dial("tcp", config.redisAddr)
	if err != nil {
		log.Fatalf("error connecting to redis server %s", err)
	}

	return &Server{
		port:      config.port,
		version:   version,
		engine:    engine,
		redisConn: redisConn,
	}
}

func (s *Server) Start() error {
	s.engine.Use(gin.Recovery(), gin.Logger())
	s.engine.GET("/", s.index)
	s.engine.GET("/redis", s.redis)
	log.Printf("kuberdbs started - listening on port %v", s.port)
	if err := s.engine.Run(fmt.Sprintf(":%v", s.port)); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func (s *Server) redis(c *gin.Context) {
	r, err := s.redisConn.Do("SELECT", 1)
	if err != nil {
		fmt.Println(err)
	}
	database, _ := redis.String(r, nil)
	c.JSON(http.StatusOK, gin.H{"database": database})
}

func (s *Server) index(c *gin.Context) {
	pong, err := s.redisConn.Do("PING")
	if err != nil {
		fmt.Println(err)
	}
	pong, _ = redis.String(pong, nil)
	c.JSON(http.StatusOK, gin.H{"name": "kuberdbs", "description": "ondemand databases on top of kubernetes", "version": s.version, "redis": pong})
}
