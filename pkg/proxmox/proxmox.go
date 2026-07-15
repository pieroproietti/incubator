package proxmox

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

// AgentStatus decodifica il JSON di risposta di Proxmox
type AgentStatus struct {
	Exited   interface{} `json:"exited"` // Può essere bool o float64 dipendendo da Proxmox
	ExitCode int         `json:"exitcode"`
	OutData  string      `json:"out-data"`
	ErrData  string      `json:"err-data"`
}

// RunCommand esegue un comando di sistema e cattura Stdout e Stderr uniti
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// WaitForAgent pinga il guest ininterrottamente fino al successo o al timeout
func WaitForAgent(vmid string, timeoutSeconds int) bool {
	waited := 0
	okCount := 0
	for waited < timeoutSeconds {
		_, err := RunCommand("qm", "agent", vmid, "ping")
		if err == nil {
			okCount++
			if okCount >= 3 {
				return true
			}
			time.Sleep(2 * time.Second)
			waited += 2
		} else {
			okCount = 0
			time.Sleep(3 * time.Second)
			waited += 3
		}
	}
	return false
}

// ExtractPid fa il parsing brutale dell'output iniziale di qm guest exec
func ExtractPid(out string) string {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"pid\"") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				pidRaw := strings.TrimSpace(parts[1])
				pidRaw = strings.Trim(pidRaw, " },")
				return pidRaw
			}
		}
	}
	return ""
}
