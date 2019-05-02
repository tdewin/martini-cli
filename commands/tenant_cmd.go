package commands

import (
	"fmt"
	"log"
	"syscall"

	"github.com/tdewin/martini-cli/core"
	"github.com/tdewin/martini-cli/tenant"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

//seperation without using global variables
func GetTenantCommands() *cli.Command {
	return &cli.Command{
		Name:    "tenant",
		Aliases: []string{"t"},
		Usage:   "tenant management",
		Subcommands: []cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create tenant",
				Action: func(c *cli.Context) error {
					err := ValidateArray([]ValidString{
						ValidString{c.String("tenant"), "tenant", "."},
						ValidString{c.String("email"), "email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`},
						ValidString{c.String("fqdn"), "fqdn", `.`},
						ValidString{c.String("port"), "port", `[0-9]+`},
						ValidString{c.String("username"), "username", `.`},
					})
					if err == nil {
						conn := core.NewConnectionFromCLIContext(c)
						err = conn.Auth(nil, false)
						if err == nil {
							pw := c.String("password")
							if pw == "" {
								fmt.Print("Enter tenant server password: ")
								dbbytePassword, errp := terminal.ReadPassword(int(syscall.Stdin))
								for errp != nil || len(string(dbbytePassword)) < 3 {
									fmt.Println()
									fmt.Print("Password can not be empty (min 3 char):")
									dbbytePassword, errp = terminal.ReadPassword(int(syscall.Stdin))
								}
								pw = string(dbbytePassword)
							}

							t := tenant.MartiniTenant{c.String("tenant"), c.String("email"), c.String("fqdn"), c.String("port"), c.String("username"), pw, "-1"}

							err = t.Create(conn)
							if err != nil {
								log.Println("Error ", err)
							}
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "tenant, t",
						Value: "",
						Usage: "Tenant name",
					},
					cli.StringFlag{
						Name:  "email, e",
						Value: "",
						Usage: "Email",
					},
					cli.StringFlag{
						Name:  "fqdn, f",
						Value: "",
						Usage: "FQDN instance",
					},
					cli.StringFlag{
						Name:  "port",
						Value: "4443",
						Usage: "FQDN port",
					},
					cli.StringFlag{
						Name:  "username, u",
						Value: "",
						Usage: "Username instance",
					},
					cli.StringFlag{
						Name:  "password, p",
						Value: "",
						Usage: "Password instance",
					},
				},
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "deploy tenant (will create a new installation instead of just adding it to martini)",
				Subcommands: []cli.Command{
					{
						Name:    "amazon",
						Aliases: []string{"a"},
						Usage:   "deploy an amazon EC2 container",
						Action: func(c *cli.Context) error {
							var err error

							err = ValidateArray([]ValidString{
								ValidString{c.GlobalString("tenant"), "tenant", "."},
								ValidString{c.GlobalString("email"), "email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`},
								ValidString{c.String("region"), "region", `.`},
							})
							if err == nil {
								conn := core.NewConnectionFromCLIContext(c)
								err = conn.Auth(nil, false)
								if err == nil {
									t := tenant.NewAWSConfig(c.GlobalString("tenant"), c.GlobalString("email"), c.String("region"))

									err = t.Deploy(conn)
								}
							}
							return err
						},
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "region",
								Value: "",
								Usage: "AWS Region",
							},
						},
					},
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "tenant, t",
						Value: "",
						Usage: "Tenant name",
					},
					cli.StringFlag{
						Name:  "email, e",
						Value: "",
						Usage: "Email",
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tenants",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth(nil, false)
					if err == nil {
						tenants, err := tenant.List(conn)
						if err == nil {

							for i := 0; i < 12; i++ {
								fmt.Print("##########")
							}

							for _, t := range tenants {
								fmt.Printf("\n| %5s | %15s | %29s | %30s | %25s |", t.Id, t.Name, t.Email, t.Instancefqdn, t.Instanceusername)
							}
							fmt.Print("\n")
							for i := 0; i < 12; i++ {
								fmt.Print("##########")
							}
							fmt.Print("\n")
						} else {
							fmt.Println(err)
						}
					}
					return err
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"x"},
				Usage:   "delete a tenant",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth(nil, false)
					if err == nil {
						err := tenant.Delete(conn, c.String("id"))
						if err != nil {
							fmt.Println(err)
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
					},
				},
			},
			{
				Name:    "broker",
				Aliases: []string{"b"},
				Usage:   "broker an rdp connection via the martini server to a tenant",
				Action: func(c *cli.Context) error {
					conn := core.NewConnectionFromCLIContext(c)
					err := conn.Auth(nil, false)
					if err == nil {
						err := tenant.Broker(conn, c.String("id"), c.String("clientip"))
						if err != nil {
							fmt.Println(err)
						}
					}
					return err
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "Id of tenant",
					},
					cli.StringFlag{
						Name:  "clientip, c",
						Value: "",
						Usage: "IP of your local break-out towards the server. If empty, the server will try to autodetect",
					},
				},
			},
			//more commands indent here
		},
	}
}
