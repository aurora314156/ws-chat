#!/bin/bash

# ------------------------------
# 一鍵安裝 WSL + zsh + mise + Python 3
# ------------------------------

set -e  # 發生錯誤就停止

echo "🚀 更新套件列表..."
sudo apt update
sudo apt install -y curl git unzip build-essential libssl-dev zlib1g-dev \
libbz2-dev libreadline-dev libsqlite3-dev wget llvm libncurses5-dev libncursesw5-dev \
xz-utils tk-dev libffi-dev liblzma-dev zsh

# 安裝 mise
echo "🚀 安裝 mise..."
curl -s https://mise.jdx.dev/install.sh | sh

# 設定 zsh
ZSHRC="$HOME/.zshrc"
echo "🚀 設定 zsh 與 PATH..."
grep -qxF 'export PATH="$HOME/.local/share/mise/shims:$PATH"' $ZSHRC || echo 'export PATH="$HOME/.local/share/mise/shims:$PATH"' >> $ZSHRC
grep -qxF 'eval "$(~/.local/bin/mise activate zsh)"' $ZSHRC || echo 'eval "$(~/.local/bin/mise activate zsh)"' >> $ZSHRC

# 重新載入 zsh
source $ZSHRC

# 安裝 Python 3
PYTHON_VERSION="3.12"
echo "🚀 安裝 Python $PYTHON_VERSION..."
mise install python@$PYTHON_VERSION
mise use -g python@$PYTHON_VERSION

# 驗證
echo "✅ 驗證 Python 與 pip"
python --version
pip --version

# 安裝 venv 模組
echo "🚀 安裝 venv..."
sudo apt install -y python3-venv

echo "🎉 安裝完成！"
echo "你可以用 'python -m venv myenv' 建立虛擬環境"
