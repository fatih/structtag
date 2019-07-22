workflow "Build and Test" {
  on = "pull_request"
  resolves = ["go test"]
}

action "go build" {
  uses = "docker://golang:1.12.7"
  runs = "go"
  args = "build ./..."
  env = {
    GOPROXY = "https://proxy.golang.org"
  }
}

action "go test" {
  uses = "docker://golang:1.12.7"
  needs = ["go build"]
  runs = "go"
  args = "test -v -cover -race ./..."
  env = {
    GOPROXY = "https://proxy.golang.org"
  }
}
