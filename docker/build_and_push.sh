bazel run //:goto_image -- --norun &&
  docker tag bazel:goto_image linmounong/goto &&
  docker push linmounong/goto
