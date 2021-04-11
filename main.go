/*
 * Author: Gugabit
 * Website: gugabit.com
 * Date: 11/04/2021
 */

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"hawx.me/code/img/blur"
)

func main() {
	a := app.New()
	w := a.NewWindow("BrawlBlur Wallpapers")

	titleMessage := widget.NewLabelWithStyle("Please select Brawlhalla root folder\nUsually located at: C:\\Program Files (x86)\\Steam\\steamapps\\common\\Brawlhalla", fyne.TextAlignCenter, fyne.TextStyle{})
	blurLevel := binding.NewFloat()
	blurLevel.Set(15)
	progress := binding.NewFloat()
	bgFiles := []os.FileInfo{}
	bgBackupFiles := []os.FileInfo{}
	backgroundsPath := ""
	backgroundsBackupPath := ""
	backgroundsFolderIsEmpty := false
	backgroundsBackupFolderExists := false
	backgroundsBackupFolderIsEmpty := true
	ggbWebsite, _ := url.Parse("https://gugabit.com")

	blurButton := widget.NewButton("Blur Backgrounds", func() {
		bgFiles, _ = readDir(backgroundsPath)
		if len(Filter(bgFiles, func(file os.FileInfo) bool { return file.Name() != "Backup" })) == 0 {
			backgroundsFolderIsEmpty = true
		} else {
			backgroundsFolderIsEmpty = false
		}

		if _, err := os.Stat(backgroundsBackupPath); os.IsNotExist(err) {
			backgroundsBackupFolderExists = false
		} else {
			backgroundsBackupFolderExists = true
		}

		if backgroundsBackupFolderExists {
			bgBackupFiles, _ = readDir(backgroundsBackupPath)
			if len(bgBackupFiles) != 0 {
				backgroundsBackupFolderIsEmpty = false
			} else {
				backgroundsBackupFolderIsEmpty = true
			}
		}

		if !backgroundsBackupFolderExists {
			os.Mkdir(backgroundsBackupPath, os.FileMode(0777))
		}

		if backgroundsBackupFolderIsEmpty && !backgroundsFolderIsEmpty {
			bgFiles, _ = readDir(backgroundsPath)
			for _, file := range bgFiles {
				os.Rename(fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()), fmt.Sprintf("%s%c%s", backgroundsBackupPath, os.PathSeparator, file.Name()))
			}
		}

		if backgroundsBackupFolderExists && !backgroundsBackupFolderIsEmpty {
			bgFiles, _ = readDir(backgroundsPath)
			for _, file := range bgFiles {
				if !file.IsDir() {
					os.Remove(fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()))
				}
			}
		}

		bgBackupFiles, _ = readDir(backgroundsBackupPath)
		nFiles := len(bgBackupFiles)
		dialog.Message("Blurring %d images, this might take a while, please wait.", nFiles).Info()

		for index, file := range bgBackupFiles {
			progress.Set(float64(index) / float64(nFiles))
			blurImage(fmt.Sprintf("%s%c%s", backgroundsBackupPath, os.PathSeparator, file.Name()), fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()), blurLevel)
		}
		progress.Set(1)
		dialog.Message("Wallpapers blurred successfully!").Title("SUCCESS!!!").Info()
	})
	blurButton.Disable()

	undoBlurButton := widget.NewButton("Unblur Backgrounds", func() {
		if _, err := os.Stat(backgroundsBackupPath); os.IsNotExist(err) {
			backgroundsBackupFolderExists = false
		} else {
			backgroundsBackupFolderExists = true
		}

		if backgroundsBackupFolderExists {
			bgBackupFiles, _ = readDir(backgroundsBackupPath)
			if len(bgBackupFiles) != 0 {
				backgroundsBackupFolderIsEmpty = false
			} else {
				backgroundsBackupFolderIsEmpty = true
			}
		}

		continueDialog := dialog.Message("Are you sure you want to undo the changes?").YesNo()
		if continueDialog {
			bgFiles, _ = readDir(backgroundsPath)

			if backgroundsBackupFolderExists && !backgroundsBackupFolderIsEmpty {
				for _, file := range bgFiles {
					if !file.IsDir() {
						os.Remove(fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()))
					}
				}
				bgBackupFiles, _ = readDir(backgroundsBackupPath)
				for _, file := range bgBackupFiles {
					os.Rename(fmt.Sprintf("%s%c%s", backgroundsBackupPath, os.PathSeparator, file.Name()), fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()))
				}

				os.Remove(backgroundsBackupPath)
				dialog.Message("Blur undone!").Title("SUCCESS!!!").Info()
				return
			}
			continueDeleteDialog := dialog.Message("Backup folder not found, you wish to delete all backgrounds instead?").YesNo()
			if continueDeleteDialog {
				for _, file := range bgFiles {
					if !file.IsDir() {
						os.Remove(fmt.Sprintf("%s%c%s", backgroundsPath, os.PathSeparator, file.Name()))
					}
				}
				dialog.Message("Please verify game files on steam.").Info()
			}
		}
	})
	undoBlurButton.Disable()

	w.SetContent(container.NewVBox(
		titleMessage,
		widget.NewButton("Select Folder", func() {
			brawlhallaPath, _ := dialog.Directory().Title("Now find a dir").Browse()

			backgroundsFolderIsEmpty = false
			backgroundsBackupFolderExists = false
			backgroundsBackupFolderIsEmpty = true
			fmt.Println(brawlhallaPath)

			backgroundsPath = fmt.Sprintf("%s%cmapArt%cBackgrounds", brawlhallaPath, os.PathSeparator, os.PathSeparator)
			backgroundsBackupPath = fmt.Sprintf("%s%cBackup", backgroundsPath, os.PathSeparator)

			if _, err := os.Stat(backgroundsPath); os.IsNotExist(err) {
				titleMessage.SetText("Backgrounds folder not found.")
				return
			}

			titleMessage.SetText(fmt.Sprintf("Backgrounds folder selected:\n%s", backgroundsPath))

			bgFiles, _ = readDir(backgroundsPath)
			if len(Filter(bgFiles, func(file os.FileInfo) bool { return file.Name() != "Backup" })) == 0 {
				backgroundsFolderIsEmpty = true
			} else {
				backgroundsFolderIsEmpty = false
			}

			if _, err := os.Stat(backgroundsBackupPath); os.IsNotExist(err) {
				backgroundsBackupFolderExists = false
			} else {
				backgroundsBackupFolderExists = true
			}

			if backgroundsBackupFolderExists {
				bgBackupFiles, _ = readDir(backgroundsBackupPath)
				if len(bgBackupFiles) != 0 {
					backgroundsBackupFolderIsEmpty = false
				} else {
					backgroundsBackupFolderIsEmpty = true
				}
			}

			if backgroundsFolderIsEmpty {
				if !backgroundsBackupFolderExists {
					titleMessage.SetText("Backgrounds folder is empty.")
				}
				if backgroundsBackupFolderIsEmpty {
					titleMessage.SetText("Backgrounds and Backup folders are empty.")
				}
			}

			blurButton.Enable()
			undoBlurButton.Enable()
		}),
		widget.NewLabelWithStyle("Blur Level\n(The greater, more computational power will be required. thefore more time.)", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewSliderWithData(1, 30, blurLevel),
		blurButton,
		widget.NewLabelWithStyle("Progress", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewProgressBarWithData(progress),
		undoBlurButton,
		widget.NewHyperlinkWithStyle("Made with Love by: Gugabit", ggbWebsite, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("Toss A Coin To Your Dev\n\nBTC: bc1qg332s70ga9qnz5hxvwgme075pcy22wucjv8a2c", fyne.TextAlignCenter, fyne.TextStyle{}),
	))

	w.ShowAndRun()
}

func blurImage(imgPath string, imageDestination string, blurLevelBind binding.Float) {
	imagePath, _ := os.Open(imgPath)
	defer imagePath.Close()
	srcImage, _, _ := image.Decode(imagePath)

	blurLevel, _ := blurLevelBind.Get()
	blurred := blur.Box(srcImage, int(blurLevel), blur.CLAMP)

	newImage, _ := os.Create(imageDestination)
	defer newImage.Close()
	jpeg.Encode(newImage, blurred, &jpeg.Options{Quality: 80})
}

func readDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	list, err := f.Readdir(-1)

	if err != nil {
		return nil, err
	}
	return list, nil
}

func Filter(arr []os.FileInfo, cond func(os.FileInfo) bool) []os.FileInfo {
	result := []os.FileInfo{}
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}
