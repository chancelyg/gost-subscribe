package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const Version = "V24.06.04"

type Server struct {
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

type Response struct {
	Servers []Server `json:"servers"`
}

// LogInit initializes the logrus logger
func LogInit() {
	bytesWriter := &bytes.Buffer{}
	stdoutWriter := os.Stdout
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z",
		FullTimestamp:   true})
	log.SetOutput(io.MultiWriter(bytesWriter, stdoutWriter))
	log.SetLevel(log.InfoLevel)
}

func IsSubstringPresent(slice []string, parentString string) bool {
	for _, v := range slice {
		if strings.Contains(parentString, v) {
			return true
		}
	}
	return false
}

func main() {
	// Define command line flags
	hlepFlag := flag.Bool("h", false, "--help")
	urlSubscription := flag.String("u", "", "A shadowsocks format URL subscription link")
	serverProxyPort := flag.Int("p", 1080, "Specyify port of the proxy")
	proxyStrategy := flag.String("s", "fifo", "Proxy strategy [round|rand|fifo|hash]")
	proxyFailTimeout := flag.Int("t", 600, "Specify the timeout duration in seconds")
	proxyMaxFails := flag.Int("m", 1, "Specify the number of failed attempts")
	printVersion := flag.Bool("V", false, "Show version")
	filterSubscribedKeywords := flag.String("f", "套餐|重置|剩余|更新", "Filter the results to include keywords separated by '|'")
	containsKeywords := flag.String("k", "", "Only return results that contain the keywords separated by '|'")
	outputPath := flag.String("o", "config.yml", "Output file path")

	if *hlepFlag {
		flag.Usage()
		return
	}

	if *printVersion {
		fmt.Println(fmt.Printf("Current version: %s", Version))
		return
	}

	flag.Parse()

	if *urlSubscription == "" {
		flag.Usage()
		log.Fatalln("please check if your parameter inputs are correct.")
	}

	LogInit()

	log.WithFields(log.Fields{
		"serverProxyPort":  *serverProxyPort,
		"proxyStrategy":    *proxyStrategy,
		"proxyFailTimeout": *proxyFailTimeout,
		"proxyMaxFails":    *proxyMaxFails,
		"outputPath":       *outputPath,
	}).Info("init successful")

	resp, err := http.Get(*urlSubscription)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response
	json.Unmarshal(body, &response)

	log.WithFields(log.Fields{"servers": len(response.Servers)}).Info("http request successful")

	var nodes []map[string]interface{}
outerLoop:
	for _, server := range response.Servers {
		keys := strings.Split(*filterSubscribedKeywords, "|")
		if *filterSubscribedKeywords != "" && IsSubstringPresent(keys, server.Remarks) {
			log.WithFields(log.Fields{"remarks": server.Remarks}).Info("ignore server")
			continue outerLoop
		}
		keys = strings.Split(*containsKeywords, "|")
		if *containsKeywords != "" && !IsSubstringPresent(keys, server.Remarks) {
			log.WithFields(log.Fields{"remarks": server.Remarks}).Info("ignore server")
			continue outerLoop
		}
		node := map[string]interface{}{
			"name": server.Remarks,
			"addr": fmt.Sprintf("%s:%d", server.Server, server.ServerPort),
			"connector": map[string]interface{}{
				"type": "ss",
				"auth": map[string]string{
					"username": server.Method,
					"password": server.Password,
				},
			},
		}
		nodes = append(nodes, node)
		log.WithFields(log.Fields{"name": server.Remarks}).Info("add subscription success")
	}

	yamlData := map[string]interface{}{
		"services": []map[string]interface{}{
			{
				"name": strconv.Itoa(*serverProxyPort),
				"addr": fmt.Sprintf(":%d", *serverProxyPort),
				"handler": map[string]string{
					"type":  "auto",
					"chain": "chain-0",
				},
				"listener": map[string]string{
					"type": "tcp",
				},
			},
		},
		"chains": []map[string]interface{}{
			{
				"name": "chain-0",
				"hops": []map[string]interface{}{{
					"name": "hop-0",
					"selector": map[string]interface{}{
						"strategy":    proxyStrategy,
						"maxFails":    strconv.Itoa(*proxyMaxFails),
						"failTimeout": fmt.Sprintf("%ds", *proxyFailTimeout),
					},
					"nodes": nodes,
				},
				},
			},
		},
	}

	file, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(yamlData)
	if err != nil {
		log.Fatalln(err)
	}
	defer encoder.Close()

	log.WithFields(log.Fields{
		"outputPath":    *outputPath,
		"serverNumbers": len(nodes),
	}).Info("write successful")

}
