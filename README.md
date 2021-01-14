# adboverssh

Remote control Android device.

On your Android device:

```
# build and upload to Android device
GOOS=linux GOARCH=arm64 go build -v ./cmd/adboverssh && adb push ./adboverssh /data/local/tmp

# you may need to generate and upload a key
ssh-keygen -f id_rsa
adb push id_rsa /data/local/tmp

# run
adb shell

# check if adb port is 5555, use -adb option if you have different port
getprop service.adb.tcp.port

cd /data/local/tmp
./adboverssh -i id_rsa root@1.2.3.4:22

# if you don't specify -l option, random port will be used
# 2021/01/14 05:49:48 connected to 1.2.3.4:22
# 2021/01/14 05:49:48 listening 127.0.0.1:35651
```

On your SSH server:

```
# download adb on https://developer.android.com/studio/releases/platform-tools
adb connect 127.0.0.1:35651
adb shell
```

You can run `go run ./cmd/generatekey` to generate `mobile/key.go` file.
