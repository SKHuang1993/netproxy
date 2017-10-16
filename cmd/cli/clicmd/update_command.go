package clicmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fagongzi/netproxy/pkg/conf"
	"github.com/fagongzi/netproxy/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	inLossRate   int
	inDelayMs    int
	outLossRate  int
	outDelayMs   int
	updateClient = false
)

// NewUpdateCommand returns the cobra command for "update".
func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <proxyAddr>",
		Short: "Update the proxy ctl",
		Run:   updateCommandFunc,
	}

	cmd.Flags().IntVar(&inLossRate, "in-lossRate", 0, "set the client receive packet loss rate.")
	cmd.Flags().IntVar(&inDelayMs, "in-delayMs", 0, "set the client receive packet delay.")
	cmd.Flags().IntVar(&outLossRate, "out-lossRate", 0, "set the client sent packet loss rate.")
	cmd.Flags().IntVar(&outDelayMs, "out-delayMs", 0, "set the client sent packet delay.")

	return cmd
}

// updateCommandFunc executes the "update" command.
func updateCommandFunc(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("http://%s%s", Global.Endpoints, proxy.APIProxies)

	cli := &http.Client{}
	ctl := &conf.Ctl{
		Address: args[0],
		In: &conf.CtlUnit{
			LossRate: inLossRate,
			DelayMs:  inDelayMs,
		},
		Out: &conf.CtlUnit{
			LossRate: outLossRate,
			DelayMs:  outDelayMs,
		},
	}
	data := ctl.Marshal()
	fmt.Printf("<%s> send to server\n", data)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(data))
	rsp, err := cli.Do(request)

	if err != nil {
		fmt.Println(err)
	} else {
		defer rsp.Body.Close()
		data, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("server return %s\n", data)
		}
	}
}
