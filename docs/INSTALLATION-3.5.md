# 安装步骤

### 克隆代码库

### 切换至 Tag: v3.5

### 切换至工作目录 releases/3.5

已经包含 actions-private, actions-public, addons-private 和 addons-public 四个文件

文件按行处理，每行内容为：name@version，例如 mysql@5.7.23

### 导出 (export) 所有扩展到当前目录的 backup 目录下

执行 `bash releases/ext_export.sh ${dice_version}`

- 例如：`bash releases/ext_export.sh 3.5`
- 数据流向：远端扩展市场元数据 -> 本地 backup 目录 (元信息 + 镜像 blob)

### 将 backup 目录和 releases 目录(脚本) 通过某种方式传输至私有集群内

### 导入 (import) 所有扩展，包括镜像推送和市场元数据更新(ext push)
  - `bash releases/ext_import.sh ${dice_version} ${destRegistryHost}`
  - 例如：`bash releases/ext_import.sh 3.5 localhost:5000`
  - 数据流向：本地 backup 目录 (元信息 + 镜像 blob) -> 镜像推送到目标 registry，元数据推送到远端扩展市场
