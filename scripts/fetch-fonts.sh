#!/usr/bin/env bash
# 安装 Noto Sans SC 到 ~/.local/share/fonts（H5 截图、Android 渲染、设计 fixture 都用这套字体）
set -euo pipefail
tmp=$(mktemp -d)
cd "$tmp"
npm pack @expo-google-fonts/noto-sans-sc --silent >/dev/null
tar xzf ./*.tgz
mkdir -p ~/.local/share/fonts
for w in 400Regular 500Medium 600SemiBold 700Bold; do cp "package/$w/NotoSansSC_$w.ttf" ~/.local/share/fonts/; done
fc-cache -f >/dev/null 2>&1 || true
rm -rf "$tmp"
echo "✓ Noto Sans SC installed"
