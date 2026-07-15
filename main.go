package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Config rappresenta le variabili d'ambiente e i parametri (mantiene l'interfaccia Bash)
type Config struct {
	TargetDir  string
	BaseVMID   int
	Firmware   string
	FsType     string
	Storage    string
	IsoStorage string
}

// IsoTask rappresenta il singolo lavoro da eseguire su una ISO
type IsoTask struct {
	IsoName string
	VMID    int
}

func main() {
	// 1. Lettura dell'Interfaccia (Argomenti e Variabili d'Ambiente)
	if len(os.Args) < 2 {
		log.Fatalf("Uso: %s /percorso/della/directory/iso\n", os.Args[0])
	}
	
	// Valori di default o presi dall'ambiente, esattamente come in Bash
	config := Config{
		TargetDir:  os.Args[1],
		BaseVMID:   getEnvAsInt("VMID", 150),
		Firmware:   getEnv("FIRMWARE", "bios"),
		FsType:     getEnv("FSTYPE", "ext4"),
		Storage:    getEnv("STORAGE", "father-zfs"),
		IsoStorage: getEnv("ISO_STORAGE", "father-local"),
	}

	// 2. Troviamo tutte le ISO nella cartella
	isos, err := filepath.Glob(filepath.Join(config.TargetDir, "*.iso"))
	if err != nil || len(isos) == 0 {
		log.Fatalf("Nessuna ISO trovata in %s", config.TargetDir)
	}
	fmt.Printf("Batch iniziato: trovate %d ISO in %s\n", len(isos), config.TargetDir)

	// 3. Setup del Worker Pool (Il Semaforo)
	// Qui decidiamo quante VM lanciare contemporaneamente (es. 3)
	const numWorkers = 3
	tasks := make(chan IsoTask, len(isos))
	var wg sync.WaitGroup

	// 4. Lanciamo i "Muratori" (Goroutines)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// Ogni worker avrà il suo offset per il VMID, per non accavallarsi
		workerVMID := config.BaseVMID + i 
		go worker(workerVMID, config, tasks, &wg)
	}

	// 5. Inseriamo i lavori nel tubo (Channel)
	for _, isoPath := range isos {
		tasks <- IsoTask{
			IsoName: filepath.Base(isoPath),
		}
	}
	close(tasks) // Chiudiamo il tubo: non ci sono più lavori

	// 6. Aspettiamo che tutti abbiano finito
	wg.Wait()
	fmt.Println("=== TUTTI I COLLAUDI COMPLETATI ===")
}

// worker è la funzione eseguita in parallelo. Pesca dal canale finché ci sono ISO.
func worker(workerVMID int, cfg Config, tasks <-chan IsoTask, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		task.VMID = workerVMID
		fmt.Printf("[Worker VMID:%d] Inizio collaudo ISO: %s\n", task.VMID, task.IsoName)
		
		// Qui mettiamo tutta la logica del test (creazione VM, accensione, qm agent, ecc.)
		success := runIncubatorTest(task, cfg)
		
		if success {
			fmt.Printf("[Worker VMID:%d] SUCCESS: %s\n", task.VMID, task.IsoName)
		} else {
			fmt.Printf("[Worker VMID:%d] FAILED: %s\n", task.VMID, task.IsoName)
		}
	}
}

// runIncubatorTest conterrà la traduzione dei comandi qm
func runIncubatorTest(task IsoTask, cfg Config) bool {
	// Simulazione della pulizia vecchia VM
	// execQmCommand("stop", fmt.Sprintf("%d", task.VMID))
	// execQmCommand("destroy", fmt.Sprintf("%d", task.VMID))
	
	// Simulazione creazione e test...
	time.Sleep(3 * time.Second) // Finge di lavorare
	
	return true 
}

// Funzioni di utilità per leggere l'ambiente
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	// Omettiamo per brevità la logica di strconv.Atoi
	return fallback
}

