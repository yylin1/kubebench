// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package util provides various helper routines.
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Pformat returns a pretty format output of any value that can be marshaled to JSON.
func Pformat(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	valueJSON, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Warningf("Couldn't pretty format %v, error: %v", value, err)
		return fmt.Sprintf("%v", value)
	}
	return string(valueJSON)
}

// src is variable initialized with random value.
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString generates a random string of the desired length.
//
// The string is DNS-1035 label compliant; i.e. its only alphanumeric lowercase.
// From: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Run a subprocess and return stdout
func Run(command []string, dir string, env []string) ([]byte, error) {
	log.Infof("Execute command: %s", strings.Join(command, " "))
	if len(command) == 0 {
		return nil, errors.New("Command cannot be empty.")
	}
	name := command[0]
	args := command[1:]
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Command failed: %s", err)
		log.Errorf("Subprocess error:\n%s", string(stderr.Bytes()))
		log.Infof("Subprocess output:\n%s", string(stdout.Bytes()))
		return nil, err
	}
	output := stdout.Bytes()
	log.Infof("Subprocess output:\n%s", string(output))

	return output, err
}
