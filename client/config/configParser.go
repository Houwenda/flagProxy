package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
)

type Config struct {
	ChallengeConf Challenge `yaml:"challenge"`
	ServerConf    Server    `yaml:"server"`
	LogConf       Log       `yaml:"log"`
}

type Challenge struct {
	Address   string `yaml:"address"`
	FlagRegex string `yaml:"flag_regex"`
}

type Server struct {
	Url         string `yaml:"url"`
	ChallengeId string `yaml:"challenge_id"`
	Key         string `yaml:"key"`
}

type Log struct {
	Path string `yaml:"path"`
}

func (c *Config) Parse(path string) (*Config, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// check file existence
	if _, existErr := os.Stat(path); existErr != nil && os.IsNotExist(existErr) {
		return nil, errors.New("config file does not exist at " + path)
	}
	// parse file
	err = yaml.Unmarshal(configFile, c)
	if err != nil {
		return nil, errors.New("config file cannot be parsed, check your syntax ")
	}
	return c, err
}

func (c *Config) Validate() error {
	// challenge
	// valid challenge address
	if _, err := net.ResolveTCPAddr("tcp4", c.ChallengeConf.Address); err != nil {
		return err
	}
	// valid regex expression
	if _, err := regexp.Compile(c.ChallengeConf.FlagRegex); err != nil {
		return err
	}

	// server
	// connection check
	resp, err := http.Get(c.ServerConf.Url)
	if err != nil {
		return err
	}
	data := make([]byte, 1024)
	if _, _ = resp.Body.Read(data); string(data[:9]) != "flagProxy" {
		fmt.Println("server returns :", []byte(string(data[:9])))
		return errors.New("invalid server url")
	}
	return nil
}
