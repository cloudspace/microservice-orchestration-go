package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"
)

type Service struct {
    Id string
    DockerImage string
    Inputs map[string]map[string]string
}

// go run orchestrator.go "[{\"id\":\"1\",\"dockerImage\":\"cloudspace/url-lengthener-go\",\"inputs\":{\"url\":{\"service\":\"input\",\"key\":\"url\"}}},{\"id\":\"2\",\"dockerImage\":\"cloudspace/utm-stripper-go\",\"inputs\":{\"url\":{\"service\":\"1\",\"key\":\"url\"}}}]" "{\"url\":\"http://t.co/wnpJFiP7ls\"}"

func main() {
    var flow []Service
    orchestrationSpecJSON := os.Args[1]
    fmt.Println("ORCHESTRATION:", orchestrationSpecJSON)

    // parse the orchestration JSON
    orchestrationSpec := []byte(orchestrationSpecJSON)

    err1 := json.Unmarshal(orchestrationSpec, &flow)
    if err1 != nil {
        fmt.Println(err1)
        return
    }

    fmt.Println("FLOW:\n")
    fmt.Printf("%+v", flow)
    fmt.Println("\n")

    // parse user input
    fmt.Println("PARSING OPTIONS JSON")

    var options map[string]string
    optionsJSON := os.Args[2]
    fmt.Println("RAW OPTIONS:", optionsJSON)

    optionsBytes := []byte(optionsJSON)

    err2 := json.Unmarshal(optionsBytes, &options)
    if err2 != nil {
        fmt.Println(err2)
        return
    }

    fmt.Println("PARSED OPTIONS:")
    fmt.Printf("%+v", options)
    fmt.Println("\n")

    var user_input map[string]map[string]string

    // prelaunch
    var outputs map[string]map[string]chan string

    // outputs["0"] := map[string]chan

    for key, input := range options {
        chn := make(chan string, 1)
        chn <- input
        outputs["0"][key] = chn
    }
 
    /////////////////////////////////////////////////////
    fmt.Println("RUNNING...")

    output := user_input
    for index, service := range flow {
        var input map[string]chan string
        for key,value := range service.Inputs {
            input[key] = outputs[value["service"]][value["key"]]
        }

        outputs[service.Id] = runService(input, service.DockerImage)

        output = runService(output, service.DockerImage)
        fmt.Println("ID:",service.Id)
    }


    // c := make(chan string)
    
    // out := runService(c, flow[0].DockerImage)
    
    // fmt.Println("HELLO:",flow[0].Inputs["url"]["key"], "\n")

    // c <- flow[0].Inputs["url"]["key"]
    // fmt.Printf("%s", <-out)
}

func runService(inputs map[string]chan string, dockerImage string) <-chan string {
    result := make(chan string)
    go func() {
        args := <-arg1
        out, err := exec.Command("docker", "run", dockerImage, args).Output()
    
        if err != nil {
            fmt.Println("ERROR:")
            fmt.Println(err)
            return
        }

        result <- fmt.Sprintf("%s", out)
    }()
    return result
}
