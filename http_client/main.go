package main

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

	"github.com/google/uuid"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Conns int
var URL string
var Interval time.Duration

var rootCmd = &cobra.Command{
	Use: "http_client [flags]",
	Run: runCmd,
}

func startClient(doSend <-chan any) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/hello", URL)

	for range doSend {
		uuidStr := uuid.New().String()
		startTs := time.Now()
		resp, err := client.Post(url, "text/plain", strings.NewReader(uuidStr))
		if err != nil {
			reportTime(startTs, "Error sending request", err)
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			reportTime(startTs, "Error reading response", err)
			continue
		}
		reportTime(startTs, fmt.Sprintf("%v", string(respBytes)), nil)
		resp.Body.Close()
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCmd(cmd *cobra.Command, args []string) {

	// prepare for nice stop in case con CTRL+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	doSend := make(chan any)

	var wg sync.WaitGroup
	for i := 1; i <= Conns; i++ {
		wg.Add(1)
		go startClient(doSend)
	}

	t := time.NewTicker(Interval)
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

func reportTime(startTs time.Time, msg string, err error) {
	timeNow := time.Now()
	if err != nil {
		logrus.Errorf("Start:%v, Finish:%v, Elapsed(ms):%v Message: %s", startTs, timeNow, float64(timeNow.Sub(startTs).Nanoseconds())/1e6, err)
	} else {
		logrus.Infof("Start:%v, Finish:%v, Elapsed(ms):%v Message: %s", startTs, timeNow, float64(timeNow.Sub(startTs).Nanoseconds())/1e6, msg)
	}
}

func setCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&Conns, "connections", "c", 1, "number of client connections working concurrently")
	cmd.PersistentFlags().StringVarP(&URL, "url", "u", "localhost:8090", "URL to connect to")
	cmd.PersistentFlags().DurationVarP(&Interval, "interval", "i", 1*time.Second, "Interval between requests (regardless of the response time and number of clients)")
}

func init() {
	setCommonFlags(rootCmd)

	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiBlue + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})
}
