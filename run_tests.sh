#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}NowYouSeeMe - Diary API 测试运行器${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 检查测试数据库是否存在
echo -e "${YELLOW}[1/4] 检查测试数据库...${NC}"
if psql -lqt | cut -d \| -f 1 | grep -qw nowyouseeme_test; then
    echo -e "${GREEN}✓ 测试数据库 nowyouseeme_test 已存在${NC}"
else
    echo -e "${YELLOW}! 测试数据库不存在，正在创建...${NC}"
    createdb nowyouseeme_test
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 成功创建测试数据库${NC}"
    else
        echo -e "${RED}✗ 创建测试数据库失败${NC}"
        exit 1
    fi
fi
echo ""

# 运行迁移
echo -e "${YELLOW}[2/4] 运行数据库迁移...${NC}"
if [ -f "backend/migrations/001_create_event_sourcing_schema.sql" ]; then
    psql -d nowyouseeme_test -f backend/migrations/001_create_event_sourcing_schema.sql -q 2>&1 | grep -v "ERROR.*already exists" || true
    echo -e "${GREEN}✓ 数据库迁移完成${NC}"
else
    echo -e "${RED}✗ 找不到迁移文件${NC}"
    exit 1
fi
echo ""

# 安装依赖
echo -e "${YELLOW}[3/4] 检查Go依赖...${NC}"
cd backend
go mod tidy > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Go依赖已更新${NC}"
else
    echo -e "${RED}✗ 更新Go依赖失败${NC}"
    exit 1
fi
echo ""

# 运行测试
echo -e "${YELLOW}[4/4] 运行测试...${NC}"
echo ""
echo -e "${YELLOW}----------------------------------------${NC}"
go test ./api -v
TEST_RESULT=$?
echo -e "${YELLOW}----------------------------------------${NC}"
echo ""

# 显示结果
if [ $TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo -e "${GREEN}========================================${NC}"
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}✗ 测试失败${NC}"
    echo -e "${RED}========================================${NC}"
    exit 1
fi
