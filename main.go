package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type AsciiData struct {
	Text    string
	Banner  string
	Output  string
	Message string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Starting server at localhost:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := AsciiData{}

	if r.Method == http.MethodPost {
		data.Text = r.FormValue("text")
		data.Banner = r.FormValue("banner")

		if !isAscii(data.Text) || data.Text == "" {
			data.Message = "Error: Please enter valid ASCII text."
		} else {
			output, err := generateAsciiArt(data.Text, data.Banner)
			if err != nil {
				data.Message = "Error: Unable to generate ASCII art. Please try again."
			} else {
				data.Output = output
				data.Message = "ASCII art generated successfully!"
			}
		}
	}

	tmpl.Execute(w, data)
}

// The rest of the functions (generateAsciiArt, printLine, isAscii) remain the same

func generateAsciiArt(text, banner string) (string, error) {
	filename := "ArtStyles/" + banner + ".txt"
	var outputBuffer strings.Builder

	strArr := strings.Split(text, "\n")
	for i := 0; i <= len(strArr)-1; i++ {
		if i-1 >= 0 {
			if strArr[i] == "" && i == len(strArr)-1 {
				continue
			}
		}
		if strArr[i] == "" {
			outputBuffer.WriteString("\n")
			continue
		}
		runes := []rune(strArr[i])
		for j := 0; j <= 8; j++ {
			for k := 0; k <= len(runes)-1; k++ {
				line := 2 + 9*(int(runes[k])-32) + j
				err := printLine(filename, line, &outputBuffer)
				if err != nil {
					return "", err
				}
			}
			if j < 8 {
				outputBuffer.WriteString("\n")
			}
		}
		if i < len(strArr)-1 {
			outputBuffer.WriteString("\n")
		}
	}

	return outputBuffer.String(), nil
}

func printLine(filename string, line int, output *strings.Builder) error {
	styleFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer styleFile.Close()

	styleFile.Seek(0, 0)

	scanner := bufio.NewScanner(styleFile)
	counter := 0
	for scanner.Scan() {
		counter++
		if counter == line {
			output.WriteString(scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func isAscii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
