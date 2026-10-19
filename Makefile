# figmafe demo — 常用命令
.PHONY: build test fixture scan map catalog fonts

build:            ## 编译 CLI 到 bin/figmafe
	cd cli && go build -o ../bin/figmafe ./cmd/figmafe

test:             ## CLI 单元测试
	cd cli && go vet ./... && go test ./...

fonts:            ## 安装 Noto Sans SC（截图与设计图对齐需要）
	scripts/fetch-fonts.sh

fixture: build    ## 从设计源生成 Figma fixture（仅在没有真实 Figma 时使用）
	bin/figmafe fixture design/source/order-detail.html -o design/figma-export/order-detail

scan: build       ## Phase 0：盘点两端组件资产
	bin/figmafe scan

map: build        ## Phase 1：映射建议 + 检查
	bin/figmafe map suggest && bin/figmafe map check

catalog: build    ## 生成 Skill 的组件手册
	bin/figmafe catalog
