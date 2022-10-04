// Copyright (c) 2020-2021, Aviatrix Systems, Inc.

package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

const CONFIG_PATH = "/etc/twamp/config.yml"

// Config struct holds the deserialized information from the config.yml file.
type Config struct {
	GrpcServer struct {
		GrpcHost string `yaml:"grpc_host"`
		GrpcPort string `yaml:"grpc_port"`
	} `yaml:"grpc_server"`
	TwampServer struct {
		TwampServerHost string `yaml:"twamp_server_host"`
		TwampServerPort string `yaml:"twamp_server_port"`
	} `yaml:"twamp_server"`
}

// ReadConfig returns a Config object after deserialization
func ReadConfig() Config {
	var config Config

	// Open YAML file
	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode YAML file to struct
	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Println(err.Error())
		}
	}

	return config
}
