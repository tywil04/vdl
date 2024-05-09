package main

import (
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"vdl/internal/pkg/distrobox"
	"vdl/internal/pkg/vscode"
)

const appName = "uk.co.tylerw.vdl"

//go:embed ui.ui
var ui string

func main() {
	app := gtk.NewApplication(appName, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		builder := gtk.NewBuilderFromString(ui, len(ui))
		window := builder.GetObject("main-window").Cast().(*adw.Window)

		app.AddWindow(&window.Window)

		containersCombo := builder.GetObject("containers-row").Cast().(*adw.ComboRow)
		directoryOpen := builder.GetObject("directory-open-button").Cast().(*gtk.Button)
		directoryEntry := builder.GetObject("directory-row").Cast().(*adw.EntryRow)
		startButton := builder.GetObject("start-button").Cast().(*gtk.Button)
		cancelButton := builder.GetObject("cancel-button").Cast().(*gtk.Button)

		containers, err := distrobox.GetContainers()
		if err != nil {
			panic(err)
		}
		containersComboList := gtk.NewStringList(containers)
		containersCombo.SetModel(containersComboList)

		directoryOpen.Connect("clicked", func() {
			directoryDialog := gtk.NewFileChooserNative("Select a Directory", &window.Window, gtk.FileChooserActionSelectFolder, "Select", "Cancel")
			directoryDialog.Connect("response", func() {
				directoryEntry.SetText(directoryDialog.File().Path())
			})
			directoryDialog.Show()
		})

		startButton.Connect("clicked", func() {
			go func() {
				selectedContainer := containersCombo.SelectedItem().Cast().(*gtk.StringObject).String()
				selectedDirectory := directoryEntry.Text()

				if selectedContainer == "" || strings.TrimSpace(selectedDirectory) == "" {
					log.Println("container or directory not selected")
					return
				}

				window.Hide()

				containerId := distrobox.GetContainerId(selectedContainer)

				err := distrobox.StartContainer(selectedContainer)
				if err != nil {
					log.Println("distrobox ", err)
				}

				err = vscode.OpenContainer(containerId, selectedDirectory)
				if err != nil {
					log.Println("vscode ", err)
				}

				os.Exit(0)
			}()
		})

		cancelButton.Connect("clicked", func() {
			os.Exit(0)
		})

		window.Show()
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
