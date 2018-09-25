package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"time"

	"./logs"

	"github.com/BurntSushi/toml"
	"github.com/sethvargo/go-password/password"
)

type serverConfig struct {
	Config      bool              `toml:"configured" json:"configured"`
	ConfigPath  string            `toml:"-" json:"-"`
	Port        string            `toml:"port" json:"-"`
	Logger      *logs.Logs        `toml:"logger"`
	TokenConfig *tokenKeys        `toml:"tokenConfig" json:"-"`
	UserConfig  *modifiableConfig `toml:"modifiableConfig" json:"user_config"`
}

type tokenKeys struct {
	Signature string            `toml:"signature" json:"-"`
	AccessMap map[string]string `toml:"accessMap" json:"access_map"`
}

type modifiableConfig struct {
	Name            string   `toml:"name" json:"server_name"`
	Debug           bool     `toml:"debug" json:"debug"`
	RequireAuth     bool     `toml:"require_auth" json:"require_auth"`
	AuthKey         string   `toml:"auth_key" json:"auth_key"`
	TokenExpiration string   `toml:"tokenExpiration" json:"token_expiration"`
	AccessModes     []string `toml:"accessModes" json:"access_modes"`
	goDuration      *time.Duration
}

func parseConfig(filePath string) {
	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		panic(err)
	}

	if c.Config {
		expirationDuration, err := time.ParseDuration(c.UserConfig.TokenExpiration)
		if err != nil {
			panic(err)
		}
		c.UserConfig.goDuration = &expirationDuration
	}

	saveConfig(filePath, &c)
}

func updateConfig(newConfig modifiableConfig) error {

	c.UserConfig.Name = newConfig.Name
	c.UserConfig.AccessModes = newConfig.AccessModes
	c.UserConfig.Debug = newConfig.Debug
	c.UserConfig.RequireAuth = newConfig.RequireAuth

	if c.UserConfig.RequireAuth {
		if newConfig.AuthKey == "" {
			return errors.New("must update \"auth_key\" field")
		}
		c.UserConfig.AuthKey = newConfig.AuthKey
	}

	if !c.Config {
		c.TokenConfig.AccessMap = make(map[string]string)
	}

	newAccess := make(map[string]string)
	for _, mode := range newConfig.AccessModes {
		sig, err := password.Generate(32, 10, 20, false, true)
		if err != nil {
			return errors.New("Could not generate pass for access modes")
		}
		if c.TokenConfig.AccessMap[mode] == "" {
			newAccess[mode] = sig
		} else {
			newAccess[mode] = c.TokenConfig.AccessMap[mode]
		}
	}
	c.TokenConfig.AccessMap = newAccess

	expirationDuration, err := time.ParseDuration(newConfig.TokenExpiration)
	if err != nil {
		return err
	}
	c.UserConfig.goDuration = &expirationDuration

	if c.TokenConfig.Signature == "" {
		sig, err := password.Generate(32, 10, 20, false, true)
		if err != nil {
			return errors.New("Could not generate JWT signature")
		}
		c.TokenConfig.Signature = sig
	}

	c.Config = true

	return nil
}

func saveConfig(filePath string, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	w.Flush()
	f.Close()

	return nil
}
