package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	t.Run("get custom config path", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "TestGetConfigDir")
		if err != nil {
			t.Fatal(err)
		}

		os.Setenv("STM_CONFIG_PATH", path.Join(tmpDir, "stm"))

		want := path.Join(tmpDir, "stm")
		got := GetConfigDir()
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}

		os.Setenv("STM_CONFIG_PATH", "")
	})

	t.Run("get default config path with xdg", func(t *testing.T) {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			t.Skip("XDG_CONFIG_HOME not defined")
		}

		got := GetConfigDir()
		want := path.Join(configHome, "stm")
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("get default config path without xdg", func(t *testing.T) {
		homePath := path.Join(os.Getenv("HOME"), ".config")

		got := GetConfigDir()
		want := path.Join(homePath, "stm")
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestSetupConfigPath(t *testing.T) {
	tmpDir, terr := ioutil.TempDir("", "TestSetupConfigPath")
	if terr != nil {
		t.Fatal(terr)
	}

	os.Setenv("STM_CONFIG_PATH", path.Join(tmpDir, "stm"))

	configDir := GetConfigDir()
	SetupConfigPath()
	stat, err := os.Stat(configDir)
	if err != nil {
		t.Errorf("%s not exists", configDir)
	}

	gotMode := stat.Mode()
	wantMode := os.FileMode(0755) | os.ModeDir
	if gotMode != wantMode {
		t.Errorf("got mode %s want mode %s", gotMode, wantMode)
	}

	os.Setenv("STM_CONFIG_PATH", "")
}
