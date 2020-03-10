package main

import "os"

// IsCakeScreenshot return true if imgName is a "Cake" screenshot
func IsCakeScreenshot(imgName string) (bool, int) {
	if _, err := os.Stat(imgName); os.IsNotExist(err) {
		return false, 0
	}
	return true, 0
}
