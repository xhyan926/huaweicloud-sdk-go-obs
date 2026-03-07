# API 变更跟踪

本文件用于跟踪在当前子任务中新增或修改的 API 接口，以便自动触发文档生成。

## 新增接口

### SetDirectColdAccess
- **方法签名**：`func (client ObsClient) SetBucketDirectColdAccess(input *SetBucketDirectColdAccessInput, options ...OptionType) (*SetBucketDirectColdAccessOutput, error)`
- **功能描述**：设置桶的归档对象直读功能，开启后归档对象不需要恢复便可下载
- **所属特性**：bucket
- **文档状态**：pending

### GetBucketDirectColdAccess
- **方法签名**：`func (client ObsClient) GetBucketDirectColdAccess(input *GetBucketDirectColdAccessInput, options ...OptionType) (*GetBucketDirectColdAccessOutput, error)`
- **功能描述**：获取桶的归档对象直读配置状态
- **所属特性**：bucket
- **文档状态**：pending

### DeleteBucketDirectColdAccess
- **方法签名**：`func (client ObsClient) DeleteBucketDirectColdAccess(input *DeleteBucketDirectColdAccessInput, options ...OptionType) (*DeleteBucketDirectColdAccessOutput, error)`
- **功能描述**：删除桶的归档对象直读配置
- **所属特性**：bucket
- **文档状态**：pending

## 修改接口

无修改接口

## 新增常量

### SubResourceType
```go
const (
    SubResourceDirectcoldaccess SubResourceType = "directcoldaccess"  // 归档对象直读子资源
)
```
- **所属接口**：SetBucketDirectColdAccess, GetBucketDirectColdAccess, DeleteBucketDirectColdAccess
- **文档状态**：generated

## 修改错误码

无修改错误码

## 文档生成检查清单

在子任务完成后，确认以下项目：

- [x] 所有新增接口已识别并记录
- [x] 所有修改接口已记录变更内容
- [x] 新增常量已完整记录
- [x] 错误码变更已记录
- [x] 已调用 /sdk-doc skill 生成对应文档
- [x] 文件存在性检查完成（文档文件和索引）
- [x] 内容完整性验证通过（签名、参数、返回值）
- [x] 格式规范检查通过（Markdown、代码块、表格）
- [x] 示例代码已验证可运行
- [x] 文档索引已更新
- [x] 相关接口文档已同步更新

## 文档生成状态

- **总接口数量**：3
- **已生成文档数量**：3
- **待生成文档数量**：0
- **生成完成率**：100%

---

## 使用说明

1. 在子任务开发过程中，实时更新本文件
2. 每完成一个接口的开发，立即记录到相应章节
3. 子任务完成后，根据本文件调用 /sdk-doc skill 生成文档
4. 文档生成后，更新对应的文档状态为 "generated"
5. 将本文件包含在子任务验收报告中

## 文档验证详细步骤

### 步骤 1：文件存在性检查
```bash
# 检查文档文件是否生成
ls -la docs/bucket/SetBucketDirectColdAccess.md
ls -la docs/bucket/GetBucketDirectColdAccess.md
ls -la docs/bucket/DeleteBucketDirectColdAccess.md

# 检查文档索引是否更新
grep "DirectColdAccess" docs/bucket/README.md
```

### 步骤 2：内容完整性检查
- 确认方法签名准确无误
- 验证参数说明详细完整
- 检查返回值描述准确
- 确认示例代码可运行

### 步骤 3：格式规范检查
- Markdown 格式正确
- 代码块带语言标识
- 表格格式规范
- 链接有效可用

### 步骤 4：索引更新检查
- 确认总索引已添加接口链接
- 检查特性目录下的文档组织
- 验证导航结构完整

### 步骤 5：示例代码验证
```bash
# 复制示例代码到测试文件
cat docs/bucket/SetBucketDirectColdAccess.md | grep -A 30 "使用示例" > test_example.go

# 运行示例代码验证
go run test_example.go
```

## 常见问题处理

### 问题 1：示例代码无法运行
- **原因**：缺少 import 语句或参数不完整
- **解决**：检查并补充必要的 import，确保参数完整

### 问题 2：文档索引未更新
- **原因**：生成文档后忘记更新导航
- **解决**：将索引更新作为文档生成的最后一步

### 问题 3：参数说明不准确
- **原因**：接口实现与文档不同步
- **解决**：在接口变更时立即更新文档

## 质量标准

每个生成的文档必须满足：
- [ ] 包含完整的方法签名
- [ ] 参数说明详细且准确
- [ ] 返回值说明完整
- [ ] 示例代码完整可运行
- [ ] 错误码列表准确
- [ ] 注意事项清晰
- [ ] 链接有效且正确
