//
// Copyright (C) 2020 white <white@Whites-Mac-Air.local>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	//"reflect"
	"encoding/json"
	//"fmt"
	"runtime"
)

const AccessLog = "access.log"
const MainLog = "command.log"

type Info struct {
	CMD string `json:cmd`
}

func getLogFile(name string) *os.File {
	logFile, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("open file error !")
	}
	return logFile
}

func get_allowed_cmd() []string {
	// TODO: http get request
	cmds := []string{"ls", "ls -al", "sleep 3"}
	return cmds
}

func check_in(ele string, all []string) bool {
	for _, v := range all {
		if ele == v {
			return true
		}
	}
	return false
}
func run_cmd(cmd string) (status int, err error) {
	if runtime.GOOS == "windows" {
		status = 1
		err = errors.New("apiServer platform not supported")
		return
	}
	if cmd != "" {
		cmds := get_allowed_cmd()
		if check_in(cmd, cmds) {
			r := strings.Split(cmd, " ")
			cmd_name := r[0]
			cmd_args := r[1:]
			OS_CMD := exec.Command(cmd_name, cmd_args...)
			c_err := OS_CMD.Run()
			if c_err == nil {
				status = 0
				err = nil
				return
			} else {
				status = 5
				err = errors.New("failed to execute command")
				return
			}
		} else {
			status = 2
			err = errors.New("command is not allowed to execute")
			return
		}
	} else {
		status = 3
		err = errors.New("no command provided")
		return
	}
}
func launchHandler(w http.ResponseWriter, r *http.Request) {
	mainLog := getLogFile(MainLog)
	defer mainLog.Close()
	logger := log.New(mainLog, "[INFO]", log.LstdFlags)

	// TODO; authen check
	//fmt.Println(r.Header.Get("User-Agent"))
	//token := r.Header.Get("Authorization"))
	//json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	var info Info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(reflect.ValueOf(info.CMD))
	//fmt.Println(reflect.TypeOf(info.CMD))
	logger.Printf("executing command: '%s'", info.CMD)
	status, err := run_cmd(info.CMD)
	if status == 0 && err == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "success"})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": status, "msg": err.Error()})
	}
}

func main() {
	accessLog := getLogFile("access.log")
	defer accessLog.Close()

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/launch", launchHandler).Methods("POST")
	loggedRouter := handlers.LoggingHandler(accessLog, r)
	http.ListenAndServe(":11000", loggedRouter)
}
