package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Workspace struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
}

type Config struct {
	Current    string      `json:"current"`
	Workspaces []Workspace `json:"workspaces"`
	current    *Workspace
}

var ConfigFilename string

func config() (cfg Config) {
	hd, _ := os.UserHomeDir()
	ConfigFilename = hd + "/projects.json"

	b, err := ioutil.ReadFile(ConfigFilename)
	if os.IsNotExist(err) {
		cfg.update()
		return
	}

	if err = json.Unmarshal(b, &cfg); err != nil {
		cfg.update()
		return
	}

	cfg.current = cfg.getWorkspace(cfg.Current)

	return
}

func (c *Config) getWorkspace(name string) (workspace *Workspace) {
	for _, w := range c.Workspaces {
		if w.Name == name {
			return &w
		}
	}

	return
}

func (c *Config) changeCurrentWorkspace(name string) bool {
	found := c.getWorkspace(name) != nil

	if found {
		c.Current = name
		c.update()
	}

	return found
}

func (c *Config) addWorkspace(workspace Workspace) error {
	if c.getWorkspace(workspace.Name) != nil {
		return fmt.Errorf(`Workspace "%s" already exists`, workspace.Name)
	}

	// _, err := os.ReadDir(workspace.Directory)
	// if os.IsNotExist(err) {
	// 	return fmt.Errorf(`Workspace does not exist`)
	// }
	c.Workspaces = append(c.Workspaces, workspace)
	if len(c.Workspaces) == 1 {
		c.Current = workspace.Name
	}
	c.update()

	return nil
}

func (c *Config) removeWorkspace(workspace Workspace) bool {
	var found bool
	newWorkspaces := make([]Workspace, 0)

	for _, w := range c.Workspaces {
		if w.Name != workspace.Name {
			newWorkspaces = append(newWorkspaces, w)
		} else {
			found = true
		}
	}

	c.Workspaces = newWorkspaces

	c.update()
	return found
}

func (c *Config) deleteProject(project string) bool {
	return os.RemoveAll(c.current.Directory+"/"+project) == nil
}

func (c *Config) update() {
	b, _ := json.Marshal(c)
	ioutil.WriteFile(ConfigFilename, b, 0644)
}
