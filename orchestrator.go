package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"
)

// go run orchestrator.go "{\"DockerImage\":\"cloudspace/url-lengthener-go\",\"Arguments\":\"http://google.com\"}"

// type Service struct {
//     DockerImage string
//     Arguments string
// }

// func main() {
//     var service Service
//     arguments := os.Args[1]
//     // fmt.Println("ARGUMENTS:", arguments)

//     input := []byte(arguments)

//     err := json.Unmarshal(input, &service)
//     if err != nil {
//         fmt.Println(err)
//         return
//     }

//     out, err := exec.Command("docker", "run", service.DockerImage, service.Arguments).Output()

//     if err != nil {
//         fmt.Println("ERROR:")
//         fmt.Println(err)
//         return
//     }

//     fmt.Printf("%s", out)
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Service struct {
    Name string
    DockerImage string
    Inputs map[string]string
    Outputs map[string]string
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

    channel := make(chan map[string]string)

    out, err := exec.Command("docker", "run", flow[0].DockerImage, flow[0].Inputs["url"]).Output()
    
    if err != nil {
        fmt.Println("ERROR:")
        fmt.Println(err)
        return
    }

    fmt.Printf("%s", out)
}

func runService(service Service, c chan) {

}
