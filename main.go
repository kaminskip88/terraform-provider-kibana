package main

import (
	"context"
	"flag"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/kaminskip88/terraform-provider-kibana/v8/kb"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%\n",
	})

}

func main() {

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: kb.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/disaster37/kibana", opts)

		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)

}
