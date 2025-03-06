package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"http_tester/log"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var logFormat string
var port int

var RootCmd = &cobra.Command{
	Use:     "http_server [flags]",
	Run:     runCmd,
	Short:   "Simple HTTP server for testing purposes",
	Long:    `This simple tool listens on the indicated port and responds to POST requests on the /hello endpoint.`,
	Example: `  http_server --port 8090`,
}

func hello(w http.ResponseWriter, req *http.Request) {

	startTs := time.Now()

	//get the request method
	if req.Method == "POST" {

		// The request body just contain a string with a unique value
		// We can read the body of the request and print it
		readBody := make([]byte, req.ContentLength)
		err := error(nil)
		readBody, err = io.ReadAll(req.Body)
		if err != nil {
			log.ReportTime(startTs, "Error reading request body", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", readBody)))
		log.ReportTime(startTs, fmt.Sprintf("%s", readBody), nil)

	} else {
		log.Logger.Errorf("Only POST please. Request method: %s", req.Method)
		w.Write([]byte("Only POST please"))
		http.Error(w, "Only POST please", http.StatusBadRequest)
	}
}

func setCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 8090, "Interval between requests (regardless of the response time and number of clients)")
	cmd.PersistentFlags().StringVarP(&logFormat, "logformat", "l", "console", "log format (console or json)")

	_ = cmd.RegisterFlagCompletionFunc("logformat", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"console\tHuman readable format", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

}

func runCmd(cmd *cobra.Command, args []string) {
	if err := log.LogInit("INFO", logFormat); err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}
	http.HandleFunc("/hello", hello)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}
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
