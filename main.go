package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	portFlag              = "port"
	version               = "0.0.1"
	redisAddrFlag         = "redis-addr"
	redisAuthPasswordFlag = "redis-auth-password"
	mysqlAddrFlag         = "mysql-addr"
	mysqlUsernameFlag     = "mysql-username"
	mysqlPasswordFlag     = "mysql-password"

	defaultPort          = 8080
	defaultRedisAddr     = "localhost:6379"
	defaultMysqlAddr     = "localhost:3306"
	defaultMysqlUsername = "root"
)

func main() {
	app := cli.NewApp()
	configureCli(app)
	app.Action = mainAction
	app.Run(os.Args)
}

func mainAction(c *cli.Context) error {
	conf, err := validateConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return NewServer(conf).Start()
}

func configureCli(app *cli.App) {
	app.Name = "kuberdbs"
	app.Usage = "ondemand databases on top of kubernetes"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  fmt.Sprintf("%s, p", portFlag),
			Value: defaultPort,
			Usage: "The `port` to start the webserver on",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", redisAddrFlag),
			Value: defaultRedisAddr,
			Usage: "Redis server address",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", redisAuthPasswordFlag),
			Usage: "Redis authentication password",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", mysqlAddrFlag),
			Value: defaultMysqlAddr,
			Usage: "Mysql server address",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", mysqlUsernameFlag),
			Value: defaultMysqlUsername,
			Usage: "Mysql username",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", mysqlPasswordFlag),
			Usage: "Mysql password",
		},
	}
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

usage: {{.HelpName}} [options]
{{if .VisibleFlags}}
options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Version}}
version: {{.Version}}{{end}}
`
}

func validateConfig(c *cli.Context) (*ServerConfig, error) {
	// set defaults for non-required flags if not specified
	var port = defaultPort
	var redisAddr = defaultRedisAddr
	var redisAuthPassword = os.Getenv("REDIS_PASSWORD")
	var mysqlAddr = defaultMysqlAddr
	var mysqlUsername = defaultMysqlUsername
	var mysqlPassword = os.Getenv("MYSQL_PASSWORD")

	if c.IsSet(portFlag) {
		port = c.Int(portFlag)
	}
	if c.IsSet(redisAddrFlag) {
		redisAddr = c.String(redisAddrFlag)
	}
	if c.IsSet(redisAuthPasswordFlag) {
		redisAuthPassword = c.String(redisAuthPasswordFlag)
	}
	if c.IsSet(mysqlAddrFlag) {
		mysqlAddr = c.String(mysqlAddrFlag)
	}
	if c.IsSet(mysqlUsernameFlag) {
		mysqlUsername = c.String(mysqlUsernameFlag)
	}
	if c.IsSet(mysqlPasswordFlag) {
		mysqlPassword = c.String(mysqlPasswordFlag)
	}

	return &ServerConfig{
		port:              port,
		redisAddr:         redisAddr,
		redisAuthPassword: redisAuthPassword,
		mysqlAddr:         mysqlAddr,
		mysqlUsername:     mysqlUsername,
		mysqlPassword:     mysqlPassword,
	}, nil
}
