load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_binary", "go_library", "go_test")
load("@io_bazel_rules_go//examples/bindata:bindata.bzl", "bindata")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")

gazelle(
    name = "gazelle",
    prefix = "github.com/linmounong/goto",
)

bindata(
    name = "s",
    srcs = glob(["s/**"]),
    package = "main",
)

go_library(
    name = "go_default_library",
    srcs = [
        "db.go",
        "handlers.go",
        "http_protocol.go",
        "main.go",
        "utils.go",
        ":s",
    ],
    importpath = "github.com/linmounong/goto",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_elazarl_go_bindata_assetfs//:go_default_library",
        "@com_github_jinzhu_gorm//:go_default_library",
        "@com_github_jinzhu_gorm//dialects/sqlite:go_default_library",
    ],
)

go_binary(
    name = "goto",
    embed = [":go_default_library"],
    importpath = "github.com/linmounong/goto",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "db_test.go",
        "handlers_test.go",
    ],
    embed = [":go_default_library"],
    importpath = "github.com/linmounong/goto",
)

go_image(
    name = "goto_image",
    embed = [":go_default_library"],
    importpath = "github.com/linmounong/goto",
    visibility = ["//visibility:public"],
)
