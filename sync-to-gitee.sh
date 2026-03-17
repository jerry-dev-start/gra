#!/bin/bash
# 将代码从 GitHub 同步推送到 Gitee

GITEE_REMOTE="gitee"
GITEE_URL="git@gitee.com:sgcnpm/gra-api.git"

# 检查 gitee remote 是否已存在，不存在则添加
if ! git remote get-url "$GITEE_REMOTE" &>/dev/null; then
    echo "添加 Gitee remote..."
    git remote add "$GITEE_REMOTE" "$GITEE_URL"
fi

# 获取当前分支名
BRANCH=$(git symbolic-ref --short HEAD)

echo "正在推送分支 [$BRANCH] 到 Gitee..."
git push "$GITEE_REMOTE" "$BRANCH" --tags

if [ $? -eq 0 ]; then
    echo "推送成功！"
else
    echo "推送失败，请检查网络或 SSH 密钥配置。"
    exit 1
fi
