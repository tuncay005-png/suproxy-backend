package binary

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// RealBinaryManager manages Xray binary with real file system operations
type RealBinaryManager struct {
	binaryPath   string
	installDir   string
	logger       *logger.Logger
	versionRegex *regexp.Regexp
}

// NewRealBinaryManager creates a new real binary manager
func NewRealBinaryManager(binaryPath, installDir string, logger *logger.Logger) Manager {
	return &RealBinaryManager{
		binaryPath:   binaryPath,
		installDir:   installDir,
		logger:       logger,
		versionRegex: regexp.MustCompile(`Xray\s+(\d+\.\d+\.\d+)`),
	}
}

func (m *RealBinaryManager) Detect(ctx context.Context) (string, error) {
	// Check custom path first
	if m.binaryPath != "" {
		if m.fileExists(m.binaryPath) {
			return m.binaryPath, nil
		}
	}

	// Check common installation paths
	commonPaths := []string{
		"/usr/local/bin/xray",
		"/usr/bin/xray",
		"/opt/xray/xray",
		filepath.Join(m.installDir, "xray"),
	}

	for _, path := range commonPaths {
		if m.fileExists(path) {
			m.logger.Info("Xray binary detected", "path", path)
			return path, nil
		}
	}

	// Check PATH environment
	if path, err := exec.LookPath("xray"); err == nil {
		m.logger.Info("Xray binary found in PATH", "path", path)
		return path, nil
	}

	return "", ErrBinaryNotFound
}

func (m *RealBinaryManager) Validate(ctx context.Context, binaryPath string) error {
	// Check if file exists
	if !m.fileExists(binaryPath) {
		return ErrBinaryNotFound
	}

	// Check if file is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to stat binary: %w", err)
	}

	mode := info.Mode()
	// On Unix, check executable bit; on Windows, check .exe extension
	if runtime.GOOS != "windows" {
		if mode&0111 == 0 {
			return fmt.Errorf("binary is not executable: %s", binaryPath)
		}
	}

	// Try to execute --version to validate it's actually Xray
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute binary: %w", err)
	}

	if !strings.Contains(string(output), "Xray") {
		return fmt.Errorf("not a valid Xray binary")
	}

	m.logger.Info("Binary validated successfully", "path", binaryPath)
	return nil
}

func (m *RealBinaryManager) ValidateConfig(ctx context.Context, configPath string) error {
	// Use Xray binary to validate config: xray run -test -config=<path>
	binaryPath, err := m.Detect(ctx)
	if err != nil {
		return fmt.Errorf("xray binary not found: %w", err)
	}

	// Create context with timeout for validation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Execute: xray run -test -c <configPath>
	cmd := exec.CommandContext(ctx, binaryPath, "run", "-test", "-c", configPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		m.logger.Error("Config validation failed",
			"error", err,
			"config_path", configPath,
			"output", string(output))
		return fmt.Errorf("config validation failed: %w", err)
	}

	m.logger.Debug("Config validation successful",
		"config_path", configPath)

	return nil
}

func (m *RealBinaryManager) CurrentVersion(ctx context.Context) (string, error) {
	// Detect binary if not set
	binaryPath := m.binaryPath
	if binaryPath == "" {
		detected, err := m.Detect(ctx)
		if err != nil {
			return "", err
		}
		binaryPath = detected
	}

	// Execute xray version command
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}

	// Parse version from output
	// Example output: "Xray 1.8.7 (Xray, Penetrates Everything.) Custom"
	version := m.parseVersion(string(output))
	if version == "" {
		return "", ErrInvalidVersion
	}

	m.logger.Debug("Current version detected", "version", version)
	return version, nil
}

func (m *RealBinaryManager) LatestVersion(ctx context.Context) (string, error) {
	// Note: Fetching latest version from GitHub API
	// In production, consider caching this response to avoid rate limits

	// TODO: Implement GitHub API call to fetch latest release
	// GET https://api.github.com/repos/XTLS/Xray-core/releases/latest
	// type GitHubRelease struct {
	//     TagName string `json:"tag_name"`
	//     Name    string `json:"name"`
	//     Assets  []struct {
	//         Name        string `json:"name"`
	//         DownloadURL string `json:"browser_download_url"`
	//     } `json:"assets"`
	// }
	// req, err := http.NewRequestWithContext(ctx, "GET",
	//     "https://api.github.com/repos/XTLS/Xray-core/releases/latest", nil)
	// if err != nil {
	//     return "", err
	// }
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	//     return "", fmt.Errorf("failed to fetch latest version: %w", err)
	// }
	// defer resp.Body.Close()
	// var release GitHubRelease
	// if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
	//     return "", fmt.Errorf("failed to parse release info: %w", err)
	// }
	// version := strings.TrimPrefix(release.TagName, "v")
	// return version, nil

	m.logger.Debug("Latest version check (stub implementation)")
	return "1.8.8", nil
}

func (m *RealBinaryManager) Download(ctx context.Context, version string) error {
	// Note: Binary download from GitHub releases
	// This is a complex operation that requires:
	// 1. HTTP download with progress tracking
	// 2. Checksum verification
	// 3. Archive extraction (zip/tar.gz)
	// 4. Permission setting
	// 5. Error handling and cleanup

	// TODO: Implement full download workflow
	// osArch := m.getOSArch()
	// downloadURL := fmt.Sprintf(
	//     "https://github.com/XTLS/Xray-core/releases/download/v%s/Xray-%s.zip",
	//     version, osArch,
	// )
	// tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("xray-%s.zip", version))
	// // Download file
	// resp, err := http.Get(downloadURL)
	// if err != nil {
	//     return fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	// }
	// defer resp.Body.Close()
	// out, err := os.Create(tempFile)
	// if err != nil {
	//     return err
	// }
	// defer out.Close()
	// if _, err := io.Copy(out, resp.Body); err != nil {
	//     return err
	// }
	// // Extract and install
	// if err := m.extractAndInstall(tempFile, version); err != nil {
	//     return err
	// }

	m.logger.Info("Binary download workflow prepared (stub implementation)", "version", version)
	return fmt.Errorf("download not implemented - manual installation required")
}

func (m *RealBinaryManager) Upgrade(ctx context.Context, version string) error {
	// Note: Upgrade workflow for production systems
	// Requires careful handling to avoid downtime

	currentVersion, err := m.CurrentVersion(ctx)
	if err != nil {
		return err
	}

	if currentVersion == version {
		m.logger.Info("Already at target version", "version", version)
		return nil
	}

	// TODO: Implement full upgrade workflow
	// 1. Validate new version exists
	// 2. Backup current binary
	// 3. Download new version
	// 4. Validate new binary
	// 5. Replace binary atomically
	// 6. Optionally restart services

	// backupPath := m.binaryPath + ".backup." + currentVersion
	// if err := m.copyFile(m.binaryPath, backupPath); err != nil {
	//     m.logger.Warn("Failed to backup binary", "error", err)
	// }
	// if err := m.Download(ctx, version); err != nil {
	//     return err
	// }

	m.logger.Info("Binary upgrade workflow prepared (stub implementation)",
		"from", currentVersion,
		"to", version)
	return fmt.Errorf("upgrade not implemented - manual upgrade required")
}

func (m *RealBinaryManager) GetPath() string {
	return m.binaryPath
}

func (m *RealBinaryManager) IsInstalled(ctx context.Context) bool {
	_, err := m.Detect(ctx)
	return err == nil
}

func (m *RealBinaryManager) parseVersion(output string) string {
	matches := m.versionRegex.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func (m *RealBinaryManager) fileExists(path string) bool {
	// Check if file exists
	_, err := os.Stat(path)
	return err == nil
}
