# bp

bp 是 buildpack 缩写，包括 build、pack 两个步骤。

## build

`build` means `build your code`.

`build` 目录下分为多个子目录，每个子目录代表一种 build 方式。

每个子目录需要提供 `Dockerfile` 文件及其依赖文件，action 会执行 `Dockerfile` 来编译应用。

## pack

`pack` means `pack to image`.

`pack` 目录和 `build` 一样，也分为多个子目录，每个子目录代表一种 pack 方式。

每个子目录需要提供 `Dockerfile` 文件及其依赖文件，action 会执行 `Dockerfile` 来制作业务镜像。
