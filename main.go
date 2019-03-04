package main

import (
	"flag"
	"os"
	"text/template"
)

// Generator is a generator
type Generator struct {
	Name string
	Port int
}

func (gen Generator) createDirectories() {
	err := os.MkdirAll(gen.Name+"/cmd/"+gen.Name, 0755)
	if err != nil {
		panic(err)
	}

	os.Mkdir(gen.Name+"/.circleci", 0755)
}

func (gen Generator) createFiles() {
	gen.renderToFile(gen.Name+"/cmd/"+gen.Name+"/Dockerfile", dockerFileTemplate)
	gen.renderToFile(gen.Name+"/.circleci/config.yml", circleFileTemplate)
	gen.renderToFile(gen.Name+"/build.sh", buildFileTemplate)
	gen.renderToFile(gen.Name+"/cmd/"+gen.Name+"/main.go", mainFileTemplate)
	gen.renderToFile(gen.Name+"/"+gen.Name+".go", programFileTemplate)
	gen.renderToFile(gen.Name+"/README.md", readmeFileTemplate)
	gen.renderToFile(gen.Name+"/.gitignore", gitignoreFileTemplate)

	buildFile, _ := os.Open(gen.Name + "/build.sh")
	defer buildFile.Close()
	buildFile.Chmod(0744)
}

func (gen Generator) renderToFile(filePath string, templateData string) {
	tmpl, err := template.New("template").Parse(templateData)
	if err != nil {
		panic(err)
	}
	fp, _ := os.Create(filePath)
	defer fp.Close()

	tmpl.Execute(fp, gen)
}

func touch(fileName string) {
	f, _ := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	defer f.Close()
}

func main() {
	var genName = flag.String("n", "", "The name of the generator")
	var port = flag.Int("p", 0, "The port for the generator's API")

	flag.Parse()

	generator := Generator{*genName, *port}

	generator.createDirectories()
	generator.createFiles()
}
