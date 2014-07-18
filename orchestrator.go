package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"
)

type Service struct {
    Name string
    DockerImage string
    Inputs map[string]string
    // Outputs map[string]string
}

// go run orchestrator.go "[{\"name\":\"service1\",\"DockerImage\":\"cloudspace/url-lengthener-go\",\"Inputs\":{\"url\":\"http://google.com\"},\"Outputs\":{\"url\":\"\"}},{\"name\":\"service1\",\"DockerImage\":\"cloudspace/url-lengthener-go\",\"Inputs\":{\"url\":\"http://google.com\"},\"Outputs\":{\"url\":\"\"}}]"

func main() {
    var flow []Service
    arguments := os.Args[1]
    fmt.Println("ARGUMENTS:", arguments)

    input := []byte(arguments)

    err := json.Unmarshal(input, &flow)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("%+v", flow)

    c := make(chan string)
    
    out := runService(c, flow[0].DockerImage)
    c <- flow[0].Inputs["url"]
    fmt.Printf("%s", <-out)
}

func runService(arg1 <-chan string, dockerImage string) <-chan string {
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
