package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"net/http"

	"github.com/argoproj/argo-workflows/v3/util/errors"
	log "github.com/sirupsen/logrus"

	// load authentication plugin for obtaining credentials from cloud providers.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/argoproj/argo-workflows/v3/cmd/argoexec/commands"
	"github.com/argoproj/argo-workflows/v3/util"

	pprofutil "github.com/argoproj/argo-workflows/v3/util/pprof"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	pprofutil.Init()
    go func() {
	    if os.Getenv("ARGO_CONTAINER_NAME") == "wait" {
            log.Println(http.ListenAndServe(":6060", nil))
        } else {
            log.Info("not listening, this is not container wait")
        }
    }()

	err := commands.NewRootCommand().ExecuteContext(ctx)
	if err != nil {
		if exitError, ok := err.(errors.Exited); ok {
			if exitError.ExitCode() >= 0 {
				os.Exit(exitError.ExitCode())
			} else {
				os.Exit(137) // probably SIGTERM or SIGKILL
			}
		} else {
			util.WriteTerminateMessage(err.Error()) // we don't want to overwrite any other message
			println(err.Error())
			os.Exit(64)
		}
	}
}
