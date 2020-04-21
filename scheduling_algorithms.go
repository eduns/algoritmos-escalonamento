package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"strconv"
)

type Process struct {
	id             int
	arrivalTime    int
	burst          int
	remainingTime  int
	waitingTime    int
	responseTime   int
	turnaroundTime int
	quantum        int
	wait           int
}

func main() {
	for {
		fmt.Println("\n==========================")
		fmt.Println("1 - Escalonar os processos")
		fmt.Println("0 - Sair")
		fmt.Print("==========================\n> ")

		var option int
		fmt.Scanln(&option)

		if option == 1 {
			fmt.Println("\n============================")
			fmt.Println("Algoritmo de escalonamento:")
			fmt.Println("1 - FCFS")
			fmt.Println("2 - SJF")
			fmt.Println("3 - SRTF")
			fmt.Println("4 - Round Robin")
			fmt.Println("5 - Multinível")
			fmt.Print("============================\n> ")

			var alg int
			fmt.Scanln(&alg)

			for alg < 1 || alg > 5 {
				fmt.Print("\nOpção Inválida\n> ")
				fmt.Scanln(&alg)
			}

			switch alg {
				case 1:
					fcfs()

				case 2:
					sjf()

				case 3:
					srtf()

				case 4:
					rr()

				case 5:
					multilevel()
			}
		} else {
			break
		}
	}
}

func createProcesses() ([]Process, bool) {
	var processes []Process
	var q, op int
	var readFile bool
	var filename string

	fmt.Println("===================")
	fmt.Println("1 - Arquivo\n2 - Manual")
	fmt.Print("===================\n> ")
	fmt.Scanln(&op)

	if op == 1 {
		fmt.Print("Nome do arquivo > ")
		fmt.Scanln(&filename)
		
		file, err := os.Open(filename)

		if err != nil {
			log.Fatalf("Falha ao abrir o arquivo %s", err)
		} else {
			readFile = true
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			file.Close()

			for _, line := range lines {
				var x = strings.Split(line, ";")
				var id,_ = strconv.Atoi(strings.Split(x[0],"p")[1])
				var arrivalTime, _ = strconv.Atoi(x[2])
				var burst, _ = strconv.Atoi(x[1])
				var quantum, _ = strconv.Atoi(strings.Split(x[4],"&")[0])
				var process = Process{id, arrivalTime, burst, burst, 0, 0, 0, quantum, 0}
				processes = append(processes, process)
			}
		}
	} else {
		fmt.Print("Quantidade > ")
		fmt.Scanln(&q)

		readFile = false

		fmt.Println("======================")
		for i := 1; i <= q; i++ {
			var arrivalTime, burst int

			fmt.Printf("Tempo de chegada de P%d > ", i)
			fmt.Scanln(&arrivalTime)
			fmt.Printf("Burst do P%d > ", i)
			fmt.Scanln(&burst)
			fmt.Printf("\n")

			var process = Process{i, arrivalTime, burst, burst, 0, 0, 0, 0, 0}
			processes = append(processes, process)
		}
		fmt.Println("======================")
	}
	return processes, readFile
}

//============ ALGORITMOS DE ESCALONAMENTO ============
func fcfs() {
	var processes, _ = createProcesses()
	calc(processes)
	showStats(processes)
}

func sjf() {
	var currentProcess Process
	var totalExecutionTime int
	var pending, arrived, processed []Process

	var processes, _ = createProcesses()
	totalExecutionTime = calcExecutionTime(processes, 0)

	arrived = checkArrived(0, processes)
	sortProcessesBy(arrived, "burst")
	currentProcess = arrived[0]
	arrived = arrived[1:]
	if len(arrived) > 0 {
		pending = append(pending, arrived...)
	}

	for i := 1; i <= totalExecutionTime; i++ {
		arrived = checkArrived(i, processes)

		if len(arrived) > 0 {
			pending = append(pending, arrived...)
			sortProcessesBy(pending, "burst")
		}

		if currentProcess.remainingTime == 0 {
			processed = append(processed, currentProcess)
			if len(pending) > 0 {
				currentProcess = pending[0]
				pending = pending[1:]
			}
		}
		currentProcess.remainingTime -= 1
	}

	processed = append(processed, currentProcess)
	calc(processed)
	sortProcessesBy(processed, "id")
	showStats(processed)
}

func srtf() {
	var currentProcess Process
	var totalExecutionTime int
	var arrived, pending, processed []Process

	var processes, _ = createProcesses()
	totalExecutionTime = calcExecutionTime(processes, 0)
	currentProcess = processes[0]
	processes = processes[1:]
	
	for j := 0; j <= totalExecutionTime; j++ {
		arrived = checkArrived(j, processes)

		if len(arrived) > 0 {
			sortProcessesBy(arrived, "remainingTime")
			pending = append(pending, arrived...)
			sortProcessesBy(pending, "remainingTime")
		}

		if currentProcess.remainingTime > 0 {
			if len(pending) > 0 && pending[0].remainingTime < currentProcess.remainingTime {
				currentProcess.waitingTime += currentProcess.wait
				currentProcess.wait = 0	

				pending = append(pending, currentProcess)
				currentProcess = pending[0]
				pending = pending[1:]
				sortProcessesBy(pending, "remainingTime")
			}
		} else {
			currentProcess.turnaroundTime = j - currentProcess.arrivalTime;
			currentProcess.waitingTime += currentProcess.wait
			currentProcess.wait = 0

			processed = append(processed, currentProcess);
			if len(pending) > 0 {
				currentProcess = pending[0]
				pending = pending[1:]
			}
		}
		for k := 0; k < len(pending); k++ {
			pending[k].wait += 1
		}
		currentProcess.remainingTime -= 1
		if currentProcess.remainingTime + 1  == currentProcess.burst {
			currentProcess.responseTime = j - currentProcess.arrivalTime
		}
	}
	sortProcessesBy(processed, "id")
	showStats(processed)
}

func rr() {
	var q, quantum, totalExecutionTime int
	var currentProcess Process
	var arrived, pending, processed []Process
	var processes, readFile = createProcesses()

	if readFile {
		q = processes[0].quantum
	} else  {
		fmt.Print("Quantum > ")
		fmt.Scanln(&q)
	}

	quantum = q
	totalExecutionTime = calcExecutionTime(processes, 0)
	currentProcess = processes[0]
	processes = processes[1:]
	
	for i := 0; i <= totalExecutionTime; i++ {
		arrived = checkArrived(i, processes)

		if len(arrived) > 0 {
			pending = append(pending, arrived...)
		}

		if currentProcess.remainingTime == 0 {
			quantum = q
			currentProcess.waitingTime += currentProcess.wait
			currentProcess.wait = 0

			currentProcess.turnaroundTime = i - currentProcess.arrivalTime
			processed = append(processed, currentProcess)

			if len(pending) > 0 {
				currentProcess = pending[0]
				pending = pending[1:]
			}
			
		} else if quantum == 0 {
			quantum = q
			currentProcess.waitingTime += currentProcess.wait
			currentProcess.wait = 0

			pending = append(pending, currentProcess)

			if len(pending) > 0 {
				currentProcess = pending[0]
				pending = pending[1:]
			}
		}
		quantum -= 1
		currentProcess.remainingTime -= 1

		for j := 0; j < len(pending); j++ {
			pending[j].wait += 1
		}

		if currentProcess.remainingTime + 1  == currentProcess.burst {
			currentProcess.responseTime = i - currentProcess.arrivalTime
		}
	}
	sortProcessesBy(processed, "id")
	showStats(processed)
}

func multilevel() {
	var q, quantum, totalExecutionTime int
	var currentProcess Process
	var processed, fcfsQueue []Process
	var processes, readFile = createProcesses()

	if readFile {
		q = processes[0].quantum
	} else {
		fmt.Print("Quantum > ")
		fmt.Scanln(&q)
	}
	quantum = q
	totalExecutionTime = calcExecutionTime(processes, quantum)
	currentProcess = processes[0]
	processes = processes[1:]
	
	for i := 0; i <= totalExecutionTime; i++ {
		for k := 0; k < len(fcfsQueue); k++ {
			fcfsQueue[k].wait += 1
		}

		if currentProcess.remainingTime == 0 {
			quantum = q
			currentProcess.waitingTime += currentProcess.wait - currentProcess.arrivalTime
			currentProcess.wait = 0
			currentProcess.turnaroundTime = i - currentProcess.arrivalTime
			processed = append(processed, currentProcess)

			if len(processes) > 0 {
				currentProcess = processes[0]
				processes = processes[1:]
			}
		} else if quantum == 0 {
			quantum = q

			currentProcess.turnaroundTime = i
			currentProcess.waitingTime += currentProcess.wait
			currentProcess.wait = 0

			fcfsQueue = append(fcfsQueue, currentProcess)
			if len(processes) > 0 {
				currentProcess = processes[0]
				processes = processes[1:]
			}
		}
		quantum -= 1
		currentProcess.remainingTime -= 1

		if currentProcess.remainingTime + 1  == currentProcess.burst {
			currentProcess.responseTime = i - currentProcess.arrivalTime
		}

		for j := 0; j < len(processes); j++ {
			processes[j].wait += 1
		}
	}

	for l := 0; l < len(fcfsQueue); l++ {
		fcfsQueue[l].waitingTime += fcfsQueue[l].wait
		fcfsQueue[l].turnaroundTime += fcfsQueue[l].wait + fcfsQueue[l].remainingTime
		fcfsQueue[l].wait = 0
	}
	processed = append(processed, fcfsQueue...)
	sortProcessesBy(processed, "id")
	showStats(processed)
}

func calcExecutionTime(processes []Process, quantum int) int {
	var totalExecutionTime int
	if quantum > 0 {
		for i := 0; i < len(processes); i++ {
			if processes[i].burst > quantum {
				totalExecutionTime += processes[i].burst - (processes[i].burst - quantum)
			} else {
				totalExecutionTime += processes[i].burst
			}
		}
	} else {
		for i := 0; i < len(processes); i++ {
			totalExecutionTime += processes[i].burst
		}
	}
	return totalExecutionTime
}

func sortProcessesBy(processes []Process, attr string) {
	if attr == "id" {
		sort.Slice(processes, func (i, j int) bool {
			return processes[i].id < processes[j].id
		})
	} else if attr == "arrivalTime" {
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].arrivalTime < processes[j].arrivalTime
		})
	} else if attr == "burst" {
		sort.Slice(processes, func (i, j int) bool {
			return processes[i].burst < processes[j].burst
		})
	} else if attr == "remainingTime" {
		sort.Slice(processes, func (i, j int) bool {
			return processes[i].remainingTime < processes[j].remainingTime
		})
	} 
}

func checkArrived(executionTime int, processes []Process) []Process {
	var arrived []Process
	for i := 0; i < len(processes); i++ {
		if processes[i].arrivalTime == executionTime {
			arrived = append(arrived, processes[i])
		}
	}
	return arrived
}

func calc(processes []Process) {
	var processesSize = len(processes)
	// =========== WAITING TIME ============
	for i := 0; i < processesSize; i++ {
		for j := 0; j < i; j++ {
			processes[i].waitingTime += processes[j].burst
		}
		processes[i].waitingTime -= processes[i].arrivalTime
	}
	// ============ TURNAROUND TIME ===========
	for k := 0; k < processesSize; k++ {
		processes[k].turnaroundTime += processes[k].burst + processes[k].waitingTime
	}
	// ============== RESPONSE TIME =====================
	for l := 1; l < processesSize; l++ {
		for m := 1; m <= l; m++ {
			processes[l].responseTime += processes[m-1].burst
		}
		processes[l].responseTime -= processes[l].arrivalTime
	}
}

func showStats(processes []Process) {
	var avgWaitingTime, avgTurnaroundTime float32
	var qtyprocesses int = len(processes)

	fmt.Println("-------------------------------------------------")
	fmt.Println(" P | Burst | Arrival | Turnaround | Waiting Time")
	fmt.Println("-------------------------------------------------")
	for i := 0; i < qtyprocesses; i++ {
			fmt.Printf(" %d | %5d | %7d | %10d | %12d\n",
				processes[i].id, processes[i].burst,
				processes[i].arrivalTime, processes[i].turnaroundTime,
				processes[i].waitingTime)

		avgTurnaroundTime += float32(processes[i].turnaroundTime)
		avgWaitingTime += float32(processes[i].waitingTime)
	}
	avgWaitingTime /= float32(qtyprocesses)
	avgTurnaroundTime /= float32(qtyprocesses)
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Average Waiting Time: %.2f\n", avgWaitingTime)
	fmt.Printf("Average Turnaround Time: %.2f\n", avgTurnaroundTime)
	fmt.Println("-------------------------------------------------")
}