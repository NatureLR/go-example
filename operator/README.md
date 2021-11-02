# operator开发

```shell
kubebuilder init --repo github.com/naturelr/code-example/operator --domain naturelr.cc --skip-go-version-check

# 创建 api
kubebuilder create api --group appx --version v1 --kind Appx

# 创建webhook
kubebuilder create webhook --group nodes --version v1 --kind Appx --defaulting --programmatic-validation

# 生成文件
make manifests generate

# 安装crd等文件
make install

# 本地调试运行
make run
```
