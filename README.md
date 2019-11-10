# concurrent test
```
go build -o curr -ldflags "-s -w" main.go

upx curr

fyne package -os darwin -icon icon.png

```