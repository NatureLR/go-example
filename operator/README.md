# 初始化项目
kubebuilder init --repo github.com/naturelr/code-example/operator --domain naturelr.cc --skip-go-version-check

# 创建 api
kubebuilder create api --group appx --version v1 --kind Appx

# 生成文件
make manifests generate

# 安装crd等文件
make install

# 本地调试运行
make run
