// Copyright 2019 John Darrington johnw.darrington@gmail.com

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package heimdall encompasses all functions related to the short-lived process manager by the same name. Heimdall
// was the ever-vigilant guardian of the gods' stronghold, Asgard - now he will be the guardian of whichever program you choose.
// Heimdall is designed as both launcher and monitor of short-lived CLI tools and programs. Heimdall provides the ability
// to automatically repeat a process, kill a hung process started with the tool, and log the programs output (filtering logs
// is also possible). It is hoped that heimdall and bifrost will be a tool you reach for again and again when developing your CLI tool.

package main

import "github.com/dnoberon/heimdall/cmd"

func main() {
	cmd.Execute()
}
