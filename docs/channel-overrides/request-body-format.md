# 渠道参数覆盖：请求体格式

`set_request_body` 用于选择上游请求体编码方式，可通过 `conditions` 为不同模型设置不同格式。模型名称和字段映射均写在渠道配置中，不需要增加 Go 模型分支。

当前 LCON 视频适配器已接入此能力。未匹配此操作时，继续使用原有的 multipart 文件上传流程。未接入的适配器遇到匹配的 `set_request_body` 会返回配置错误，不会忽略它。

## 配置示例

将以下操作加入渠道参数覆盖的 `operations` 数组。`original_model` 对应客户端使用的模型名（模型映射前）。

```json
{
  "mode": "set_request_body",
  "value": {
    "format": "json",
    "field_types": {
      "duration": "integer"
    }
  },
  "conditions": [
    {
      "path": "original_model",
      "mode": "full",
      "value": "你的模型名"
    }
  ]
}
```

| 配置项 | 说明 |
| --- | --- |
| `format: "json"` | 使用 JSON 请求体，自动设置 `application/json`。 |
| `format: "multipart"` | 把覆盖后的 JSON 对象编码为文本表单，自动生成 multipart boundary。 |
| `field_types` | 可选。指定字段类型转换，支持点路径及数组索引，例如 `input.duration`、`images.0`。 |
| `conditions`、`logic` | 沿用参数覆盖的匹配规则。多条 `set_request_body` 命中时，最后一条配置生效，不合并。 |

`field_types` 支持 `integer`、`number`、`boolean`、`string`。例如把 `seconds` 移到 `duration` 后，将字符串 `"5"` 转为 JSON 整数 `5`。整数转换要求整数文本；小数或超出 int64 范围的值会报错。数值转换只接受有限数值。不存在的字段保持缺省，`null` 保持 `null`，显式的 `0`、`false` 不丢失；对象和数组不能作为标量转换。字段路径不支持通配符和查询表达式。

## 字段与文件行为

- 请求格式操作只写入中继元数据，`format`、`field_types` 不会被添加到上游请求体。
- 先完成所有字段覆盖操作，再编码。字段新增用 `set`，改名用 `move`，复制图片用 `copy`，移除旧字段用 `delete`。编码器保留覆盖后的所有字段，不使用模型字段白名单。
- `model` 使用渠道映射后的上游模型名。比例、分辨率、时长没有编码器默认值，需要通过渠道覆盖设置。
- 新编码流程的输入必须为 JSON 对象。客户端直接上传 multipart 文件时继续走原有适配器流程。
- 配置 `multipart` 时，字符串按普通文本写入；数组、对象、数字、布尔值和 `null` 按 JSON 文本写入各个表单项。图片 URL 不会在此流程中自动下载为文件，文件上传需求使用原有适配器流程。
- `Content-Type` 由编码器设置，避免另外使用请求头覆盖改变它，尤其不能手动填写 multipart boundary。
- 使用新请求体格式的请求发生重试时，重新使用客户端原始请求和当前渠道配置，避免沿用其他渠道改名、删除后的字段或请求格式。未启用新格式的请求保留原有字段覆盖流程。

渠道 13 的完整示例见 [channel-13-lcon-seedance-mini.json](channel-13-lcon-seedance-mini.json)，包括原有模型规则、JSON 格式、比例和时长改名、分辨率、复制第二张图及删除旧字段。此文件需要保存到渠道参数覆盖配置后才生效；现有运行程序需要包含此代码改动。

## 适配器接入

公共能力位于 `relay/common/request_body_override.go`。其他任务适配器接入时，需要：

1. 实现 `TaskRequestBodyOverrideAdaptor.SupportsRequestBodyOverride()` 并返回 `true`。
2. 对需要从原始客户端 JSON 生成上游请求的适配器，实现 `TaskPreBuildParamOverrideAdaptor`，在构建前执行覆盖。
3. 在构建请求体时，如果 `info.RequestBodyOverride` 存在，调用 `BuildRequestBodyOverride`，使用返回的请求体和 `Content-Type`；否则保留适配器原有流程。

同一适配器后续增加模型，只需要增加渠道的模型配置及参数覆盖规则。
