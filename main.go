package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Job struct {
	Args     []string
	N, Total int
}

var (
	quiet       bool
	workerCount int
	csvFile     string
)

func producer(jobs [][]string) <-chan Job {
	ch := make(chan Job)
	go func() {
		for n, args := range jobs {
			ch <- Job{args, n, len(jobs)}
		}
		close(ch)
	}()
	return ch
}

func worker(id int, wg *sync.WaitGroup, cmd string, queue <-chan Job) {
	for job := range queue {
		if !quiet {
			fmt.Printf("[%d] %d/%d: %s %s\n", id, job.N, job.Total, cmd, strings.Join(job.Args, " "))
		}
		if exec.Command(cmd, job.Args...).Run() != nil && !quiet {
			fmt.Printf("[%d] %d/%d: %s %s failed to execute\n", id, job.N, job.Total, cmd, strings.Join(job.Args, " "))
		}
	}
	wg.Done()
}

func main() {
	flag.BoolVar(&quiet, "q", false, "Quiet")
	flag.IntVar(&workerCount, "n", 8, "Number of `workers`")
	flag.StringVar(&csvFile, "csv", "", "CSV file containing arguments for each job")
	flag.Parse()

	if workerCount < 1 || (csvFile == "" && flag.NArg() < 2) || (csvFile != "" && flag.NArg() < 1) {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0], "[-q] [-n <workers>] <cmd> <job1> [<job2> ...]")
		fmt.Println(os.Args[0], "[-q] [-n <workers>] -csv <csv-file> <cmd>")
		return
	}

	var jobs [][]string
	var err error
	cmd := flag.Arg(0)
	if csvFile == "" {
		for _, arg := range flag.Args()[1:] {
			jobs = append(jobs, []string{arg})
		}
	} else {
		var f io.Reader
		if csvFile == "-" {
			f = os.Stdin
		} else {
			f, err := os.Open(csvFile)
			if err != nil {
				fmt.Println("Failed to open", csvFile, err)
				return
			}
			defer f.Close()
		}
		jobs, err = csv.NewReader(f).ReadAll()
		if err != nil && err != io.EOF {
			fmt.Println("Failed to parse", csvFile, err)
			return
		}
	}

	queue := producer(jobs)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := 1; id <= workerCount; id++ {
		go worker(id, &wg, cmd, queue)
	}
	wg.Wait()
}
