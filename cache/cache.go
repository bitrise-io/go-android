package cache

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/hashicorp/go-version"
)

// Level defines the extent to which caching should be used.
// - LevelNone: no caching
// - LevelDeps: only dependencies will be cached
// - LevelAll: caching will include gradle and android build cache
type Level string

// Cache level
const (
	LevelNone = Level("none")
	LevelDeps = Level("only_deps")
	LevelAll  = Level("all")
)

func parseGradlewVersion(out string) (string, error) {
	/*
	   ------------------------------------------------------------
	   Gradle 6.1.1
	   ------------------------------------------------------------

	   Build time:   2020-01-24 22:30:24 UTC
	   Revision:     a8c3750babb99d1894378073499d6716a1a1fa5d

	   Kotlin:       1.3.61
	   Groovy:       2.5.8
	   Ant:          Apache Ant(TM) version 1.10.7 compiled on September 1 2019
	   JVM:          1.8.0_241 (Oracle Corporation 25.241-b07)
	   OS:           Mac OS X 10.15.5 x86_64
	*/

	pattern := `Gradle (.*)`
	exp := regexp.MustCompile(pattern)
	matches := exp.FindStringSubmatch(out)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find Gradle version in output:\n%s\nusing a pattern:%s", out, pattern)
	}
	return matches[1], nil
}

func projectGradleVersion(projectPth string) (string, error) {
	gradlewPth := filepath.Join(projectPth, "gradlew")
	exist, err := pathutil.IsPathExists(gradlewPth)
	if err != nil {
		return "", fmt.Errorf("failed to check if %s exists: %s", gradlewPth, err)
	}
	if !exist {
		return "", fmt.Errorf("no gradlew found at: %s", gradlewPth)
	}

	versionCmd := command.New("./gradlew", "-version")
	versionCmd.SetDir(filepath.Dir(gradlewPth))
	out, err := versionCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return "", fmt.Errorf("%s failed: %s", versionCmd.PrintableCommandArgs(), out)
		}
		return "", fmt.Errorf("%s failed: %s", versionCmd.PrintableCommandArgs(), err)
	}

	return parseGradlewVersion(out)
}

func oldGradleExcludePaths(homeDir, currentGradleVersion string) ([]string, error) {
	var excludes []string

	gradlewDir := filepath.Join(homeDir, ".gradle")
	exist, err := pathutil.IsPathExists(gradlewDir)
	if err != nil {
		return nil, fmt.Errorf("failed to check if %s exist: %s", gradlewDir, err)
	}
	if !exist {
		return nil, nil
	}

	{
		// exclude old wrappers, like ~/.gradle/wrapper/dists/gradle-5.1.1-all
		wrapperDistrDir := filepath.Join(gradlewDir, "wrapper", "dists")
		entries, err := ioutil.ReadDir(wrapperDistrDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read entries of %s: %s", wrapperDistrDir, err)
		}
		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "gradle-"+currentGradleVersion) {
				excludes = append(excludes, "!"+filepath.Join(wrapperDistrDir, e.Name()))
			}
		}
	}

	{
		// exclude old caches, like ~/.gradle/caches/5.1.1
		cachesDir := filepath.Join(gradlewDir, "caches")
		entries, err := ioutil.ReadDir(cachesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read entries of %s: %s", cachesDir, err)
		}
		for _, e := range entries {
			v, err := version.NewVersion(e.Name())
			if err != nil || v == nil {
				continue
			}

			if e.Name() != currentGradleVersion {
				excludes = append(excludes, "!"+filepath.Join(cachesDir, e.Name()))
			}
		}
	}

	{
		// exclude old daemon, like ~/.gradle/daemon/5.1.1
		daemonDir := filepath.Join(gradlewDir, "daemon")
		entries, err := ioutil.ReadDir(daemonDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read entries of %s: %s", daemonDir, err)
		}
		for _, e := range entries {
			v, err := version.NewVersion(e.Name())
			if err != nil || v == nil {
				continue
			}

			if e.Name() != currentGradleVersion {
				excludes = append(excludes, "!"+filepath.Join(daemonDir, e.Name()))
			}
		}
	}

	return excludes, nil
}

// Collect walks the directory tree underneath projectRoot and registers matching
// paths for caching based on the value of cacheLevel. Returns an error if there
// was an underlying error that would lead to a corrupted cache file, otherwise
// the given path is skipped.
func Collect(projectRoot string, cacheLevel Level) error {
	if cacheLevel == LevelNone {
		return nil
	}

	gradleCache := cache.New()

	var includePths []string
	excludePths := []string{
		"~/.gradle/**",
		"!~/.gradle/daemon/*/daemon-*.out.log", // excludes Gradle daemon logs, like: ~/.gradle/daemon/6.1.1/daemon-3122.out.log
		"~/.android/build-cache/**",
		"*.lock",
		"*.bin",
		"*/build/*.json",
		"*/build/*.html",
		"*/build/*.xml",
		"*/build/*.properties",
		"*/build/*/zip-cache/*",
		"*.log",
		"*.txt",
		"*.rawproto",
		"!*.ap_",
		"!*.apk",
	}

	homeDir := pathutil.UserHomeDir()

	projectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return fmt.Errorf("cache collection skipped: failed to determine project root path")
	}

	ver, err := projectGradleVersion(projectRoot)
	if err != nil {
		log.Warnf("failed to get project gradle version: %s", err)
	} else {
		excludes, err := oldGradleExcludePaths(homeDir, ver)
		if err != nil {
			log.Warnf("failed to collect old gradle exclude paths: %s", err)
		} else {
			fmt.Printf("old gradle exclude path:\n%s\n", strings.Join(excludes, "\n"))
			excludePths = append(excludePths, excludes...)
		}
	}

	lockFilePath := filepath.Join(projectRoot, "gradle.deps")

	lockfileContent := ""
	if err := filepath.Walk(projectRoot, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".gradle") && !strings.Contains(path, "node_modules") {
			if md5Hash, err := computeMD5String(path); err != nil {
				log.Warnf("Failed to compute MD5 hash of file(%s), error: %s", path, err)
			} else {
				lockfileContent += md5Hash
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("dependency map generation skipped: failed to collect dependencies")
	}
	if err := fileutil.WriteStringToFile(lockFilePath, lockfileContent); err != nil {
		return fmt.Errorf("dependency map generation skipped: failed to write lockfile, error: %s", err)
	}

	includePths = append(includePths, fmt.Sprintf("%s -> %s", filepath.Join(homeDir, ".gradle"), lockFilePath))
	includePths = append(includePths, fmt.Sprintf("%s -> %s", filepath.Join(homeDir, ".kotlin"), lockFilePath))
	includePths = append(includePths, fmt.Sprintf("%s -> %s", filepath.Join(homeDir, ".m2"), lockFilePath))

	if cacheLevel == LevelAll {
		includePths = append(includePths, fmt.Sprintf("%s -> %s", filepath.Join(homeDir, ".android", "build-cache"), lockFilePath))

		if err := filepath.Walk(projectRoot, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				if f.Name() == "build" {
					includePths = append(includePths, path)
				}
				if f.Name() == ".gradle" {
					includePths = append(includePths, path)
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("cache collection skipped: failed to determine cache paths")
		}
	}

	gradleCache.IncludePath(strings.Join(includePths, "\n"))
	gradleCache.ExcludePath(strings.Join(excludePths, "\n"))

	if err := gradleCache.Commit(); err != nil {
		return fmt.Errorf("cache collection skipped: failed to commit cache paths")
	}

	return nil
}

func computeMD5String(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Failed to close file(%s), error: %s", filePath, err)
		}
	}()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
