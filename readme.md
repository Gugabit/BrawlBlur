# BrawlBlur

Utility for blurring Brawlhalla backgrounds.

## Use instructions

1. Download BrawlBlur.exe release [here](https://github.com/Gugabit/BrawlBlur/releases/download/v0.1/BrawlBlur.exe)
2. Open the program
3. Click the button 'Select Folder' and select Brawlhalla folder (not wallpapers folder).
4. Select the blur level in the slider.
5. Click the Blur Backgrounds button.

> Done!

To undo the changes simply click on 'Unblur backgrounds', or copy the files from the backup folder to the Backgrounds folder.

---

## Build instructions

- Generate .syso command:

```sh
rsrc.exe -manifest ".\BrawlBlur.exe.manifest" -ico ".\favicon.ico" -o "BrawlBlur.syso"
```

- Build command:

```sh
go build -o BrawlBlur.exe -ldflags -H=windowsgui
```
