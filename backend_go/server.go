package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


type RAM struct {
	TotalRam string `json:"totalram"`
	FreeRam  string `json:"freeram"`
}

type MemorySegment struct {
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
	Size         int    `json:"size_kb"`
	Permissions  string `json:"permissions"`
	Device       string `json:"device,omitempty"`
	FileName     string `json:"file_name,omitempty"`
}

type SystemMemory struct {
	TotalRAM int `json:"total_ram_mb"`
}

type MemorySegment2 struct {
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
}

type ProcessMemory struct {
	ResidentMemory int `json:"resident_memory_mb"`
	VirtualMemory  int `json:"virtual_memory_mb"`
	RAMPercentage  float64 `json:"ram_percentage"`
	MemorySegments []MemorySegment2  `json:"memory_segments"`
}


func main() {
	router := mux.NewRouter()

	// Endpoint para obtener información de la RAM
	router.HandleFunc("/api/ram", handleRequest).Methods("GET")

	// Endpoint para obtener información de la CPU
	router.HandleFunc("/api/cpu", handleCPURequest).Methods("GET")

	// Endpoint para obtener los segmentos de memoria de una carpeta específica
	router.HandleFunc("/api/memoria/{folder}", getMemorySegments).Methods("GET")

	// Endpoint para matar un proceso por su ID
	router.HandleFunc("/api/kill/{id}", handleKill).Methods("GET")

	router.HandleFunc("/api/memoryprocess/{id}", getProcessMemory).Methods("GET")

	fmt.Println("Servidor en ejecución en http://localhost:8080")
	// Agregar el middleware CORS a todos los endpoints
	handler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	http.ListenAndServe(":8080", handler)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Ejecutar el comando para obtener la información de la RAM
	cmd := exec.Command("sh", "-c", "cat /proc/ram_grupo19")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al ejecutar el comando", http.StatusInternalServerError)
		return
	}

	outputRam := string(out[:])

	var ram RAM
	err = json.Unmarshal([]byte(outputRam), &ram)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al analizar el JSON", http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras adecuadas para la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Devolver el JSON generado
	fmt.Fprintf(w, "%s", outputRam)
}

func handleCPURequest(w http.ResponseWriter, r *http.Request) {
	// Ejecutar el comando para obtener la información de la CPU
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_grupo19")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al ejecutar el comando", http.StatusInternalServerError)
		return
	}

	outputCpu := string(out[:])

	// Establecer las cabeceras adecuadas para la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Devolver el valor de outputCpu
	fmt.Fprintf(w, "%s", outputCpu)
}

func getMemorySegments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Obtener el parámetro de la URL para buscar informacion del proceso
	folder := vars["folder"]

	filePath := "/proc/" + folder + "/maps"
	memorySegments, err := readMemorySegments(filePath)
	if err != nil {
		http.Error(w, "Error al leer los segmentos de memoria", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(memorySegments)
	if err != nil {
		http.Error(w, "Error al convertir los segmentos de memoria a JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// readMemorySegments lee los segmentos de memoria desde el archivo maps y los devuelve como una lista de MemorySegment
func readMemorySegments(filePath string) ([]MemorySegment, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	memorySegments := make([]MemorySegment, 0)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			addressRange := strings.Split(fields[0], "-")
			if len(addressRange) == 2 {
				startAddress := addressRange[0]
				endAddress := addressRange[1]

				size := calculateSegmentSize(startAddress, endAddress)

				segment := MemorySegment{
					StartAddress: startAddress,
					EndAddress:   endAddress,
					Size:         size,
				}

				// Set Permissions field
				if len(fields) >= 5 {
					segment.Permissions = fields[1]
				}

				// Set Device field if applicable
				if len(fields) >= 6 {
					segment.Device = fields[5]
				}

				// Set FileName field if applicable
				if len(fields) >= 6 {
					segment.FileName = fields[5]
				}

				memorySegments = append(memorySegments, segment)
			}
		}
	}

	return memorySegments, nil
}

// calculateSegmentSize calcula el tamaño de un segmento de memoria en kilobytes
func calculateSegmentSize(startAddress, endAddress string) int {
	start, _ := strconv.ParseUint(startAddress, 16, 64)
	end, _ := strconv.ParseUint(endAddress, 16, 64)
	size := (end - start) / 1024 // Convertir de bytes a kilobytes

	return int(size)
}

func handleKill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]

	// Ejecutar el comando "kill {processID}" para terminar el proceso	
	cmd := exec.Command("sh", "-c", fmt.Sprintf("kill %s", processID))

	err := cmd.Run()
	if err != nil {
		http.Error(w, "Error al matar el proceso", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Proceso %s cerrado", processID)
}

func getProcessMemory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]

	systemMemory, err := getSystemMemory()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al obtener el total de memoria RAM del sistema", http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("/proc/%s/smaps", processID)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al leer el archivo smaps", http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(content), "\n")
	residentMemoryBytes := 0
	virtualMemoryBytes := 0
	var memorySegments []MemorySegment2
	var currentSegment MemorySegment2

	for _, line := range lines {
		if strings.HasPrefix(line, "Rss:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memorySizeBytes, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Error al convertir el tamaño de la memoria residente", http.StatusInternalServerError)
					return
				}
				residentMemoryBytes += memorySizeBytes
			}
		} else if strings.HasPrefix(line, "Size:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memorySizeBytes, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Error al convertir el tamaño de la memoria virtual", http.StatusInternalServerError)
					return
				}
				virtualMemoryBytes += memorySizeBytes
			}
		} else if strings.HasPrefix(line, "Start:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				startAddress := fields[1]
				currentSegment = MemorySegment2{
					StartAddress: startAddress,
				}
			}
		} else if strings.HasPrefix(line, "End:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				endAddress := fields[1]
				currentSegment.EndAddress = endAddress
				memorySegments = append(memorySegments, currentSegment)
			}
		}
	}

	residentMemoryMB := residentMemoryBytes / 1024
	virtualMemoryMB := virtualMemoryBytes / 1024
	ramPercentage := float64(residentMemoryMB) / float64(systemMemory.TotalRAM) * 100

	processMemory := ProcessMemory{
		ResidentMemory:  residentMemoryMB,
		VirtualMemory:   virtualMemoryMB,
		RAMPercentage:   ramPercentage,
		MemorySegments:  memorySegments,
	}

	jsonData, err := json.Marshal(processMemory)
	if err != nil {
		http.Error(w, "Error al convertir el resultado a JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getSystemMemory() (SystemMemory, error) {
	cmd := exec.Command("sh", "-c", "free -m | awk 'NR==2{print $2}'")
	out, err := cmd.Output()
	if err != nil {
		return SystemMemory{}, err
	}

	totalRAM, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return SystemMemory{}, err
	}

	systemMemory := SystemMemory{
		TotalRAM: totalRAM,
	}

	return systemMemory, nil
}