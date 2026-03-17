@echo off
chcp 65001 >nul
REM 将代码从 GitHub 同步推送到 Gitee

set GITEE_REMOTE=gitee
set GITEE_URL=git@gitee.com:sgcnpm/gra-api.git

REM 检查 gitee remote 是否已存在，不存在则添加
git remote get-url %GITEE_REMOTE% >nul 2>&1
if %errorlevel% neq 0 (
    echo 添加 Gitee remote...
    git remote add %GITEE_REMOTE% %GITEE_URL%
)

REM 获取当前分支名
for /f "tokens=*" %%b in ('git symbolic-ref --short HEAD') do set BRANCH=%%b

echo 正在推送分支 [%BRANCH%] 到 Gitee...
git push %GITEE_REMOTE% %BRANCH% --tags

if %errorlevel% equ 0 (
    echo 推送成功！
) else (
    echo 推送失败，请检查网络或 SSH 密钥配置。
    exit /b 1
)
