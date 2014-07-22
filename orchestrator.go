package main

import (
    // "fmt"
    "bytes"
    "os"
    "os/exec"
    "encoding/json"
    "log"
)

type Service struct {
    Id string
    DockerImage string
    Inputs map[string]map[string]string
    Outputs []string
}

// go run orchestrator.go "[{\"id\":\"1\",\"dockerImage\":\"cloudspace/url-lengthener-go\",\"inputs\":{\"url\":{\"service\":\"input\",\"key\":\"url\"}},\"outputs\":[\"url\"]},{\"id\":\"2\",\"dockerImage\":\"cloudspace/utm-stripper-go\",\"inputs\":{\"url\":{\"service\":\"1\",\"key\":\"url\"}},\"outputs\":[\"url\"]}]" "{\"url\":\"http://t.co/wnpJFiP7ls\"}"


func main() {
    var buf bytes.Buffer
    logger := log.New(&buf, "logger: ", log.Lshortfile)

    var flow []Service
    orchestrationSpecJSON := os.Args[1]
    logger.Println("ORCHESTRATION:", orchestrationSpecJSON)

    // parse the orchestration JSON
    orchestrationSpec := []byte(orchestrationSpecJSON)
    err1 := json.Unmarshal(orchestrationSpec, &flow)
    if err1 != nil {
        logger.Println(err1)
        return
    }
    logger.Printf("FLOW: %+v\n", flow)

    // parse user input
    logger.Println("PARSING OPTIONS JSON")

    var options map[string]string
    optionsJSON := os.Args[2]
    logger.Println("RAW OPTIONS:", optionsJSON)

    optionsBytes := []byte(optionsJSON)

    err2 := json.Unmarshal(optionsBytes, &options)
    if err2 != nil {
        logger.Println(err2)
        return
    }
    logger.Printf("PARSED OPTIONS: %+v\n", options)

    // prelaunch
    outputs := make(map[string]map[string]chan string)

    user_input := make(map[string]chan string)
    for key, input := range options {
        chn := make(chan string, 1)
        chn <- input
        user_input[key] = chn
    }
    outputs["input"] = user_input

    /////////////////////////////////////////////////////
    logger.Println("RUNNING...")

    for _, service := range flow {
        input := make(map[string]chan string)
        for key, value := range service.Inputs {
            input[key] = outputs[value["service"]][value["key"]]
        }

        outputs[service.Id] = runService(input, service, logger)

        logger.Println("ID:",service.Id)
    }

    // Output results from final service as json
    enc := json.NewEncoder(os.Stdout)
    for k, v := range outputs[flow[len(flow)-1].Id] {
      enc.Encode(map[string]string{k: <-v})
    }


    // // fmt.Println(k, <-v)
    // jsn, _ := json.Marshal([]string{k, <-v})
    // fmt.Println(string(jsn))
    // fmt.Println(<- outputs[flow[len(flow)-1].Id]["url"])
    // fmt.Print(&buf) // Uncomment to view logger output
}

func runService(inputs map[string]chan string, service Service, logger *log.Logger) map[string]chan string {
    result := make(map[string]chan string)
    for _, key := range service.Outputs {
      result[key] = make(chan string, 1)
    }
    go func() {
        logger.Println(service.Id)
        args := <-inputs["url"]
        logger.Println(service.Id, args)
        out, err := exec.Command("docker", "run", service.DockerImage, args).Output()
        if err != nil {
            logger.Println("ERROR: %s", err)
            return
        }
        logger.Println(out)

        var dat map[string]string
        if err := json.Unmarshal([]byte(out), &dat); err != nil {
          panic(err)
        }
        for key, value := range dat {
          logger.Println(service.Id, key, value)
          result[key] <- value
        }
    }()
    return result
}
