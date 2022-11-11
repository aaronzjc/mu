package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/aaronzjc/mu/pkg/flow"
)

const (
	AppVersion = "8.0"
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
	npm = npm.WithEnvVariable("APP_VERSION", AppVersion)
	npm = npm.WithMountedDirectory("/src/web", src).WithWorkdir("/src/web")
	npm = npm.WithEnvVariable("VERSION", AppVersion)
	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"cat", ".env"},
	})
	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "install", "--sass_binary_site=https://npm.taobao.org/mirrors/node-sass/"},
	})
	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "run", "build"},
	})
	build, err := npm.Stdout().Contents(ctx)
	if err != nil {
		return err
	}
	if _, err := npm.Directory("/src/public").Export(ctx, "dagger/frontend"); err != nil {
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
	src := client.Host().Workdir()
	golang := client.Container().From("golang:1.19-alpine3.15")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")
	golang = golang.WithEnvVariable("GO111MODULE", "on")
	golang = golang.WithEnvVariable("GOPROXY", "https://goproxy.cn,direct")
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")
	path := "dagger/backend/"
	for _, target := range []string{"api", "commander", "agent"} {
		golang = golang.Exec(dagger.ContainerExecOpts{
			Args: []string{
				"go",
				"build",
				"-ldflags", "-X main.version=" + AppVersion,
				"-o", path + target,
				"cmd/" + target + "/main.go"},
		})
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

	src := client.Host().Workdir()
	for _, v := range []string{"api", "commander", "agent"} {
		docker := client.Container()
		docker = docker.Build(src, dagger.ContainerBuildOpts{Dockerfile: "./scripts/dockerfiles/" + v + ".Dockerfile"})
		resp, err := docker.Publish(ctx, fmt.Sprintf("aaronzjc/mu-%s:%s", v, AppVersion))
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

	kubeconfig := client.Host().Workdir().File("./scripts/kubeconf.yaml")
	deployments := client.Host().Workdir().Directory("./scripts/k8s")
	kubectl := client.Container().From("bitnami/kubectl")
	kubectl = kubectl.WithMountedFile("/.kube/config", kubeconfig)
	kubectl = kubectl.WithMountedDirectory("/tmp", deployments)
	kubectl = kubectl.Exec(dagger.ContainerExecOpts{
		Args: []string{"apply", "-f", "/tmp/mu-api.yaml", "-n", "k3s-apps"},
	})
	kubectl = kubectl.Exec(dagger.ContainerExecOpts{
		Args: []string{"apply", "-f", "/tmp/mu-agent.yaml", "-n", "k3s-apps"},
	})
	kubectl = kubectl.Exec(dagger.ContainerExecOpts{
		Args: []string{"apply", "-f", "/tmp/mu-commander.yaml", "-n", "k3s-apps"},
	})
	logs, err := kubectl.Stdout().Contents(ctx)
	if err != nil {
		return err
	}
	fmt.Println(logs)
	return nil
}
