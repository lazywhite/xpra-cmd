//
// Copyright (C) 2020 white <white@Whites-Mac-Air.local>
//
// Distributed under terms of the MIT license.
//

package main

import (
    "time"
    "errors"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
    //"reflect"
    "encoding/json"
    //"fmt"
    "runtime"
    "crypto/tls"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/Nerzal/gocloak/v4"
    "github.com/go-resty/resty/v2"
)

var mainLogger *log.Logger
var mainLogFile *os.File

func initMainLogger(){
    mainLog := getLogFile("command.log")
    mainLogger = log.New(mainLog, "[INFO]", log.LstdFlags)
}

//get keycloak openid
func authUser(token string) (ok bool, err error){
    if token == ""{
        ok = false
        err = errors.New("no authorization header provided")
        return
    }
    realmName := os.Getenv("REALM_NAME")
    serverURL := os.Getenv("SERVER_URL")

    client := gocloak.NewClient(serverURL)
    restyClient := client.RestyClient()
    //restyClient.SetDebug(true)
    restyClient.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })

    userInfo, _err := client.GetUserInfo(token, realmName)
    if userInfo != nil && _err == nil{
        ok = true
        err = nil
        return
    }else{
        ok = false
        err = errors.New("invalid token")
        return
    }
}

func getLogFile(name string) *os.File {
    logFile, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        log.Fatalln("open file error !")
    }
    return logFile
}

func getAllowedCommand(token string) []string {
    apiURL := os.Getenv("API_URL")
    path := "avlcloud/api/apps/?enabled=true"
    url := apiURL + path

    //cmds := []string{"xclock"}

    var cmds []string
    client := resty.New()
    resp, err := client.
                    SetTimeout(3 * time.Second).
                    R().
                    //EnableTrace().
                    SetHeader("Accept", "application/json").
                    SetHeader("Authorization", token).
                    Get(url)
    if err == nil && resp.StatusCode() == 200{
        //fmt.Println(reflect.TypeOf(resp.Body()))
        err := json.Unmarshal(resp.Body(), &cmds)
        if err != nil {
            mainLogger.SetPrefix("[Error]")
            mainLogger.Printf("failed to parse reponse body: %v", resp.Body())
        }
    }else{
        mainLogger.SetPrefix("[Error]")
        mainLogger.Printf("failed to get response: %s", url)
    }
    cmds = append(cmds, "xclock")
    return cmds
}

func inTest(ele string, all []string) bool {
    for _, v := range all {
        if ele == v {
            return true
        }
    }
    return false
}
func runCommand(cmd string, token string) (status int, err error) {
    if runtime.GOOS == "windows" {
        status = 1
        err = errors.New("apiServer platform not supported")
        return
    }
    if cmd != "" {
        cmds := getAllowedCommand(token)
        if inTest(cmd, cmds) {
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
    token := r.Header.Get("Authorization")
    //1. auth user
    if ok, err := authUser(token); ok{
        var info CommandInfo
        err := json.NewDecoder(r.Body).Decode(&info)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        mainLogger.SetPrefix("[INFO]")
        mainLogger.Printf("executing command: '%s'", info.CMD)

        //2. execute command
        status, err := runCommand(info.CMD, token)
        if status == 0 && err == nil {
            json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "success"})
        } else {
            json.NewEncoder(w).Encode(map[string]interface{}{"code": status, "msg": err.Error()})
        }

    }else{
        json.NewEncoder(w).Encode(map[string]interface{}{"code": 5, "msg": err.Error()})

    }
}

func main() {
    accessLogFile := getLogFile("access.log")
    initMainLogger()
    defer accessLogFile.Close()
    defer mainLogFile.Close()

    r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/launch", launchHandler).Methods("POST")
    loggedRouter := handlers.LoggingHandler(accessLogFile, r)
    http.ListenAndServe(":11000", loggedRouter)
}
