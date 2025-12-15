package inspection

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
)

type V2VSession struct {
	NBDURL   string
	cmd      *exec.Cmd
	authFile string // Path to libvirt auth file (for cleanup)
}

func OpenWithVirtV2V(
	ctx context.Context,
	vmName string,
	datacenter string,
	snapshotName string,
	vcenterURL string,
	username string,
	password string,
) (*V2VSession, error) {

	parsedURL, err := url.Parse(vcenterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vCenter URL: %w", err)
	}

	vcenterHost := parsedURL.Hostname()

	if datacenter == "" {
		return nil, fmt.Errorf("datacenter cannot be empty")
	}

	// Create libvirt auth file with credentials for security
	// This avoids embedding password in the URL
	authFile, err := createLibvirtAuthFile(vcenterHost, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create libvirt auth file: %w", err)
	}

	// Build vpx source URL WITHOUT credentials (credentials via auth file)
	vpxURL := fmt.Sprintf(
		"vpx://%s/%s/%s?snapshot=%s&no_verify=1",
		vcenterHost,
		datacenter,
		vmName,
		snapshotName,
	)

	args := []string{
		"-it", "vddk",
		vpxURL,
		"-o", "nbd",
	}

	cmd := exec.CommandContext(ctx, "virt-v2v-open", args...)

	// Set LIBVIRT_AUTH_FILE environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("LIBVIRT_AUTH_FILE=%s", authFile))

	// Discard stdout/stderr to avoid credential leaks in logs
	// virt-v2v-open doesn't produce useful output that needs logging
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		os.Remove(authFile) // Clean up auth file on error
		return nil, fmt.Errorf("failed to start virt-v2v-open: %w", err)
	}

	// Default port used by virt-v2v-open
	nbdURL := "nbd://localhost:10809"

	return &V2VSession{
		NBDURL:   nbdURL,
		cmd:      cmd,
		authFile: authFile,
	}, nil
}

func (s *V2VSession) Close() {
	if s == nil {
		return
	}

	// Kill the process
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_, _ = s.cmd.Process.Wait()
	}

	// Clean up the auth file
	if s.authFile != "" {
		_ = os.Remove(s.authFile)
	}
}

// createLibvirtAuthFile creates a libvirt auth.conf file with username and password
// This allows us to avoid embedding credentials in the vpx:// URL
// Format: https://libvirt.org/auth.html#client-configuration
func createLibvirtAuthFile(vcenterHost, username, password string) (string, error) {
	tmpFile, err := os.CreateTemp("", "libvirt-auth-*.conf")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary auth file: %w", err)
	}

	// Create libvirt auth.conf content
	authContent := fmt.Sprintf(`[credentials-vcenter]
username=%s
password=%s

[auth-esx-%s]
credentials=vcenter
`, username, password, vcenterHost)

	// Write auth content to file
	if _, err := tmpFile.WriteString(authContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write auth file: %w", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to close auth file: %w", err)
	}

	// Set restrictive permissions (read only by owner)
	if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to set auth file permissions: %w", err)
	}

	return tmpFile.Name(), nil
}
