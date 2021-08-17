package config

import (
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type singleton struct {
	Viper viper.Viper
}

type ZarfFile struct {
	Url        string
	Shasum     string
	Target     string
	Executable bool
}

type ZarfChart struct {
	Name    string
	Url     string
	Version string
}

type ZarfFeature struct {
	Name        string
	Description string
	Default     bool
	Manifests   string
	Images      []string
	Files       []ZarfFile
	Charts      []ZarfChart
}

type ZarfMetatdata struct {
	Name         string
	Description  string
	Version      string
	Uncompressed bool
}

type ZarfData struct {
	Source string
	Target struct {
		Namespace string
		Selector  string
		Path      string
	}
}

const K3sBinary = "/usr/local/bin/k3s"
const K3sChartPath = "/var/lib/rancher/k3s/server/static/charts"
const K3sManifestPath = "/var/lib/rancher/k3s/server/manifests"
const K3sImagePath = "/var/lib/rancher/k3s/agent/images"
const PackageInitName = "zarf-init.tar.zst"
const PackagePrefix = "zarf-package-"
const ZarfLocal = "zarf.localhost"

var instance *singleton
var once sync.Once

func getInstance() *singleton {
	once.Do(func() {
		instance = &singleton{Viper: *viper.New()}
		setupViper()
	})
	return instance
}

func IsZarfInitConfig() bool {
	var kind string
	getInstance().Viper.UnmarshalKey("kind", &kind)
	return strings.ToLower(kind) == "zarfinitconfig"
}

func GetPackageName() string {
	metadata := GetMetaData()
	if metadata.Uncompressed {
		return PackagePrefix + metadata.Name + ".tar"
	} else {
		return PackagePrefix + metadata.Name + ".tar.zst"
	}
}

func GetDataInjections() []ZarfData {
	var data []ZarfData
	getInstance().Viper.UnmarshalKey("data", &data)
	return data
}

func GetMetaData() ZarfMetatdata {
	var metatdata ZarfMetatdata
	getInstance().Viper.UnmarshalKey("metadata.name", &metatdata.Name)
	getInstance().Viper.UnmarshalKey("metatdata.description", &metatdata.Description)
	getInstance().Viper.UnmarshalKey("metatdata.version", &metatdata.Version)
	getInstance().Viper.UnmarshalKey("metadata.uncompressed", &metatdata.Uncompressed)
	return metatdata
}

func GetLocalCharts() []ZarfChart {
	var charts []ZarfChart
	getInstance().Viper.UnmarshalKey("local.charts", &charts)
	return charts
}

func GetLocalFiles() []ZarfFile {
	var files []ZarfFile
	getInstance().Viper.UnmarshalKey("local.files", &files)
	return files
}

func GetLocalImages() []string {
	var images []string
	getInstance().Viper.UnmarshalKey("local.images", &images)
	return images
}

func GetLocalManifests() string {
	var manifests string
	getInstance().Viper.UnmarshalKey("local.manifests", &manifests)
	return manifests
}

func GetInitFeatures() []ZarfFeature {
	var features []ZarfFeature
	getInstance().Viper.UnmarshalKey("features", &features)
	return features
}

func GetRemoteImages() []string {
	var images []string
	getInstance().Viper.UnmarshalKey("remote.images", &images)
	return images
}

func GetRemoteRepos() []string {
	var repos []string
	getInstance().Viper.UnmarshalKey("remote.repos", &repos)
	return repos
}

func DynamicConfigLoad(path string) {
	logContext := logrus.WithField("path", path)
	logContext.Info("Loading dynamic config")
	getInstance().Viper.SetConfigFile(path)
	if err := getInstance().Viper.MergeInConfig(); err != nil {
		logContext.Warn("Unable to load the config file")
	}
}

func WriteConfig(path string) {
	now := time.Now()
	currentUser, userErr := user.Current()
	hostname, hostErr := os.Hostname()

	// Record the time of package creation
	getInstance().Viper.Set("package.timestamp", now.Format(time.RFC1123Z))
	if hostErr == nil {
		// Record the hostname of the package creation terminal
		getInstance().Viper.Set("package.terminal", hostname)
	}
	if userErr == nil {
		// Record the name of the user creating the package
		getInstance().Viper.Set("package.user", currentUser.Name)
	}
	// Save the parsed output to the config path given
	if err := getInstance().Viper.WriteConfigAs(path); err != nil {
		logrus.WithField("path", path).Fatal("Unable to write the config file")
	}
}

func setupViper() {
	instance.Viper.AddConfigPath(".")
	instance.Viper.SetConfigName("zarf-config")

	// If a config file is found, read it in.
	if err := instance.Viper.ReadInConfig(); err == nil {
		logrus.WithField("path", instance.Viper.ConfigFileUsed()).Info("Config file loaded")
	}
}
