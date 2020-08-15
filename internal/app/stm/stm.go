package stm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/google/shlex"
)

type ManagerList []Manager

func (ml ManagerList) Names() []string {
	ret := []string{}
	for _, m := range ml {
		ret = append(ret, m.Name)
	}
	return ret
}

type ToolList []Tool

func (tl ToolList) FilterByManager(m string) ToolList {
	ret := ToolList{}
	for _, t := range tl {
		if m != t.Manager {
			continue
		}
		ret = append(ret, t)
	}
	return ret
}

type Config struct {
	Managers ManagerList `json:"managers"`
	Tools    ToolList    `json:"tools"`
}

func LoadConfig() (Config, error) {
	var config Config

	content, ferr := ioutil.ReadFile(path.Join(AppDir(), "config.json"))
	if ferr != nil {
		return config, ferr
	}

	if jerr := json.Unmarshal(content, &config); jerr != nil {
		return config, jerr
	}

	return config, nil
}

func (c Config) FindManager(m string) (Manager, error) {
	for _, cm := range c.Managers {
		if m == cm.Name {
			return cm, nil
		}
	}

	return Manager{}, errors.New(fmt.Sprintf("%s manager not found", m))
}

type Manager struct {
	Name       string `json:"name"`
	InstallCmd string `json:"install_command"`
	UpdateCmd  string `json:"update_command"`
}

func (m Manager) InstallPackages(pkgs []string) error {
	c, err := renderTemplate(m.Name, m.InstallCmd, pkgs)
	if err != nil {
		return err
	}

	cmd := prepareCommand(c)
	cmd.Run()
	return nil
}

func (m Manager) UpdatePackages(pkgs []string) error {
	c, err := renderTemplate(m.Name, m.UpdateCmd, pkgs)
	if err != nil {
		return err
	}

	cmd := prepareCommand(c)
	cmd.Run()
	return nil
}

type CmdParams struct {
	Manager  string
	Packages string
}

func renderTemplate(manager string, tmpl string, pkgs []string) (string, error) {
	t := template.Must(template.New(manager).Parse(tmpl))

	c := CmdParams{
		Manager:  manager,
		Packages: strings.Join(pkgs, " "),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, c); err != nil {
		return "", err
	}

	return os.ExpandEnv(buf.String()), nil
}

func prepareCommand(c string) *exec.Cmd {
	splCmd, err := shlex.Split(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(splCmd) < 1 {
		fmt.Println("install_command is required")
		os.Exit(1)
	}

	b := splCmd[0]
	args := []string{}
	if len(splCmd) > 1 {
		args = splCmd[1:]
	}

	cmd := exec.Command(b, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

type Tool struct {
	Package string `json:"package"`
	Binary  string `json:"binary,omitempty"`
	Path    string `json:"path,omitempty"`
	Manager string `json:"manager"`
}

func (t Tool) Installed() bool {
	b := os.ExpandEnv(t.Binary)
	_, err := exec.LookPath(b)
	if err == nil {
		return true
	}

	p := os.ExpandEnv(t.Path)
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}

	return false
}

func AppDir() string {
	appDir := os.Getenv("STM_CONFIG_PATH")
	if appDir != "" {
		return appDir
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(configHome, "stm")
}

func init() {
	appDir := AppDir()
	os.Setenv("STM_CONFIG_PATH", appDir)
	os.MkdirAll(appDir, 0755)
}
