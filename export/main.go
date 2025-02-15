package main

import (
	"fmt"
	"github.com/orestonce/go2cpp"
	"m3u8d"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	BuildCliBinary()   // 编译二进制
	CreateLibForQtUi() // 创建Qt需要使用的.a库文件
}

func BuildCliBinary() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	type buildCfg struct {
		GOOS   string
		GOARCH string
	}
	var list = []buildCfg{
		{
			GOOS:   "linux",
			GOARCH: "386",
		},
		{
			GOOS:   "linux",
			GOARCH: "arm",
		},
		{
			GOOS:   "linux",
			GOARCH: "mipsle",
		},
		{
			GOOS:   "darwin",
			GOARCH: "amd64",
		},
	}
	for _, cfg := range list {
		name := "m3u8d_cli_v1.5.0_" + cfg.GOOS + "_" + cfg.GOARCH
		cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", filepath.Join(wd, "bin", name))
		cmd.Dir = filepath.Join(wd, "cmd")
		cmd.Env = append(os.Environ(), "GOOS="+cfg.GOOS)
		cmd.Env = append(cmd.Env, "GOARCH="+cfg.GOARCH)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(cmd.Dir)
			panic(err)
		}
		fmt.Println("done", name)
	}
}

func CreateLibForQtUi() {
	ctx := go2cpp.NewGo2cppContext(go2cpp.NewGo2cppContext_Req{
		CppBaseName:                 "m3u8d",
		EnableQtClass_RunOnUiThread: true,
		EnableQtClass_Toast:         true,
	})
	ctx.Generate1(m3u8d.RunDownload)
	ctx.Generate1(m3u8d.CloseOldEnv)
	ctx.Generate1(m3u8d.GetProgress)
	ctx.Generate1(m3u8d.GetWd)
	ctx.Generate1(m3u8d.ParseCurlStr)
	ctx.Generate1(m3u8d.RunDownload_Req_ToCurlStr)
	ctx.MustCreateAmd64LibraryInDir("m3u8d-qt")
}
