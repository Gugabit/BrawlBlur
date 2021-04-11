# BrawlBlur

- Generate .syso command:

```sh
rsrc.exe -manifest ".\BrawlBlur.exe.manifest" -ico ".\favicon.ico" -o "BrawlBlur.syso"
```

- Build command:

```sh
go build -o BrawlBlur.exe -ldflags -H=windowsgui
```
