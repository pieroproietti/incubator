package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pieroproietti/incubator/pkg/config"
	"github.com/pieroproietti/incubator/pkg/orchestrator"
)

func main() {
	targetDir := "/var/lib/vz/template/iso"
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	cfg := config.Config{
		TargetDir:  targetDir,
		BaseVMID:   config.GetEnvAsInt("VMID", 101),
		Firmware:   config.GetEnv("FIRMWARE", "bios"),
		FsType:     config.GetEnv("FSTYPE", "ext4"),
		Storage:    config.GetEnv("STORAGE", "father-zfs"),       // Ripristinato!
		IsoStorage: config.GetEnv("ISO_STORAGE", "father-local"), // Ripristinato!
		Template:   config.GetEnv("TEMPLATE", ""),
		Bridge:     config.GetEnv("BRIDGE", "vmbr0"), // Ripristinato!
	}

	isos, err := filepath.Glob(filepath.Join(cfg.TargetDir, "*.iso"))
	if err != nil || len(isos) == 0 {
		log.Fatalf("[ERRORE] Nessuna ISO trovata nella directory: %s\n", cfg.TargetDir)
	}

	fmt.Println("===================================================================")
	fmt.Printf("🐣 PENGUINS INCUBATOR (Go Edition - Modular & Reporting)\n")
	fmt.Printf("Batch iniziato: trovate %d ISO in %s\n", len(isos), cfg.TargetDir)
	fmt.Printf("Configurazione -> Firmware: %s | FsType: %s | Base VMID: %d\n", cfg.Firmware, cfg.FsType, cfg.BaseVMID)
	fmt.Println("===================================================================")

	const numWorkers = 3
	tasks := make(chan orchestrator.IsoTask, len(isos))

	// ECCO IL TUBO MANCANTE: Creiamo il canale per raccogliere i referti
	reportsChan := make(chan orchestrator.TestReport, len(isos))
	var wg sync.WaitGroup

	// Innesca i muratori passando reportsChan come quarto argomento
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerVMID := cfg.BaseVMID + i
		go orchestrator.Worker(workerVMID, cfg, tasks, reportsChan, &wg)
	}

	for _, isoFullPath := range isos {
		tasks <- orchestrator.IsoTask{
			IsoName: filepath.Base(isoFullPath),
			IsoPath: isoFullPath,
		}
	}
	close(tasks)

	wg.Wait()
	close(reportsChan) // Chiudiamo il canale dei referti

	fmt.Println("\n>>> BATCH COMPLETATO. Generazione del report finale in corso...")
	generateMarkdownSummary(reportsChan, cfg)
}

func generateMarkdownSummary(reports <-chan orchestrator.TestReport, cfg config.Config) {
	fileName := fmt.Sprintf("incubator-summary-%s-%s.md", cfg.Firmware, cfg.FsType)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Impossibile creare il file di report: %v", err)
		return
	}
	defer file.Close()

	header := fmt.Sprintf("## Risultati CI Incubator (%s / %s)\n\n", strings.ToUpper(cfg.Firmware), strings.ToUpper(cfg.FsType))
	header += "| Distro | ISO | Firmware | FileSystem | Status | Durata |\n"
	header += "|--------|-----|----------|------------|--------|--------|\n"

	file.WriteString(header)
	fmt.Print("\n" + header)

	for r := range reports {
		durata := r.Duration.Round(time.Second)
		// Convertiamo il nome della distro in MAIUSCOLO per staccare visivamente
		row := fmt.Sprintf("| **%s** | %s | %s | %s | %s | %s |\n",
			strings.ToUpper(r.Distro), r.IsoName, r.Firmware, r.FsType, r.Status, durata)

		file.WriteString(row)
		fmt.Print(row)
	}

	fmt.Printf("\nReport salvato con successo in: %s\n", fileName)
}
