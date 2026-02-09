package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	owner         = "ErickMaria"
	repo          = "envcontainer"
	githubAPIURL  = "https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest"
	cacheFileName = "envcontainer_version_cache"
	cacheDuration = 24 * time.Hour // Cache valid for 24 hours
)

type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

type VersionInfo struct {
	LatestVersion  string             `json:"latest_version"`
	CurrentVersion string             `json:"current_version"`
	HasUpdate      bool               `json:"has_update"`
	ReleaseURL     string             `json:"release_url"`
	Changelogs     []ReleaseChangelog `json:"changelogs"`
	CachedAt       int64              `json:"cached_at"`
}

type ReleaseChangelog struct {
	Version string
	TagName string
	Name    string
	Body    string
}

// GetCacheDir returns the cache directory
func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".envcontainer", "cache")
	return cacheDir, nil
}

// CheckForUpdate checks if a new version is available
func CheckForUpdate(currentVersion string) (*VersionInfo, error) {
	return CheckForUpdateForce(currentVersion, false)
}

// CheckForUpdateForce checks if a new version is available, with option to force bypass cache
func CheckForUpdateForce(currentVersion string, forceCheck bool) (*VersionInfo, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	// Try to load from cache first (unless force check is enabled)
	if !forceCheck {
		versionInfo, err := loadFromCache(cacheDir, currentVersion)
		if err == nil && isValidCache(versionInfo) {
			return versionInfo, nil
		}
	}

	// Fetch latest version from GitHub
	latestVersion, releaseURL, err := fetchLatestVersion()
	if err != nil {
		// Return cached data if fetch fails (fallback)
		versionInfo, cacheErr := loadFromCache(cacheDir, currentVersion)
		if cacheErr == nil {
			return versionInfo, nil
		}
		return nil, err
	}

	// Compare versions
	hasUpdate := compareVersions(currentVersion, latestVersion)

	// Fetch changelogs if update is available
	var changelogs []ReleaseChangelog
	if hasUpdate {
		changelogs, _ = GetChangelogsBetweenVersions(currentVersion, latestVersion)
	}

	versionInfo := &VersionInfo{
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
		HasUpdate:      hasUpdate,
		ReleaseURL:     releaseURL,
		Changelogs:     changelogs,
		CachedAt:       time.Now().Unix(),
	}

	// Save to cache
	_ = saveToCache(cacheDir, versionInfo)

	return versionInfo, nil
}

func fetchLatestVersion() (string, string, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch latest version: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse release data: %w", err)
	}

	version := strings.TrimPrefix(release.TagName, "v")
	releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", owner, repo, release.TagName)

	return version, releaseURL, nil
}

func compareVersions(current, latest string) bool {
	// Simple version comparison: split by dots and compare numerically
	currParts := strings.Split(strings.TrimPrefix(current, "v"), ".")
	latestParts := strings.Split(strings.TrimPrefix(latest, "v"), ".")

	for i := 0; i < len(currParts) && i < len(latestParts); i++ {
		var currNum, latestNum int
		fmt.Sscanf(currParts[i], "%d", &currNum)
		fmt.Sscanf(latestParts[i], "%d", &latestNum)

		if latestNum > currNum {
			return true
		}
		if latestNum < currNum {
			return false
		}
	}

	return len(latestParts) > len(currParts)
}

func loadFromCache(cacheDir, currentVersion string) (*VersionInfo, error) {
	cacheFile := filepath.Join(cacheDir, cacheFileName)

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var versionInfo VersionInfo
	err = json.Unmarshal(data, &versionInfo)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

func saveToCache(cacheDir string, versionInfo *VersionInfo) error {
	// Create cache directory if it doesn't exist
	err := os.MkdirAll(cacheDir, 0o755)
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, cacheFileName)

	data, err := json.Marshal(versionInfo)
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFile, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func isValidCache(versionInfo *VersionInfo) bool {
	if versionInfo == nil {
		return false
	}

	cacheAge := time.Since(time.Unix(versionInfo.CachedAt, 0))
	return cacheAge < cacheDuration
}

// ClearCache removes the version cache file
func ClearCache() error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, cacheFileName)
	return os.Remove(cacheFile)
}

// GetDownloadURL returns the download URL for the binary for current OS
func GetDownloadURL(release *Release) string {
	// Determine OS and arch
	osName := strings.ToLower(os.Getenv("GOOS"))
	if osName == "" {
		osName = getSystemOS()
	}

	for _, asset := range release.Assets {
		if strings.Contains(strings.ToLower(asset.Name), osName) {
			return asset.DownloadURL
		}
	}

	return ""
}

func getSystemOS() string {
	// This is a fallback; normally GOOS env var or runtime.GOOS should be used
	output, err := os.ReadFile("/etc/os-release")
	if err == nil {
		content := string(output)
		if strings.Contains(content, "ubuntu") || strings.Contains(content, "debian") {
			return "linux"
		}
	}
	return "linux" // default
}

// GetChangelogsBetweenVersions fetches all releases between current and latest version
func GetChangelogsBetweenVersions(currentVersion, latestVersion string) ([]ReleaseChangelog, error) {
	// Check cache first
	cacheDir, _ := GetCacheDir()
	if cacheDir != "" {
		cached, _ := loadFromCache(cacheDir, currentVersion)
		if cached != nil && isValidCache(cached) && cached.LatestVersion == latestVersion && len(cached.Changelogs) > 0 {
			return cached.Changelogs, nil
		}
	}

	var changelogs []ReleaseChangelog

	// Fetch all releases from GitHub
	allReleases, err := fetchAllReleases()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}

	// Filter releases between current and latest version
	current := normalizeVersion(currentVersion)
	latest := normalizeVersion(latestVersion)

	for _, release := range allReleases {
		releaseVer := normalizeVersion(strings.TrimPrefix(release.TagName, "v"))

		// Include releases that are newer than current and not newer than latest
		if isVersionGreater(releaseVer, current) && !isVersionGreater(releaseVer, latest) {
			changelog := ReleaseChangelog{
				Version: strings.TrimPrefix(release.TagName, "v"),
				TagName: release.TagName,
				Name:    release.Name,
				Body:    release.Body,
			}
			changelogs = append(changelogs, changelog)
		}
	}

	return changelogs, nil
}

func fetchAllReleases() ([]Release, error) {
	var allReleases []Release
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?page=%d&per_page=30", owner, repo, page)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			break
		}

		var releases []Release
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		err = json.Unmarshal(body, &releases)
		if err != nil {
			return nil, err
		}

		if len(releases) == 0 {
			break
		}

		allReleases = append(allReleases, releases...)
		page++
	}

	return allReleases, nil
}

func normalizeVersion(version string) [3]int {
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	var normalized [3]int

	for i := 0; i < 3 && i < len(parts); i++ {
		var num int
		fmt.Sscanf(parts[i], "%d", &num)
		normalized[i] = num
	}

	return normalized
}

func isVersionGreater(version, compare [3]int) bool {
	for i := 0; i < 3; i++ {
		if version[i] > compare[i] {
			return true
		}
		if version[i] < compare[i] {
			return false
		}
	}
	return false
}
