package orchestrator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pieroproietti/incubator/pkg/config"
	"github.com/pieroproietti/incubator/pkg/proxmox"
)

// IsoTask represents the single ISO to be tested
type IsoTask struct {
	IsoName string
	IsoPath string
	VMID    int
}

// TestReport contains the complete telemetry of the test for the dashboard
type TestReport struct {
	IsoName  string
	Distro   string
	Firmware string
	FsType   string
	Status   string
	ErrorMsg string
	Fstab    string
	SysSpecs string
	Duration time.Duration
}

// Worker fetches ISOs from the channel, runs the orchestration, and returns the report
func Worker(workerVMID int, cfg config.Config, tasks <-chan IsoTask, reports chan<- TestReport, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		// Crucial: assign the worker's license plate to the task!
		task.VMID = workerVMID

		startTime := time.Now()

		report := TestReport{
			IsoName:  task.IsoName,
			Distro:   strings.Split(strings.TrimPrefix(task.IsoName, "egg-of-"), "-")[0],
			Firmware: cfg.Firmware,
			FsType:   cfg.FsType,
		}

		fmt.Printf("\n[Worker VMID:%d] >>> STARTING TEST: %s\n", task.VMID, task.IsoName)

		success := RunIncubatorTest(task, cfg)

		report.Duration = time.Since(startTime)

		if success {
			report.Status = "✅ OK"
			fmt.Printf("[Worker VMID:%d] %s: %s\n", task.VMID, report.Status, task.IsoName)
		} else {
			report.Status = "❌ FAILED"
			report.ErrorMsg = "Error during the installation cycle"
			fmt.Printf("[Worker VMID:%d] %s: %s\n", task.VMID, report.Status, task.IsoName)
		}

		reports <- report
		fmt.Println("-------------------------------------------------------------------")
	}
}

// RunIncubatorTest manages the complete lifecycle of the installation test
func RunIncubatorTest(task IsoTask, cfg config.Config) bool {
	vmidStr := strconv.Itoa(task.VMID)
	logPrefix := fmt.Sprintf("[VMID:%s]", vmidStr)

	// PHASE 1: Purge
	fmt.Printf("%s 1. Purging old VM and destroying disks...\n", logPrefix)
	proxmox.RunCommand("qm", "stop", vmidStr, "--timeout", "15")
	proxmox.RunCommand("qm", "destroy", vmidStr, "--purge", "1", "--destroy-unreferenced-disks", "1")

	// PHASE 2: Provisioning (with dynamic UEFI/BIOS support)
	fmt.Printf("%s 2. Configuring new VM (Firmware: %s)...\n", logPrefix, strings.ToUpper(cfg.Firmware))
	if cfg.Template != "" {
		out, err := proxmox.RunCommand("qm", "clone", cfg.Template, vmidStr, "--name", "testing-krill", "--full", "0")
		if err != nil {
			fmt.Printf("%s [ERROR] Failed template cloning: %v\nProxmox Details: %s\n", logPrefix, err, out)
			return false
		}
		proxmox.RunCommand("qm", "set", vmidStr, "--ide2", cfg.IsoStorage+":iso/"+task.IsoName+",media=cdrom", "--boot", "order=scsi0;ide2")
	} else {
		// Dynamic array for common arguments
		args := []string{
			"create", vmidStr,
			"--name", "testing-krill-" + cfg.Firmware,
			"--memory", "4096",
			"--cores", "2",
			"--scsihw", "virtio-scsi-single",
			"--scsi0", cfg.Storage + ":16",
			"--net0", "virtio,bridge=" + cfg.Bridge,
			"--serial0", "socket",
			"--vga", "qxl",
			"--agent", "1",
			"--ide2", cfg.IsoStorage + ":iso/" + task.IsoName + ",media=cdrom",
			"--boot", "order=scsi0;ide2",
		}

		// Append EFI parameters if required
		if strings.ToLower(cfg.Firmware) == "uefi" {
			args = append(args,
				"--bios", "ovmf",
				"--machine", "q35",
				"--efidisk0", cfg.Storage+":0,pre-enrolled-keys=0", // Disables Secure Boot
			)
		}

		out, err := proxmox.RunCommand("qm", args...)
		if err != nil {
			fmt.Printf("%s [ERROR] Failed VM creation: %v\nProxmox Details: %s\n", logPrefix, err, out)
			return false
		}
	}

	// PHASE 3: Live Boot
	fmt.Printf("%s 3. Booting Live system and waiting for QEMU Agent...\n", logPrefix)
	proxmox.RunCommand("qm", "start", vmidStr)

	if !proxmox.WaitForAgent(vmidStr, 300) {
		fmt.Printf("%s [FATAL ERROR] Agent did not respond after 5 minutes.\n", logPrefix)
		proxmox.RunCommand("qm", "stop", vmidStr)
		return false
	}
	fmt.Printf("%s Agent connected! Settle time of 10 seconds...\n", logPrefix)
	time.Sleep(10 * time.Second)

	// PHASE 4: Krill Installation
	fmt.Printf("%s 4. Starting Krill installation in the background...\n", logPrefix)
	installCmd := fmt.Sprintf("sudo eggs sysinstall krill --unattended --fstype=%s > /tmp/krill.log 2>&1", cfg.FsType)
	out, err := proxmox.RunCommand("qm", "guest", "exec", vmidStr, "--", "/bin/sh", "-c", installCmd)
	if err != nil {
		fmt.Printf("%s [FATAL ERROR] Krill startup failed: %v\n", logPrefix, err)
		return false
	}

	pid := proxmox.ExtractPid(out)
	if pid == "" {
		fmt.Printf("%s [FATAL ERROR] Unable to extract Krill PID.\n", logPrefix)
		return false
	}
	fmt.Printf("%s Internal Krill PID: %s. Monitoring in progress...\n", logPrefix, pid)

	// PHASE 5: Monitoring and Anti-Hang Countermeasures
	statusFails := 0
	waited := 0
	installTimeout := 3600 // 1 hour limit

	for {
		vmStatus, _ := proxmox.RunCommand("qm", "status", vmidStr)
		if strings.Contains(vmStatus, "stopped") {
			fmt.Printf("%s The VM shut down autonomously. Krill installation completed!\n", logPrefix)
			break
		}

		statusOut, err := proxmox.RunCommand("qm", "guest", "exec-status", vmidStr, pid)
		if err != nil {
			statusFails++
			// GHOST FINGER: simulates ENTER keypress
			proxmox.RunCommand("qm", "sendkey", vmidStr, "ret")

			if statusFails >= 4 {
				fmt.Printf("%s [GUILLOTINE] Agent offline, VM stuck at reboot. Forced shutdown!\n", logPrefix)
				proxmox.RunCommand("qm", "stop", vmidStr, "--timeout", "15")
				break
			}
			time.Sleep(5 * time.Second)
			waited += 5
			continue
		}
		statusFails = 0

		var status proxmox.AgentStatus
		if err := json.Unmarshal([]byte(statusOut), &status); err == nil {
			isExited := false
			switch v := status.Exited.(type) {
			case bool:
				isExited = v
			case float64:
				isExited = (v == 1)
			}

			if isExited {
				if status.ExitCode != 0 {
					fmt.Printf("%s [ERROR] Krill returned ExitCode: %d\n", logPrefix, status.ExitCode)
					return false
				}
				fmt.Printf("%s Krill script finished successfully. Waiting for shutdown...\n", logPrefix)
			}
		}

		time.Sleep(5 * time.Second)
		waited += 5
		if waited >= installTimeout {
			fmt.Printf("%s [FATAL ERROR] Timeout reached (1 hour). Installation aborted.\n", logPrefix)
			proxmox.RunCommand("qm", "stop", vmidStr)
			return false
		}
	}

	// PHASE 6: Site Closure
	proxmox.RunCommand("qm", "stop", vmidStr)
	fmt.Printf("%s Test concluded successfully.\n", logPrefix)
	return true
}
