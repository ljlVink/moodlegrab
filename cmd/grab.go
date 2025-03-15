/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"moodlegrab/moodlegrab"
	grab "moodlegrab/moodlegrab"
	"net/http"
	gourl "net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var url string

var grabCmd = &cobra.Command{
	Use:   "grab",
	Short: "start grab",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grab called")
		if url == "" {
			fmt.Println("Error: url flag is required")
			return
		}
		if !strings.HasPrefix(url, "http") {
			fmt.Println("Error: url must start with http")
			return
		}
		fmt.Println("URL to grab:", url)
		proxyURL, err := gourl.Parse("http://127.0.0.1:8080") // 替换为你的代理地址
		if err != nil {
			fmt.Println("Error parsing proxy URL:", err)
			return
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
		var config moodlegrab.Yaml_config
		viper.Unmarshal(&config)

		account := config.General.Account
		passwd := config.General.Passwd
		if account == "" || passwd == "" {
			log.Fatalln("account or passwd is null! check config!!")
		}
		grabclient := &grab.GrabClient{
			MoodleUrl: url,
			UserName:  account,
			Passwd:    passwd,
			Client: http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
					TLSClientConfig: &tls.Config{
						MaxVersion: tls.VersionTLS12,
						CipherSuites: []uint16{
							tls.TLS_RSA_WITH_AES_256_CBC_SHA,
						},
					},
				},
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}
		err = grabclient.Login()
		if err != nil {
			log.Fatalln(err)
		}
		err = grabclient.GrepCourses()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(grabCmd)
	grabCmd.Flags().StringVarP(&url, "url", "u", "https://moodle.smbu.edu.cn", "URL to grab")
}
