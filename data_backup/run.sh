#!/bin/bash

# 变量定义（替换为实际值）
FILE_TO_BACKUP="/root/ezbookkeeping/data/ezbookkeeping.db"  # 要备份的文件
REPO_DIR="/root/ezbookkeeping_data_git"           # 本地仓库路径
COMMIT_MSG="Automated backup: $(date +'%Y-%m-%d %H:%M:%S')"  # 提交消息

# 进入仓库目录
cd "$REPO_DIR" || exit 1

# 复制文件到仓库（假设仓库根目录下备份为 backup_file.txt）
cp "$FILE_TO_BACKUP" ./data_backup/ezbookkeeping.db

# Git 操作
git add ezbookkeeping.db
git commit -m "$COMMIT_MSG"
git push origin data_backup  # 修改为 data_backup 分支

# 如果使用 PAT（非 SSH），需设置环境变量（在 cron 前导出）
# export GITHUB_TOKEN=your_pat_token
# git push https://$GITHUB_TOKEN@github.com/your_username/your-repo.git data_backup
