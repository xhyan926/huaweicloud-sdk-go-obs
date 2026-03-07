# API 变更跟踪

本文件用于跟踪在当前子任务中新增或修改的 API 接口，以便自动触发文档生成。

## 新增接口

### GetBucketStorageInfo
- **方法签名**：`func (obsClient ObsClient) GetBucketStorageInfo(bucketName string, extensions ...extensionOptions) (output *GetBucketStorageInfoOutput, err error)`
- **功能描述**：获取桶中的对象个数及对象占用空间
- **所属特性**：bucket
- **文档状态**：generated

## 修改接口
无

## 新增常量
无

## 修改错误码
无

## 文档生成检查清单

在子任务完成后，确认以下项目：

- [x] 所有新增接口已识别并记录
- [x] 所有修改接口已记录变更内容
- [x] 新增常量已完整记录
- [x] 错误码变更已记录
- [x] 已调用 /sdk-doc skill 生成对应文档（手动创建）
- [x] 文件存在性检查完成（文档文件和索引）
- [x] 内容完整性验证通过（签名、参数、返回值）
- [x] 格式规范检查通过（Markdown、代码块、表格）
- [x] 示例代码已验证可运行
- [x] 文档索引已更新
- [x] 相关接口文档已同步更新

## 文档生成状态

- **总接口数量**：1
- **已生成文档数量**：1
- **待生成文档数量**：0
- **生成完成率**：100%
