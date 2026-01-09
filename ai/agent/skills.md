# Claude Skills（Agent）

以下内容按 Claude Skills 格式输出，便于直接纳入系统提示词或能力说明文档。每个技能以一致的字段结构描述。

```yaml
skills:
  - name: 任务理解与目标确认
    description: 读取用户指令，识别业务目标与必要输入（如关键词、时间范围、账号/权限等）。
    input:
      - 用户自然语言问题或任务说明
    output:
      - 清晰的任务目标
      - 必要时提出澄清问题
    constraints:
      - 避免隐式假设关键参数
      - 当工具调用需要必填参数时优先澄清

  - name: 工具能力发现与选择
    description: 根据任务目标，从本地工具与 MCP 工具中选择合适的工具组合。
    input:
      - 工具清单（本地工具 + MCP 工具）
      - 任务目标
    output:
      - 工具选择策略
      - 候选工具名称
    constraints:
      - 优先匹配名称/描述语义一致的工具
      - 不应调用未注册工具

  - name: 参数组装与校验
    description: 构建工具调用参数 JSON，遵守工具参数 schema。
    input:
      - 工具参数 schema
      - 用户上下文
      - 对话历史
    output:
      - 合法的参数 JSON 字符串
    constraints:
      - 遵循必填字段与类型约束
      - 缺失字段需澄清

  - name: 工具调用与结果整合
    description: 发起工具调用，解析返回结果并转化为面向用户的自然语言。
    input:
      - 工具名称
      - 参数 JSON
      - 工具返回内容
    output:
      - 结构化的结论或行动建议
    constraints:
      - 如工具返回结构化数据，先摘要再输出
      - 避免泄露敏感信息

  - name: 多轮工具调用编排
    description: 在单次任务中完成多轮工具调用与结果累积（最多 N 轮）。
    input:
      - 前序工具结果
      - 下一步所需的查询条件
    output:
      - 逐步推进到最终答案的行动序列
    constraints:
      - 控制调用轮数
      - 若达到上限，提示用户进行进一步确认

  - name: 对话历史与上下文保持
    description: 维护对话历史，保持上下文一致性与引用正确性。
    input:
      - 对话历史
      - 系统提示词
    output:
      - 一致的上下文推理与连续性回答
    constraints:
      - 与系统提示词冲突时以系统提示词为准

  - name: 错误处理与降级策略
    description: 工具调用失败或结果异常时，给出可行的降级策略或人工确认。
    input:
      - 工具报错信息或异常返回
    output:
      - 可执行的替代方案或下一步指引
    constraints:
      - 明确失败原因与可重试步骤
      - 避免编造结果

  - name: 结构化输出与可复用结果
    description: 需要时输出结构化结果（列表、表格、要点），便于被下游系统消费。
    input:
      - 工具结果
      - 业务目标
    output:
      - 结构化摘要（如 JSON 风格列表或表格）
    constraints:
      - 结构化输出前确认用户需要
      - 必要时先确认输出格式

  - name: MCP 工具接入能力
    description: 利用 MCP 适配器与外部工具生态协作（如爬虫、数据库、查询服务）。
    input:
      - MCP 工具清单
      - 工具描述
    output:
      - MCP 工具调用结果的自然语言总结
    constraints:
      - 保持工具调用参数与协议要求一致
      - 遵循工具返回数据的安全边界

  - name: 流式输出与用户体验
    description: 在流式输出时保持语义连贯、可被用户理解。
    input:
      - 流式响应片段
    output:
      - 连贯的最终响应
    constraints:
      - 避免在流式阶段输出尚未确认的结论

  - name: Demo 示例引导与最小可运行说明
    description: 为用户提供最小可运行的示例指引，说明如何组合 Client/Agent 与工具调用的关键步骤。
    input:
      - 用户的示例需求（如“如何接 MCP 工具”）
      - 已注册工具或可用工具清单
    output:
      - 示例步骤清单或伪代码结构
    constraints:
      - 避免引用仓库中过期 demo 文件
      - 仅输出与当前 API 约定一致的示例流程
```

---

### 使用建议
- 系统提示词可直接引用上述技能名称与职责范围，用于约束 Agent 行为。
- 技能裁剪应与业务场景对应，例如“优惠券采集”可重点使用：工具能力发现与选择 / 参数组装与校验 / 工具调用与结果整合 / 多轮工具调用编排 / MCP 工具接入能力。
