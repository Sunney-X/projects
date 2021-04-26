package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func Run() error {
	cfg := config()

	app := &cli.App{
		Name:  "projects",
		Usage: "A project manager app",
		UsageText: `projects <PROJECT_NAME>   Start a project instantly
   projects command [arguments...]`,
		Action: func(c *cli.Context) (err error) {
			pn := c.Args().First()
			if pn == "" {
				return cli.ShowAppHelp(c)
			}

			files, _ := os.ReadDir(cfg.current.Directory)
			for _, f := range files {
				if strings.EqualFold(pn, f.Name()) {
					exec.Command("code", cfg.current.Directory+"/"+f.Name()).Run()
					return
				}
			}

			return
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:      "new",
			UsageText: "projects new <PROJECT_NAME>",
			Aliases:   []string{"n"},
			Usage:     "Create a new project",
			Action: func(c *cli.Context) (err error) {
				pn := c.Args().First()

				if pn == "" {
					return nil
				}
				err = os.Mkdir(cfg.current.Directory+"/"+pn, 0700)
				if os.IsExist(err) {
					return cli.Exit(fmt.Sprintf(`Project "%s" already exists`, pn), 0)
				}

				exec.Command("code", pn).Run()

				return
			},
		},
		{
			Name:      "add",
			Aliases:   []string{"a"},
			UsageText: "projects add <WORKSPACE_NAME>",
			Usage:     "Add current directory into workspaces with the given name",
			Before: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return cli.Exit("Workspace name must be provided!", 0)
				}

				return nil
			},
			Action: func(c *cli.Context) (err error) {
				dir, err := os.Getwd()
				if err != nil {
					return cli.Exit(err, 0)
				}

				if err = cfg.addWorkspace(Workspace{
					Name:      c.Args().First(),
					Directory: dir,
				}); err != nil {
					return cli.Exit(err, 0)
				}

				return
			},
		},
		{
			Name:      "remove",
			Aliases:   []string{"r", "rm"},
			UsageText: "projects remove <WORKSPACE_NAME>",
			Usage:     "Remove workspace",
			Before: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return cli.Exit("Workspace name must be provided!", 0)
				}

				return nil
			},
			Action: func(c *cli.Context) (err error) {
				wn := c.Args().First()
				if cfg.removeWorkspace(Workspace{Name: wn}) {
					fmt.Println("✔️")
				} else {
					fmt.Println("❌")
				}

				return
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d", "del"},
			UsageText: "projects delete <PROJECT_NAME>",
			Usage:     "Remove specific project from the current workspace",
			Before: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return cli.Exit("Project name must be provided!", 0)
				}

				return nil
			},
			Action: func(c *cli.Context) (err error) {
				pn := c.Args().First()
				var a string
				fmt.Printf(`This action will delete the project "%s"
Do you want to proceed? (Y/n)
> `, pn)
				if _, err := fmt.Scan(&a); err == nil && strings.EqualFold(a, "y") {
					if cfg.deleteProject(pn) {
						fmt.Println("✔️")
					} else {
						fmt.Println("❌")
					}
				}

				return
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"l"},
			UsageText: "projects list",
			Usage:     "List existing projects for the current workspace",
			Action: func(c *cli.Context) (err error) {
				var w *Workspace
				wn := c.Args().First()
				if wn == "" && cfg.current == nil {
					return cli.Exit("You must provide a workspace name or set a default workspace", 0)
				} else if wn != "" {
					w = cfg.getWorkspace(wn)
				} else {
					w = cfg.current
				}

				if w == nil {
					return cli.Exit(fmt.Sprintf(`Workspace "%s" is non existant`, wn), 0)
				}

				files, err := ioutil.ReadDir(w.Directory)
				if err != nil {
					return cli.Exit("Error: "+err.Error(), 0)
				}

				var fls []string
				for _, f := range files {
					if f.IsDir() {
						fls = append(fls, f.Name())
					}
				}

				const padding = 20
				writer := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
				var n int

				if len(fls) > 1 {
					for n+1 < len(fls) {
						fmt.Fprintf(writer, "%s\t%s\n", files[n].Name(), files[n+1].Name())
						n += 2
					}
					writer.Flush()
				} else {
					fmt.Println(fls[0])
				}

				return
			},
		},
		{
			Name:      "set",
			Aliases:   []string{"s"},
			UsageText: "projects set <WORKSPACE_NAME>",
			Usage:     "Set a default workspace",
			Before: func(c *cli.Context) error {
				if c.Args().First() == "" {
					return cli.Exit("Workspace name must be provided!", 0)
				}

				return nil
			},
			Action: func(c *cli.Context) (err error) {
				wn := c.Args().First()
				if cfg.changeCurrentWorkspace(wn) {
					fmt.Println("✔️")
				} else {
					fmt.Println("❌")
				}

				return
			},
		},
		{
			Name:      "workspaces",
			Aliases:   []string{"w"},
			UsageText: "projects workspaces",
			Usage:     "List existing workspaces",
			Action: func(c *cli.Context) (err error) {
				for _, w := range cfg.Workspaces {
					fmt.Printf("%s - %s\n\n", w.Name, w.Directory)
				}

				return
			},
		},
	}

	if cfg.current != nil {
		app.Description = "Current workspace - " + cfg.current.Name
	}

	return app.Run(os.Args)
}
