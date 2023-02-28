chcp 65001
@echo off
setlocal enabledelayedexpansion 

echo --------------------------------------
echo 	Build Script By Break
echo 	     QQ:441479573
echo --------------------------------------

SET FileName=Hello
set "prefix="

:main
	echo 请选择需要构建的平台
	echo.
	echo (1) Centos
	echo (2) Windows
	echo (3) OpenWRT
	echo (4) Android
	echo (5) MacOSX
	echo (0) Exit
	echo.
	set /p type=请输入对应数字:

    echo -------------------- x --------------------
    IF EXIST dist (
        rmdir /s/q dist
    )

	goto do%type%
	if %errorlevel%==1 goto main

:do1
	echo 正在构建 Centos 平台程序...
	go build tests\zip.go
	set CGO_ENABLED=0
	set GOOS=linux
	set GOARCH=amd64
	goto done

:do2
	echo 正在构建 Windows 平台程序...
	go build tests\zip.go
	set "prefix=.exe"
	goto done

:do3
	echo 正在构建 OpenWRT 平台程序...
	go build tests\zip.go
	set GOOS=linux
	set GOARCH=arm
	goto done

:do4
	echo 正在构建 Android 平台程序...
	go build tests\zip.go
	set CGO_ENABLED=0 GOARCH=arm GOOS=linux
	goto done
	
:do5
	echo 正在构建 MacOSX 平台程序...
	go build tests\zip.go
	SET CGO_ENABLED=0
	SET GOOS=darwin
	SET GOARCH=amd64
	goto done
	
:done
    :: 构建目录
    mkdir dist
    MD dist\config
    MD dist\resources
    MD dist\storage\cache
    MD dist\storage\logs

    :: 编译
    go build -o %FileName%%prefix% main.go

    :: 移动&复制文件
    move cli_task%prefix% dist
    copy config\app.yaml.example dist\config
    copy storage\logs\.gitignore dist\storage\logs\.gitignore
    xcopy resources dist\resources /q/s
    ::move %FileName%* dist
    echo -------------------- x --------------------

    :: 压缩
	zip.exe dist dist.zip

    rmdir /s/q dist
    del zip.exe

	echo 构建成功

:do0
	pause

