package main

import (
	"fmt"
	"math/rand"
	"time"

	"log"

	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

const maxDBNumber = 10000
const minDBNumber = 1

// Server that serves kuberdbs web service
type Server struct {
	port              int
	version           string
	engine            *gin.Engine
	redisConn         redis.Conn
	redisAddr         string
	redisAuthPassword string
}

type ServerConfig struct {
	port              int
	redisAddr         string
	redisAuthPassword string
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
		port:              config.port,
		version:           version,
		engine:            engine,
		redisConn:         redisConn,
		redisAddr:         config.redisAddr,
		redisAuthPassword: config.redisAuthPassword,
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
	// if redis password is set then supply it
	if s.redisAuthPassword != "" {
		if _, err := s.redisConn.Do("AUTH", s.redisAuthPassword); err != nil {
			s.redisConn.Close()
			fmt.Println(err)
			return
		}
	}
	// generate a random db number between 10000  and 1
	dbNumber := getDBNumber(maxDBNumber, minDBNumber)
	_, err := s.redisConn.Do("SELECT", dbNumber)
	if err != nil {
		s.redisConn.Close()
		fmt.Println(err)
		return
	}
	// create redis url
	redisURL := ""
	if s.redisAuthPassword != "" {
		redisURL = fmt.Sprintf("redis://%s@%s/%d", s.redisAuthPassword, s.redisAddr, dbNumber)
	} else {
		redisURL = fmt.Sprintf("redis://%s/%d", s.redisAddr, dbNumber)
	}

	c.String(http.StatusOK, fmt.Sprintf("REDIS_URL=%s", redisURL))
}

func getDBNumber(max int, min int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (s *Server) index(c *gin.Context) {
	// if redis password is set then supply it
	if s.redisAuthPassword != "" {
		if _, err := s.redisConn.Do("AUTH", s.redisAuthPassword); err != nil {
			fmt.Printf("couldn't issue AUTH command: %v \n", err)
			s.redisConn.Close()
			return
		}
	}
	pong, err := s.redisConn.Do("PING")
	if err != nil {
		fmt.Println(err)
	}
	pong, _ = redis.String(pong, nil)
	c.JSON(http.StatusOK, gin.H{"name": "kuberdbs", "description": "ondemand databases on top of kubernetes", "version": s.version, "redis": pong})
}
