#!/usr/bin/env bash
set -euo pipefail

# Использование: ./scripts/release.sh v0.1.0
VERSION=${1:?Укажи версию: ./scripts/release.sh v0.1.0}
BINARY=lcli
DIST=dist
REPO=pvarganov/lcli
TAP_REPO=pvarganov/homebrew-tap

# Проверяем наличие gh
if ! command -v gh &>/dev/null; then
  echo "Установи GitHub CLI: brew install gh"
  exit 1
fi

rm -rf "$DIST"
mkdir -p "$DIST"

echo "==> Сборка бинарников ($VERSION)"

build() {
  local os=$1 arch=$2
  local out="$DIST/${BINARY}_${os}_${arch}"
  echo "    $os/$arch"
  GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o "$out" .
  tar -czf "${out}.tar.gz" -C "$DIST" "${BINARY}_${os}_${arch}"
  rm "$out"
}

build darwin  amd64
build darwin  arm64
build linux   amd64
build linux   arm64

echo "==> Подсчёт SHA256"
SHA_DARWIN_AMD64=$(shasum -a 256 "$DIST/${BINARY}_darwin_amd64.tar.gz" | awk '{print $1}')
SHA_DARWIN_ARM64=$(shasum -a 256 "$DIST/${BINARY}_darwin_arm64.tar.gz" | awk '{print $1}')
SHA_LINUX_AMD64=$(shasum  -a 256 "$DIST/${BINARY}_linux_amd64.tar.gz"  | awk '{print $1}')
SHA_LINUX_ARM64=$(shasum  -a 256 "$DIST/${BINARY}_linux_arm64.tar.gz"  | awk '{print $1}')

echo "==> Создание GitHub Release $VERSION"
gh release create "$VERSION" \
  "$DIST/${BINARY}_darwin_amd64.tar.gz" \
  "$DIST/${BINARY}_darwin_arm64.tar.gz" \
  "$DIST/${BINARY}_linux_amd64.tar.gz" \
  "$DIST/${BINARY}_linux_arm64.tar.gz" \
  --repo "$REPO" \
  --title "$VERSION" \
  --generate-notes

echo "==> Генерация Homebrew Formula"
FORMULA=$(cat <<EOF
class Lcli < Formula
  desc "CLI for Linear issue tracker"
  homepage "https://github.com/pvarganov/lcli"
  version "${VERSION#v}"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/pvarganov/lcli/releases/download/${VERSION}/${BINARY}_darwin_arm64.tar.gz"
      sha256 "${SHA_DARWIN_ARM64}"
    end
    on_intel do
      url "https://github.com/pvarganov/lcli/releases/download/${VERSION}/${BINARY}_darwin_amd64.tar.gz"
      sha256 "${SHA_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pvarganov/lcli/releases/download/${VERSION}/${BINARY}_linux_arm64.tar.gz"
      sha256 "${SHA_LINUX_ARM64}"
    end
    on_intel do
      url "https://github.com/pvarganov/lcli/releases/download/${VERSION}/${BINARY}_linux_amd64.tar.gz"
      sha256 "${SHA_LINUX_AMD64}"
    end
  end

  def install
    bin.install "${BINARY}"
  end

  test do
    system "#{bin}/${BINARY}", "--version"
  end
end
EOF
)

# Обновляем формулу в tap-репозитории (если он склонирован рядом)
TAP_DIR="../homebrew-tap"
if [ -d "$TAP_DIR/Formula" ]; then
  echo "$FORMULA" > "$TAP_DIR/Formula/lcli.rb"
  git -C "$TAP_DIR" add Formula/lcli.rb
  git -C "$TAP_DIR" commit -m "lcli ${VERSION}"
  git -C "$TAP_DIR" push
  echo "==> Formula обновлена в $TAP_DIR"
else
  echo ""
  echo "==> Tap-репозиторий не найден рядом ($TAP_DIR)."
  echo "    Создай его: https://github.com/new → имя: homebrew-tap"
  echo "    Затем склонируй рядом с этим репо и запусти скрипт ещё раз."
  echo ""
  echo "==> Либо скопируй Formula вручную:"
  echo "---"
  echo "$FORMULA"
fi

echo "==> Готово! Установка:"
echo "    brew tap pvarganov/tap"
echo "    brew install lcli"
