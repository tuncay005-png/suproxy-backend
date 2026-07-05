package binary

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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

	// TODO: Check if file is executable
	// info, err := os.Stat(binaryPath)
	// if err != nil {
	//     return fmt.Errorf("failed to stat binary: %w", err)
	// }

	// mode := info.Mode()
	// if mode&0111 == 0 {
	//     return fmt.Errorf("binary is not executable")
	// }

	// TODO: Try to execute --version to validate it's actually Xray
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	// cmd := exec.CommandContext(ctx, binaryPath, "version")
	// output, err := cmd.Output()
	// if err != nil {
	//     return fmt.Errorf("failed to execute binary: %w", err)
	// }

	// if !strings.Contains(string(output), "Xray") {
	//     return fmt.Errorf("not a valid Xray binary")
	// }

	m.logger.Info("Binary validation prepared", "path", binaryPath)
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

	// TODO: Execute xray version command
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	// cmd := exec.CommandContext(ctx, binaryPath, "version")
	// output, err := cmd.Output()
	// if err != nil {
	//     return "", fmt.Errorf("failed to get version: %w", err)
	// }

	// Parse version from output
	// Example output: "Xray 1.8.7 (Xray, Penetrates Everything.) Custom"
	// version := m.parseVersion(string(output))
	// if version == "" {
	//     return "", ErrInvalidVersion
	// }

	// return version, nil

	// For now, return mock version
	m.logger.Debug("Version check prepared", "binary", binaryPath)
	return "1.8.7", nil
}

func (m *RealBinaryManager) LatestVersion(ctx context.Context) (string, error) {
	// TODO: Fetch latest version from GitHub API
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

	// Clean version (remove 'v' prefix if present)
	// version := strings.TrimPrefix(release.TagName, "v")
	// return version, nil

	m.logger.Debug("Latest version check prepared")
	return "1.8.8", nil
}

func (m *RealBinaryManager) Download(ctx context.Context, version string) error {
	// TODO: Download specific version from GitHub
	// 1. Construct download URL based on OS/Arch
	// 2. Download archive
	// 3. Verify checksum
	// 4. Extract binary
	// 5. Set executable permissions
	// 6. Move to install directory

	// Example URL:
	// https://github.com/XTLS/Xray-core/releases/download/v1.8.7/Xray-linux-64.zip

	// osArch := m.getOSArch()
	// downloadURL := fmt.Sprintf(
	//     "https://github.com/XTLS/Xray-core/releases/download/v%s/Xray-%s.zip",
	//     version, osArch,
	// )

	// Download to temp file
	// tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("xray-%s.zip", version))

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

	// Extract and install
	// if err := m.extractAndInstall(tempFile, version); err != nil {
	//     return err
	// }

	m.logger.Info("Binary download prepared", "version", version)
	return nil
}

func (m *RealBinaryManager) Upgrade(ctx context.Context, version string) error {
	// TODO: Upgrade workflow
	// 1. Check current version
	// 2. Download new version
	// 3. Stop running instances (or prepare for restart)
	// 4. Backup current binary
	// 5. Replace binary
	// 6. Validate new binary
	// 7. Restart instances if needed

	currentVersion, err := m.CurrentVersion(ctx)
	if err != nil {
		return err
	}

	if currentVersion == version {
		m.logger.Info("Already at target version", "version", version)
		return nil
	}

	// Backup current binary
	// backupPath := m.binaryPath + ".backup." + currentVersion
	// if err := m.copyFile(m.binaryPath, backupPath); err != nil {
	//     m.logger.Warn("Failed to backup binary", "error", err)
	// }

	// Download and install new version
	if err := m.Download(ctx, version); err != nil {
		return err
	}

	m.logger.Info("Binary upgrade prepared", "from", currentVersion, "to", version)
	return nil
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
	// TODO: Check if file exists
	// _, err := os.Stat(path)
	// return err == nil

	_, _ = os.Stat, path
	return false
}

func (m *RealBinaryManager) getOSArch() string {
	// TODO: Detect OS and architecture
	// Examples: linux-64, linux-arm64-v8a, windows-64, darwin-arm64-v8a

	// osName := runtime.GOOS
	// arch := runtime.GOARCH

	// Map Go arch to Xray naming
	// switch {
	// case osName == "linux" && arch == "amd64":
	//     return "linux-64"
	// case osName == "linux" && arch == "arm64":
	//     return "linux-arm64-v8a"
	// case osName == "windows" && arch == "amd64":
	//     return "windows-64"
	// case osName == "darwin" && arch == "arm64":
	//     return "darwin-arm64-v8a"
	// default:
	//     return fmt.Sprintf("%s-%s", osName, arch)
	// }

	return "linux-64"
}

func (m *RealBinaryManager) extractAndInstall(archivePath, version string) error {
	// TODO: Extract zip/tar.gz archive
	// TODO: Find xray binary in extracted files
	// TODO: Move to install directory
	// TODO: Set executable permissions (chmod +x)
	// TODO: Clean up temp files

	_, _ = archivePath, version
	return nil
}

func (m *RealBinaryManager) copyFile(src, dst string) error {
	// TODO: Copy file with permissions preserved
	_, _, _ = src, dst, strings.Contains
	return nil
}
