package core

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Config struct {
	Workers int 'json:"workers"'
	Listen int 'json:"listen"'

	Log struct {
		Tank string 'json:"tank"'
		Level string 'json:"level"'
		File string 'json:"file"'
	}'json:"log"'
}

func (c *Config) validate() error {

