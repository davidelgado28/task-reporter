package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"
)

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
}

const reportTemplate = `# Relatório de Tarefas
Gerado em: {{ .Date }}

| ID | Tarefa | Status | Responsável |
| :---: | :--- | :---: | :--- |
{{- range .Tasks }}
| {{ .ID }} | {{ .Title }} | {{ .Status }} | {{ .Assignee }} |
{{- end }}
`

func main() {
	inputPath := flag.String("input", "data/tasks.json", "Caminho para o arquivo JSON de tarefas")
	outputPath := flag.String("output", "report.md", "Caminho para salvar o relatório em Markdown")
	flag.Parse()

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao abrir arquivo de entrada: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao decodificar JSON: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar template: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo de saída: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	data := struct {
		Date  string
		Tasks []Task
	}{
		Date:  time.Now().Format("02/01/2006 15:04:05"),
		Tasks: tasks,
	}

	if err := tmpl.Execute(outFile, data); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Relatório gerado com sucesso em %s\n", *outputPath)
}
