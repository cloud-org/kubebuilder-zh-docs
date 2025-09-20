# Grafana 插件（`grafana/v1-alpha`）

Grafana 插件是一个可选插件，用于脚手架生成 Grafana Dashboard，帮助你查看由使用 [controller-runtime][controller-runtime] 的项目导出的默认指标。

<aside class="note">
<h1>示例</h1>

你可以在仓库根目录的 [testdata][testdata] 下，查看 `project-v4-with-plugins` 示例中的默认脚手架。

</aside>

## 何时使用？

- 当你希望使用 [Grafana][grafana] 通过 Prometheus 查看 [controller-metrics][controller-metrics] 导出的指标。

## 如何使用？

### 前置条件

- 项目需使用 [controller-runtime][controller-runtime] 暴露 [默认控制器指标][controller-metrics]，并被 Prometheus 采集。
- 可访问 [Prometheus][prometheus]：
  - Prometheus 需暴露可访问的 endpoint（例如 `prometheus-operator` 的 `http://prometheus-k8s.monitoring.svc:9090`）。
  - 该 endpoint 已配置为 Grafana 的数据源。参考 [Add a data source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/)。
- 可访问 [Grafana][grafana-install]，并确保：
  - 拥有仪表盘编辑权限（[Dashboard edit permission][grafana-permissions]）。
  - 已配置 Prometheus 数据源。
    ![pre][prometheus-data-source]

<aside class="note">

参考[指标文档][reference-metrics-doc]了解如何在 Kubebuilder 脚手架的项目中启用指标上报。

在 [config/prometheus][kustomize-plugin] 中可找到用于启用默认 `/metrics` 端点抓取的 ServiceMonitor。

</aside>

### 基本用法

Grafana 插件挂载在 `init` 与 `edit` 子命令上：

```sh
# 使用 grafana 插件初始化新项目
kubebuilder init --plugins grafana.kubebuilder.io/v1-alpha

# 在已有项目上启用 grafana 插件
kubebuilder edit --plugins grafana.kubebuilder.io/v1-alpha
```

插件会创建一个新目录并在其中生成 JSON 文件（例如 `grafana/controller-runtime-metrics.json`）。

#### 使用演示

如下动图展示了在项目中启用该插件：

![output](https://user-images.githubusercontent.com/18136486/175382307-9a6c3b8b-6cc7-4339-b221-2539d0fec042.gif)

#### 如何在 Grafana 中导入这些 Dashboard

1. 复制 JSON 文件内容。
2. 打开 `<your-grafana-url>/dashboard/import`，按指引[导入新的仪表盘](https://grafana.com/docs/grafana/latest/dashboards/export-import/#import-dashboard)。
3. 将 JSON 粘贴到 “Import via panel json”，点击 “Load”。
   <img width="644" src="https://user-images.githubusercontent.com/18136486/176121955-1c4aec9c-0ba4-4271-9767-e8d1726d9d9a.png">
4. 选择作为数据源的 Prometheus。
   <img width="633" src="https://user-images.githubusercontent.com/18136486/176122261-e3eab5b0-9fc4-45fc-a68c-d9ce1cfe96ee.png">
5. 成功导入后，Dashboard 即可使用。

### Dashboard 说明

#### Controller Runtime Reconciliation 总数与错误数

- 指标：
  - `controller_runtime_reconcile_total`
  - `controller_runtime_reconcile_errors_total`
- 查询：
  - `sum(rate(controller_runtime_reconcile_total{job="$job"}[5m])) by (instance, pod)`
  - `sum(rate(controller_runtime_reconcile_errors_total{job="$job"}[5m])) by (instance, pod)`
- 描述：
  - 近 5 分钟内 Reconcile 总次数的每秒速率。
  - 近 5 分钟内 Reconcile 错误次数的每秒速率。
- 示例：<img width="912" src="https://user-images.githubusercontent.com/18136486/176122555-f3493658-6c99-4ad6-a9b7-63d85620d370.png">

#### 控制器 CPU 与内存使用

- 指标：
  - `process_cpu_seconds_total`
  - `process_resident_memory_bytes`
- 查询：
  - `rate(process_cpu_seconds_total{job="$job", namespace="$namespace", pod="$pod"}[5m]) * 100`
  - `process_resident_memory_bytes{job="$job", namespace="$namespace", pod="$pod"}`
- 描述：
  - 近 5 分钟内 CPU 使用率的每秒速率。
  - 控制器进程的常驻内存字节数。
- 示例：<img width="912" src="https://user-images.githubusercontent.com/18136486/177239808-7d94b17d-692c-4166-8875-6d9332e05bcb.png">

#### P50/90/99 工作队列等待时长（秒）

- 指标：
  - `workqueue_queue_duration_seconds_bucket`
- 查询：
  - `histogram_quantile(0.50, sum(rate(workqueue_queue_duration_seconds_bucket{job="$job", namespace="$namespace"}[5m])) by (instance, name, le))`
- 描述：
  - 条目在工作队列中等待被取用的时长。
- 示例：<img width="912" src="https://user-images.githubusercontent.com/18136486/180359126-452b2a0f-a511-4ae3-844f-231d13cd27f8.png">

#### P50/90/99 工作队列处理时长（秒）

- 指标：
  - `workqueue_work_duration_seconds_bucket`
- 查询：
  - `histogram_quantile(0.50, sum(rate(workqueue_work_duration_seconds_bucket{job="$job", namespace="$namespace"}[5m])) by (instance, name, le))`
- 描述：
  - 从工作队列中取出并处理一个条目所花费的时间。
- 示例：<img width="912" src="https://user-images.githubusercontent.com/18136486/180359617-b7a59552-1e40-44f9-999f-4feb2584b2dd.png">

#### Add Rate in Work Queue

- Metrics
  - workqueue_adds_total
- Query:
  - sum(rate(workqueue_adds_total{job="$job", namespace="$namespace"}[5m])) by (instance, name)
- Description
  - Per-second rate of items added to work queue
- Sample: <img width="912" src="https://user-images.githubusercontent.com/18136486/180360073-698b6f77-a2c4-4a95-8313-fd8745ad472f.png">

#### Retries Rate in Work Queue

- Metrics
  - workqueue_retries_total
- Query:
  - sum(rate(workqueue_retries_total{job="$job", namespace="$namespace"}[5m])) by (instance, name)
- Description
  - Per-second rate of retries handled by workqueue
- Sample: <img width="912" src="https://user-images.githubusercontent.com/18136486/180360101-411c81e9-d54e-4b21-bbb0-e3f94fcf48cb.png">

#### Number of Workers in Use

- Metrics
  - controller_runtime_active_workers
- Query:
  - controller_runtime_active_workers{job="$job", namespace="$namespace"}
- Description
  - The number of active controller workers
- Sample: <img width="912" src="https://github.com/kubernetes-sigs/kubebuilder/assets/18136486/288db1b5-e2d8-48ea-9aae-30de7eeca277">

#### WorkQueue Depth

- Metrics
  - workqueue_depth
- Query:
  - workqueue_depth{job="$job", namespace="$namespace"}
- Description
  - Current depth of workqueue
- Sample: <img width="912" src="https://github.com/kubernetes-sigs/kubebuilder/assets/18136486/34f14df4-0428-460e-9658-01dd3d34aade">

#### Unfinished Seconds

- Metrics
  - workqueue_unfinished_work_seconds
- Query:
  - rate(workqueue_unfinished_work_seconds{job="$job", namespace="$namespace"}[5m])
- Description
  - How many seconds of work has done that is in progress and hasn't been observed by work_duration.
- Sample: <img width="912" src="https://github.com/kubernetes-sigs/kubebuilder/assets/18136486/081727c0-9531-4f7a-9649-87723ebc773f">

### Visualize Custom Metrics

The Grafana plugin supports scaffolding manifests for custom metrics.

#### Generate Config Template

When the plugin is triggered for the first time, `grafana/custom-metrics/config.yaml` is generated.

```yaml
---
customMetrics:
#  - metric: # Raw custom metric (required)
#    type:   # Metric type: counter/gauge/histogram (required)
#    expr:   # Prom_ql for the metric (optional)
#    unit:   # Unit of measurement, examples: s,none,bytes,percent,etc. (optional)
```

#### Add Custom Metrics to Config

You can enter multiple custom metrics in the file. For each element, you need to specify the `metric` and its `type`.
The Grafana plugin can automatically generate `expr` for visualization.
Alternatively, you can provide `expr` and the plugin will use the specified one directly.

```yaml
---
customMetrics:
  - metric: memcached_operator_reconcile_total # Raw custom metric (required)
    type: counter # Metric type: counter/gauge/histogram (required)
    unit: none
  - metric: memcached_operator_reconcile_time_seconds_bucket
    type: histogram
```

#### Scaffold Manifest

Once `config.yaml` is configured, you can run `kubebuilder edit --plugins grafana.kubebuilder.io/v1-alpha` again.
This time, the plugin will generate `grafana/custom-metrics/custom-metrics-dashboard.json`, which can be imported to Grafana UI.

#### Show case:

See an example of how to visualize your custom metrics:

![output2][show-case]

## Subcommands

The Grafana plugin implements the following subcommands:

- edit (`$ kubebuilder edit [OPTIONS]`)

- init (`$ kubebuilder init [OPTIONS]`)

## Affected files

The following scaffolds will be created or updated by this plugin:

- `grafana/*.json`

## Further resources

- Check out [video to show how it works][video]
- Checkout the [video to show how the custom metrics feature works][video-custom-metrics]
- Refer to a sample of `serviceMonitor` provided by [kustomize plugin][kustomize-plugin]
- Check the [plugin implementation][plugin-implementation]
- [Grafana Docs][grafana-docs] of importing JSON file
- The usage of serviceMonitor by [Prometheus Operator][servicemonitor]

[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[grafana]: https://grafana.com/docs/grafana/next/
[grafana-docs]: https://grafana.com/docs/grafana/latest/dashboards/export-import/#import-dashboard
[kube-prometheus]: https://github.com/prometheus-operator/kube-prometheus
[prometheus]: https://prometheus.io/docs/introduction/overview/
[prom-operator]: https://prometheus-operator.dev/docs/prologue/introduction/
[servicemonitor]: https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md#related-resources
[grafana-install]: https://grafana.com/docs/grafana/latest/setup-grafana/installation/
[grafana-permissions]: https://grafana.com/docs/grafana/next/administration/roles-and-permissions/#dashboard-permissions
[prometheus-data-source]: https://user-images.githubusercontent.com/18136486/176119794-f6d69b0b-93f0-4f9e-a53c-daf9f77dadae.gif
[video]: https://youtu.be/-w_JjcV8jXc
[video-custom-metrics]: https://youtu.be/x_0FHta2HXc
[show-case]: https://user-images.githubusercontent.com/18136486/186933170-d2e0de71-e079-4d1b-906a-99a549d66ebf.gif
[controller-metrics]: ./../../reference/metrics-reference.md
[kustomize-plugin]: ./../../../../../testdata/project-v4-with-plugins/config/prometheus/monitor.yaml
[plugin-implementation]: ./../../../../../pkg/plugins/optional/grafana/
[reference-metrics-doc]: ./../../reference/metrics.md#exporting-metrics-for-prometheus
[testdata]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/testdata/project-v4-with-plugins
