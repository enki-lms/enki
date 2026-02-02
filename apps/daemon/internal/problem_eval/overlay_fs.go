//go:build linux
// +build linux

package problem_eval

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// OverlayMount holds the paths for our overlay setup.
type OverlayMount struct {
	InstancePath string // The unique parent dir for this mount
	LowerDir     string
	UpperDir     string
	WorkDir      string
	MergedDir    string
}

// NewOverlayMount prepares the directory structure for an overlay mount.
func NewOverlayMount(baseDir, lowerDir string) (*OverlayMount, error) {
	instancePath, err := os.MkdirTemp(baseDir, "jail-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp instance dir: %w", err)
	}

	m := &OverlayMount{
		InstancePath: instancePath,
		LowerDir:     lowerDir,
		UpperDir:     filepath.Join(instancePath, "upper"),
		WorkDir:      filepath.Join(instancePath, "work"),
		MergedDir:    filepath.Join(instancePath, "merged"),
	}

	// Create the necessary directories
	for _, dir := range []string{m.UpperDir, m.WorkDir, m.MergedDir} {
		if err := os.Mkdir(dir, 0755); err != nil {
			// Cleanup partially created dirs on failure
			os.RemoveAll(instancePath)
			return nil, fmt.Errorf("failed to create overlay dir %s: %w", dir, err)
		}
	}

	return m, nil
}

// Mount performs the overlay mount operation.
func (m *OverlayMount) Mount() error {
	opts := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s",
		m.LowerDir, m.UpperDir, m.WorkDir,
	)
	return syscall.Mount("overlay", m.MergedDir, "overlay", 0, opts)
}

// Unmount unmounts and cleans up the directories.
func (m *OverlayMount) Unmount() error {
	if err := syscall.Unmount(m.MergedDir, 0); err != nil {
		return fmt.Errorf("failed to unmount %s: %w", m.MergedDir, err)
	}
	return os.RemoveAll(m.InstancePath)
}
