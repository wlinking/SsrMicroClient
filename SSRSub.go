// +build !noGui

package main

import (
	"log"
	"os"

	"SsrMicroClient/gui"
	"SsrMicroClient/init"
	"SsrMicroClient/process/lockfile"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	configPath := ssrinit.GetConfigAndSQLPath()
	ssrinit.Init(configPath)

	lockFile, err := os.Create(configPath + "/SsrMicroClientRunStatuesLockFile")
	if err != nil {
		log.Println(err)
		return
	}
	if err = lockfile.LockFile(lockFile); err != nil {
		log.Println("process is exist!\n" + err.Error())
		return
	}
	defer func() {
		_ = lockFile.Close()
		_ = os.Remove(configPath + "/SsrMicroClientRunStatuesLockFile")
	}()

	ssrMicroClientGUI, err := gui.NewSsrMicroClientGUI(configPath)
	if err != nil && ssrMicroClientGUI != nil {
		ssrMicroClientGUI.MessageBox(err.Error())
	}
	if ssrMicroClientGUI != nil {
		if ssrMicroClientGUI.App.IsSessionRestored() {
			ssrMicroClientGUI.MessageBox("restore is from before")
		}

		//ssrMicroClientGUI.MainWindow.Show()
		ssrMicroClientGUI.App.Exec()
	}
}
