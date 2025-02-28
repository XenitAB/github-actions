package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/urfave/cli/v2"
	"github.com/xenitab/github-actions/docker/go-tf-prepare/pkg/azure"
)

func main() {
	stdr.SetVerbosity(1)
	log := stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))
	log = log.WithName("tf-prepare")

	ctx := logr.NewContext(context.Background(), log)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "azure",
				Usage: "Terraform prepare for Azure",
				Flags: azure.Flags(),
				Action: func(cli *cli.Context) error {
					err := azure.Action(ctx, cli)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err, "CLI execution failed")
		os.Exit(1)
	}

	os.Exit(0)
}
