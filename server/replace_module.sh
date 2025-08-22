#!/bin/bash

# 替换项目模块名的脚本
# 用法: ./replace_module.sh "github.com/yourname/your-project"

if [ $# -eq 0 ]; then
    echo "用法: $0 新的模块名"
    echo "例如: $0 github.com/yourname/your-project"
    exit 1
fi

OLD_MODULE="github.com/flipped-aurora/gin-vue-admin/server"
NEW_MODULE="$1"

echo "开始替换模块名..."
echo "旧模块名: $OLD_MODULE"
echo "新模块名: $NEW_MODULE"

# 1. 替换 go.mod 文件
echo "正在替换 go.mod..."
sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" go.mod

# 2. 替换所有 .go 文件中的导入路径
echo "正在替换所有 .go 文件中的导入路径..."
find . -name "*.go" -type f -exec sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" {} +

# 3. 替换其他可能的配置文件
echo "正在检查其他配置文件..."
find . -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.md" | xargs grep -l "$OLD_MODULE" 2>/dev/null | while read file; do
    echo "替换文件: $file"
    sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" "$file"
done

# 4. 清理备份文件
echo "正在清理备份文件..."
find . -name "*.bak" -type f -delete

# 5. 更新模块依赖
echo "正在更新 Go 模块依赖..."
go mod tidy

echo "替换完成！"
echo ""

