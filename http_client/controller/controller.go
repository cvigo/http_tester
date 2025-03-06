package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"http_tester/log"

	"github.com/google/uuid"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var conns int
var url string
var interval time.Duration
var logFormat string

var RootCmd = &cobra.Command{
	Use:     "http_client [flags]",
	Run:     runCmd,
	Short:   "Simple HTTP client for testing purposes",
	Long:    `This simple tool sends a POST request to the indicated URL, with a random UUID string as the body.`,
	Example: `  http_client --url localhost:8090 --connections 10 --interval 1s --logformat json`,
}

func startClient(doSend <-chan any) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/hello", url)

	for range doSend {
		uuidStr := uuid.New().String()
		startTs := time.Now()
		resp, err := client.Post(url, "text/plain", strings.NewReader(uuidStr))
		if err != nil {
			log.ReportTime(startTs, "Error sending request", err)
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.ReportTime(startTs, "Error reading response", err)
			continue
		}
		log.ReportTime(startTs, fmt.Sprintf("%v", string(respBytes)), nil)
		resp.Body.Close()
	}
}

func runCmd(cmd *cobra.Command, args []string) {

	if err := log.LogInit("INFO", logFormat); err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}

	// prepare for nice stop in case con CTRL+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	doSend := make(chan any)

	var wg sync.WaitGroup
	for i := 1; i <= conns; i++ {
		wg.Add(1)
		go startClient(doSend)
	}

	t := time.NewTicker(interval)
	ticker := t.C
	defer t.Stop()

Loop:
	for {
		select {
		case <-sigs: //break
			break Loop
		case <-ticker:
			doSend <- struct{}{}
		}
	}

	close(doSend)
	time.Sleep(1 * time.Second)
}

func setCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&conns, "connections", "c", 1, "number of client connections working concurrently")
	cmd.PersistentFlags().StringVarP(&url, "url", "u", "localhost:8090", "URL to connect to")
	cmd.PersistentFlags().DurationVarP(&interval, "interval", "i", 1*time.Second, "Interval between requests (regardless of the response time and number of clients)")
	cmd.PersistentFlags().StringVarP(&logFormat, "logformat", "l", "console", "log format (console or json)")

	_ = cmd.RegisterFlagCompletionFunc("logformat", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"console\tHuman readable format", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

}

func init() {
	setCommonFlags(RootCmd)

	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiBlue + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})
}
