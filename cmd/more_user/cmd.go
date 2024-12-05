package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofrp/fp-multiuser/pkg/server"

	"github.com/spf13/cobra"
)

const version = "0.1.1"

var (
	showVersion bool

	bindAddr    string
	tokenFile   string
	remoteToken string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "l", "127.0.0.1:7200", "bind address")
	rootCmd.PersistentFlags().StringVarP(&tokenFile, "token_file", "f", "./tokens", "token file")
	rootCmd.PersistentFlags().StringVarP(&remoteToken, "remote_token", "m", "", "remote token file")
}

// rootCmd represents the root command for fp-multiuser.
var rootCmd = &cobra.Command{
	Use:   "HeloFrp_MU",
	Short: "MU is the server plugin of frp to support multiple users based on fp-multiuser",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the version flag is set
		if showVersion {
			fmt.Println(version)
			return nil
		}

		// Parse tokens from the specified file
		var tokens map[string]string
		var err error
		if remoteToken != "" {
			tokens, err = ParseTokensFromRemote(remoteToken)
		} else {
			tokens, err = ParseTokensFromFile(tokenFile)
		}
		if err != nil {
			log.Printf("parse tokens error: %v", err)
			return err
		}

		// Create a new server instance with the specified configuration
		s, err := server.New(server.Config{
			BindAddress: bindAddr,
			Tokens:      tokens,
		})
		if err != nil {
			return err
		}

		// Start the server
		s.Run()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ParseTokensFromFile(file string) (map[string]string, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	rows := strings.Split(string(buf), "\n")
	for _, row := range rows {
		kvs := strings.SplitN(row, "=", 2)
		if len(kvs) == 2 {
			ret[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
		}
	}
	return ret, nil
}

func ParseTokensFromRemote(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	rows := strings.Split(string(buf), "\n")
	for _, row := range rows {
		kvs := strings.SplitN(row, "=", 2)
		if len(kvs) == 2 {
			ret[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
		}
	}
	return ret, nil
}
