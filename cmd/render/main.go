package main

import (
	"flag"
	"log"
	"os"
	"text/template"
	"time"
	"unicode"

	"demo.test/grpc-demo/internal/client/models"
	"gopkg.in/yaml.v2"
)

var (
	renderValYmlPath string
	templPath        string
)

func initTempRenderFlagSet() {
	renderFlagSet := flag.NewFlagSet("render_flagset", flag.ExitOnError)
	renderFlagSet.StringVar(&renderValYmlPath, "renderval_path", "", "the path of render values")
	renderFlagSet.StringVar(&templPath, "temp_path", "", "the path of template file")

	// logrus.Infoln(os.Args)
	if !renderFlagSet.Parsed() {
		renderFlagSet.Parse(os.Args[1:])
	}
}

func parse(destination *os.File) error {
	templVals := models.Templ{}
	yamlData, err := os.ReadFile(renderValYmlPath)
	if err != nil {
		log.Fatalf("read file failed, error: %v", err)
		return err
	}
	err = yaml.Unmarshal(yamlData, &templVals)
	if err != nil {
		log.Fatalf("parse yaml failed, error: %v", err)
		return err
	}

	data, err := os.ReadFile(templPath)
	if err != nil {
		log.Fatalf("read template file failed, error: %v", err)
		return err
	}
	// create a template instance
	tmpl, err := template.New("template").
		Funcs(template.FuncMap{
			"firstCharLower": firstCharLower,
			"currentTime":    currentTime,
		}).
		Parse(string(data))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// render template
	err = tmpl.Execute(destination, templVals)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// converts string first letter to lowercase
func firstCharLower(str string) string {
	if len(str) > 0 {
		firstChar := []rune(str)[0]
		if unicode.IsUpper(firstChar) {
			str = string(unicode.ToLower(firstChar)) + str[1:]
		}
	}
	return str
}

func currentTime() string {
	currentTime := time.Now().UTC()
	return currentTime.Format("2006-01-02 15:04:05")
}

/*----------------------------------------------------------------------------------------*/
func main() {

	initTempRenderFlagSet()
	parse(os.Stdout)

}
