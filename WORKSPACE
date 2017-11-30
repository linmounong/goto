http_archive(
    name = "io_bazel_rules_go",
    sha256 = "341d5eacef704415386974bc82a1783a8b7ffbff2ab6ba02375e1ca20d9b031c",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.1/rules_go-0.7.1.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.3.0",
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

go_repository(
    name = "com_github_jinzhu_gorm",
    commit = "0a51f6cdc55d1650d9ed3b4c13026cfa9133b01e",
    importpath = "github.com/jinzhu/gorm",
)

go_repository(
    name = "com_github_jinzhu_inflection",
    commit = "1c35d901db3da928c72a72d8458480cc9ade058f",
    importpath = "github.com/jinzhu/inflection",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    commit = "d5ffb5c0cca8778699a929b236766f4a7af674e8",
    importpath = "github.com/mattn/go-sqlite3",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    commit = "30f82fa23fd844bd5bb1e5f216db87fd77b5eb43",
    importpath = "github.com/elazarl/go-bindata-assetfs",
)

go_repository(
    name = "com_github/jteeuwen/go_bindata",
    commit = "a0ff2567cfb70903282db057e799fd826784d41d",
    importpath = "github.com/jteeuwen/go-bindata",
)
