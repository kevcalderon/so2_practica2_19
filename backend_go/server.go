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

// Estructura para almacenar la información de la RAM
type RAM struct {
	TotalRam string `json:"totalram"`
	FreeRam  string `json:"freeram"`
}
// Estructura para almacenar información de un segmento de memoria
type MemorySegment struct {
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
	Size         int    `json:"size_kb"`
	Permissions  string `json:"permissions"`
	Device       string `json:"device,omitempty"`
	FileName     string `json:"file_name,omitempty"`
}
// Estructura para almacenar la información de la memoria del sistema
type SystemMemory struct {
	TotalRAM int `json:"total_ram_mb"`
}
// Estructura para almacenar información simplificada de un segmento de memoria
type MemorySegment2 struct {
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
}
// Estructura para almacenar información de la memoria de un proceso
type ProcessMemory struct {
	ResidentMemory int `json:"resident_memory_mb"`
	VirtualMemory  int `json:"virtual_memory_mb"`
	RAMPercentage  float64 `json:"ram_percentage"`
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
// Manejador para el endpoint de obtener información de la RAM
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
// Manejador para el endpoint de obtener información de la CPU
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
// Manejador para el endpoint de obtener los segmentos de memoria de una carpeta específica
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
		// Leer el contenido del archivo en filePath
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Dividir el contenido en líneas individuales
	lines := strings.Split(string(content), "\n")
	memorySegments := make([]MemorySegment, 0)

	for _, line := range lines {
		// Dividir la línea en campos individuales
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			// Extraer el rango de direcciones
			addressRange := strings.Split(fields[0], "-")
			if len(addressRange) == 2 {
				startAddress := addressRange[0]
				endAddress := addressRange[1]

				// Calcular el tamaño del segmento de memoria
				size := calculateSegmentSize(startAddress, endAddress)

				// Crear una nueva estructura de MemorySegment
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

				// Agregar el segmento de memoria a la lista
				memorySegments = append(memorySegments, segment)
			}
		}
	}

	// Devolver la lista de segmentos de memoria y ningún error
	return memorySegments, nil
}

// calculateSegmentSize calcula el tamaño de un segmento de memoria en kilobytes
func calculateSegmentSize(startAddress, endAddress string) int {
	start, _ := strconv.ParseUint(startAddress, 16, 64)
	end, _ := strconv.ParseUint(endAddress, 16, 64)
	size := (end - start) / 1024 // Convertir de bytes a kilobytes

	return int(size)
}

// Manejador para el endpoint de matar un proceso por su ID
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

// Manejador para el endpoint de obtener información de la memoria de un proceso
func getProcessMemory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	processID := vars["id"]

	// Obtener el total de memoria RAM del sistema
	systemMemory, err := getSystemMemory()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al obtener el total de memoria RAM del sistema", http.StatusInternalServerError)
		return
	}

	// Construir la ruta del archivo smaps para el proceso dado
	filePath := fmt.Sprintf("/proc/%s/smaps", processID)
	// Leer el contenido del archivo smaps
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al leer el archivo smaps", http.StatusInternalServerError)
		return
	}

	// Dividir el contenido en líneas individuales
	lines := strings.Split(string(content), "\n")
	residentMemoryBytes := 0
	virtualMemoryBytes := 0

	// Recorrer cada línea del contenido
	for _, line := range lines {
		// Verificar si la línea comienza con "Rss:"
		if strings.HasPrefix(line, "Rss:") {
			// Dividir la línea en campos individuales
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Convertir el tamaño de memoria residente a bytes
				memorySizeBytes, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Error al convertir el tamaño de la memoria residente", http.StatusInternalServerError)
					return
				}
				residentMemoryBytes += memorySizeBytes
			}
		} else if strings.HasPrefix(line, "Size:") {
			// Verificar si la línea comienza con "Size:"
			// Dividir la línea en campos individuales
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Convertir el tamaño de la memoria virtual a bytes
				memorySizeBytes, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Error al convertir el tamaño de la memoria virtual", http.StatusInternalServerError)
					return
				}
				virtualMemoryBytes += memorySizeBytes
			}
		}
	}

	// Convertir los tamaños de memoria a MB
	residentMemoryMB := residentMemoryBytes / 1024
	virtualMemoryMB := virtualMemoryBytes / 1024
	// Calcular el porcentaje de uso de la RAM del proceso
	ramPercentage := float64(residentMemoryMB) / float64(systemMemory.TotalRAM) * 100

	// Crear una estructura ProcessMemory con la información del proceso
	processMemory := ProcessMemory{
		ResidentMemory:  residentMemoryMB,
		VirtualMemory:   virtualMemoryMB,
		RAMPercentage:   ramPercentage,
	}

	// Convertir la estructura ProcessMemory a formato JSON
	jsonData, err := json.Marshal(processMemory)
	if err != nil {
		http.Error(w, "Error al convertir el resultado a JSON", http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras adecuadas para la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// getSystemMemory obtiene la información sobre la memoria total del sistema
func getSystemMemory() (SystemMemory, error) {
	// Ejecutar el comando para obtener la memoria total del sistema
	cmd := exec.Command("sh", "-c", "free -m | awk 'NR==2{print $2}'")
	out, err := cmd.Output()
	if err != nil {
		return SystemMemory{}, err
	}

	// Convertir el resultado obtenido a un entero
	totalRAM, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return SystemMemory{}, err
	}

	// Crear una estructura SystemMemory con la información obtenida
	systemMemory := SystemMemory{
		TotalRAM: totalRAM,
	}

	return systemMemory, nil
}
