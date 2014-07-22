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
    Outputs []string
}

// go run orchestrator.go "[{\"id\":\"1\",\"dockerImage\":\"cloudspace/url-lengthener-go\",\"inputs\":{\"url\":{\"service\":\"input\",\"key\":\"url\"}},\"outputs\":[\"url\"]},{\"id\":\"2\",\"dockerImage\":\"cloudspace/utm-stripper-go\",\"inputs\":{\"url\":{\"service\":\"1\",\"key\":\"url\"}},\"outputs\":[\"url\"]}]" "{\"url\":\"http://t.co/wnpJFiP7ls\"}"


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

    // prelaunch
    outputs := make(map[string]map[string]chan string)

    // outputs["0"] := map[string]chan

    user_input := make(map[string]chan string)
    for key, input := range options {
        chn := make(chan string, 1)
        chn <- input
        user_input[key] = chn
    }
    outputs["input"] = user_input

    /////////////////////////////////////////////////////
    fmt.Println("RUNNING...")

    for _, service := range flow {
        input := make(map[string]chan string)
        for key, value := range service.Inputs {
            input[key] = outputs[value["service"]][value["key"]]
        }

        outputs[service.Id] = runService(input, service)

        // output = runService(output, service.DockerImage)
        fmt.Println("ID:",service.Id)
    }
    fmt.Println(<- outputs["2"]["url"])

}

func runService(inputs map[string]chan string, service Service) map[string]chan string {
    result := make(map[string]chan string)
    for _, key := range service.Outputs {
      result[key] = make(chan string, 1)
    }
    go func() {
        fmt.Println(service.Id)
        args := <-inputs["url"]
        fmt.Println(service.Id, args)
        out, err := exec.Command("docker", "run", service.DockerImage, args).Output()

        fmt.Println(out)

        if err != nil {
            fmt.Println("ERROR:")
            fmt.Println(err)
            return
        }
        var dat map[string]string
        if err := json.Unmarshal([]byte(out), &dat); err != nil {
          panic(err)
        }
        for key, value := range dat {
          fmt.Println(service.Id, key, value)
          result[key] <- value
        }
    }()
    return result
}
