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
    working_directory: /go/src/github.com/ironarachne/{{.Name}}
    steps:
      - checkout
      - run: go get -v -t -d ./...
      - run: go test -v ./...
      - setup_remote_docker:
          docker_layer_caching: true
      - run: |
          cd cmd/{{.Name}}d
          TAG=0.1.$CIRCLE_BUILD_NUM
          docker build -t ironarachne/{{.Name}}d:$TAG -t ironarachne/{{.Name}}d:latest .
          docker login -u $DOCKER_USER -p $DOCKER_PASS
          docker push ironarachne/{{.Name}}d:$TAG
          docker push ironarachne/{{.Name}}d:latest
  deploy:
    machine:
        enabled: true
    steps:
      - run: curl -X POST 'https://portainer.ironarachne.com/api/webhooks/'

workflows:
  version: 2
  build-and-deploy:
    jobs:
      - build
      - deploy:
          requires:
            - build
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
  "math/rand"

  "github.com/ironarachne/{{.Name}}"
  "github.com/ironarachne/utility"
  "github.com/kataras/iris"
)

func main() {
  app := iris.New()

  app.Get("/", func(ctx iris.Context) {
    ctx.Writef("{{.Name}}d")
  })

  app.Get("/{id:int64}", func(ctx iris.Context) {
    id, err := ctx.Params().GetInt64("id")
    if err != nil {
      ctx.Writef("error while trying to parse id parameter")
      ctx.StatusCode(iris.StatusBadRequest)
      return
    }

    rand.Seed(id)
    {{.Object}} := {{.Name}}.Generate()

    ctx.JSON({{.Object}})
  })

  app.Run(iris.Addr(":{{.Port}}"))
}`

	programFileTemplate = `package {{.Name}}

import (
  "math/rand"
  "strings"

  "github.com/ironarachne/utility"
)

type {{.Object}} struct {
}

// Generate generates a {{.Object}}
func Generate() {
  {{.Object}} := {{.Object}}{}

  return {{.Object}}
}`

	readmeFileTemplate = `# {{.Name}}

Just another generator. This one's for {{.Object}}s.`
)
