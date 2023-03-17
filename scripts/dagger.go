package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/aaronzjc/mu/pkg/flow"
)

const (
	Version = "8.14"
)

func main() {
	var action string
	if len(os.Args) >= 2 {
		action = os.Args[1]
	}
	fw := flow.NewFlow(context.Background())
	switch action {
	case "backend":
		fw.Step("backend", buildBackend)
	case "frontend":
		fw.Step("frontend", buildFrontend)
	case "image":
		fw.Step("backend", buildBackend)
		fw.Step("frontend", buildFrontend)
		fw.Step("image", buildAndPushImage)
	case "deploy":
		fw.Step("backend", buildBackend)
		fw.Step("frontend", buildFrontend)
		fw.Step("image", buildAndPushImage)
		fw.Step("deploy", deploy)
	default:
		fmt.Println("usage: go run ./scripts/dagger.go [backend|frontend|image|deploy]")
		return
	}
	fw.Run()
}

func buildFrontend(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Directory("./web", dagger.HostDirectoryOpts{
		Exclude: []string{"node_modules"},
	})
	npm := client.Container().From("node:14-alpine")
	npm = npm.WithMountedDirectory("/src/web", src).WithWorkdir("/src/web")
	npm = npm.WithEnvVariable("VERSION", Version)
	npm = npm.WithExec([]string{"cat", ".env"})
	npm = npm.WithExec([]string{"npm", "install", "--sass_binary_site=https://npm.taobao.org/mirrors/node-sass/"})
	npm = npm.WithExec([]string{"npm", "run", "build"})
	build, err := npm.Stdout(ctx)
	if err != nil {
		return err
	}

	dst := "dagger/frontend"
	// 清空目标输出目录
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if _, err := npm.Directory("/src/public").Export(ctx, dst); err != nil {
		return err
	}
	// just for dev
	if _, err := npm.Directory("/src/public").Export(ctx, "public"); err != nil {
		return err
	}
	fmt.Println("npm stdout", build)
	return nil
}

func buildBackend(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// 获取本地项目路径
	src := client.Host().Directory(".")
	golang := client.Container().From("golang:1.19-alpine3.15")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")
	golang = golang.WithEnvVariable("GO111MODULE", "on")
	golang = golang.WithEnvVariable("GOPROXY", "https://goproxy.cn,direct")
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")
	path := "dagger/backend/"
	for _, target := range []string{"api", "commander", "agent"} {
		golang = golang.WithExec([]string{
			"go",
			"build",
			"-ldflags", "-X main.version=" + Version,
			"-o", path + target,
			"cmd/" + target + "/main.go"},
		)
	}
	if _, err := golang.Directory(path).Export(ctx, path); err != nil {
		return err
	}

	return nil
}

func buildAndPushImage(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Directory(".")
	for _, v := range []string{"api", "commander", "agent"} {
		docker := client.Container()
		docker = docker.Build(src, dagger.ContainerBuildOpts{Dockerfile: "./scripts/dockerfiles/" + v + ".Dockerfile"})
		resp, err := docker.Publish(ctx, fmt.Sprintf("aaronzjc/mu-%s:%s", v, Version))
		if err != nil {
			return err
		}
		fmt.Println(resp)
	}

	return nil
}

func deploy(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// 更新部署的版本
	fileData := make(map[string][]byte)
	oldTag, newTag := "latest", Version
	dir := "./scripts/k8s"
	files := []string{"mu-api.yaml", "mu-agent.yaml", "mu-commander.yaml"}
	for _, v := range files {
		file := dir + "/" + v
		if fileData[file], err = os.ReadFile(file); err != nil {
			return err
		}
		out := strings.ReplaceAll(string(fileData[file]), oldTag, newTag)
		os.WriteFile(file, []byte(out), 0666)
	}
	defer func() {
		for _, v := range files {
			file := dir + "/" + v
			os.WriteFile(file, fileData[file], 0666)
		}
	}()

	kubectl := client.Container().From("bitnami/kubectl")
	kubeconfig := client.Host().Directory(".").File("./scripts/kubeconf.yaml")
	kubectl = kubectl.WithMountedFile("/.kube/config", kubeconfig)
	deployments := client.Host().Directory(dir)
	kubectl = kubectl.WithMountedDirectory("/tmp", deployments)
	for _, f := range files {
		kubectl = kubectl.WithExec([]string{"apply", "-f", "/tmp/" + f, "-n", "k3s-apps"})
	}
	logs, err := kubectl.Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println(logs)
	return nil
}
