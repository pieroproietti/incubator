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
		Storage:    config.GetEnv("STORAGE", "father-zfs"),       // Restored!
		IsoStorage: config.GetEnv("ISO_STORAGE", "father-local"), // Restored!
		Template:   config.GetEnv("TEMPLATE", ""),
		Bridge:     config.GetEnv("BRIDGE", "vmbr0"), // Restored!
		Workers:    config.GetEnvAsInt("WORKERS", 2),
	}

	isos, err := filepath.Glob(filepath.Join(cfg.TargetDir, "*.iso"))
	if err != nil || len(isos) == 0 {
		log.Fatalf("[ERROR] No ISOs found in directory: %s\n", cfg.TargetDir)
	}

	fmt.Println("===================================================================")
	fmt.Printf("🐣 PENGUINS INCUBATOR (Go Edition - Modular & Reporting)\n")
	fmt.Printf("Batch started: found %d ISOs in %s\n", len(isos), cfg.TargetDir)
	fmt.Printf("Configuration -> Firmware: %s | FsType: %s | Base VMID: %d\n", cfg.Firmware, cfg.FsType, cfg.BaseVMID)
	fmt.Println("===================================================================")

	numWorkers := cfg.Workers
	if numWorkers < 1 {
		numWorkers = 2
	}
	tasks := make(chan orchestrator.IsoTask, len(isos))

	// THE MISSING PIPE: Create a channel to collect the reports
	reportsChan := make(chan orchestrator.TestReport, len(isos))
	var wg sync.WaitGroup

	// Trigger the workers passing reportsChan as the fourth argument
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
	close(reportsChan) // Close the reports channel

	fmt.Println("\n>>> BATCH COMPLETED. Generating final report...")
	generateMarkdownSummary(reportsChan, cfg)
}

func generateMarkdownSummary(reports <-chan orchestrator.TestReport, cfg config.Config) {
	fileName := fmt.Sprintf("incubator-summary-%s-%s.md", cfg.Firmware, cfg.FsType)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Unable to create report file: %v", err)
		return
	}
	defer file.Close()

	header := fmt.Sprintf("## Incubator CI Results (%s / %s)\n\n", strings.ToUpper(cfg.Firmware), strings.ToUpper(cfg.FsType))
	header += "| Distro | ISO | Firmware | FileSystem | Status | Duration |\n"
	header += "|--------|-----|----------|------------|--------|----------|\n"

	file.WriteString(header)
	fmt.Print("\n" + header)

	for r := range reports {
		durata := r.Duration.Round(time.Second)
		// Convert the distro name to UPPERCASE for visual separation
		row := fmt.Sprintf("| **%s** | %s | %s | %s | %s | %s |\n",
			strings.ToUpper(r.Distro), r.IsoName, r.Firmware, r.FsType, r.Status, durata)

		file.WriteString(row)
		fmt.Print(row)
	}

	fmt.Printf("\nReport successfully saved in: %s\n", fileName)
}
