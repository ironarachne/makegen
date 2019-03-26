package main

var (
	buildFileTemplate = `#!/bin/bash
go build -o build/{{.Name}}d cmd/{{.Name}}d/main.go`

	circleFileTemplate = `# Golang CircleCI 2.0 configuration file
#
# Check https://circleci.com/docs/2.0/language-go/ for more details
version: 2
jobs:
  build:
    docker:
      - image: circleci/golang:1.11
    working_directory: /go/src/github.com/{{.Organization}}/{{.Name}}
    steps:
      - checkout
      - run: go get -v -t -d ./...
      - run: go test -v ./...
      - setup_remote_docker:
          docker_layer_caching: true
      - run: |
          cd cmd/{{.Name}}d
          TAG=0.1.$CIRCLE_BUILD_NUM
          docker build -t {{.Organization}}/{{.Name}}d:$TAG -t {{.Organization}}/{{.Name}}d:latest .
          docker login -u $DOCKER_USER -p $DOCKER_PASS
          docker push {{.Organization}}/{{.Name}}d:$TAG
          docker push {{.Organization}}/{{.Name}}d:latest

workflows:
  version: 2
  build-and-publish:
    jobs:
      - build:
          filters:
            branches:
              only: master`

	dockerFileTemplate = `# build stage
FROM golang:1.11 AS build-env
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/{{.Name}}d

# final stage
FROM scratch
COPY --from=build-env /go/bin/{{.Name}}d /go/bin/{{.Name}}d
EXPOSE {{.Port}}
CMD ["/go/bin/{{.Name}}d"]`

	gitignoreFileTemplate = `build/`

	mainFileTemplate = `package main

import (
	"log"
  "fmt"
  "math/rand"
  "net/http"
  "time"

  "github.com/{{.Organization}}/{{.Name}}"
  "github.com/ironarachne/random"
  "github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func get{{.TitleObject}}(w http.ResponseWriter, r *http.Request) {
  id := chi.URLParam(r, "id")

	var new{{.TitleObject}} {{.Name}}.{{.TitleObject}}

	random.SeedFromString(id)

	new{{.TitleObject}} = {{.Name}}.Generate{{.TitleObject}}()

	json.NewEncoder(w).Encode(new{{.TitleObject}})
}

func main() {
  r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("{\"status\": \"online\"}"))
  })

  r.Get("/{id}", get{{.TitleObject}})

  fmt.Println("{{.TitleObject}} Generator API is online.")
	log.Fatal(http.ListenAndServe(":{{.Port}}", r))
}`

	programFileTemplate = `package {{.Name}}

import (
  "math/rand"
  "strings"

  "github.com/ironarachne/random"
)

// {{.TitleObject}} is a {{.Object}}
type {{.TitleObject}} struct {
}

// Generate generates a {{.Object}}
func Generate() {
  {{.Object}} := {{.TitleObject}}{}

  return {{.Object}}
}`

	readmeFileTemplate = `# {{.Name}}

Just another generator. This one's for {{.Object}}s.`
)
