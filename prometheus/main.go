package main

import (
	"fmt"
	"sort"
)

// MetricData 模拟从Prometheus中获取的指标数据
type MetricData struct {
	Cluster      string
	Namespace    string
	Pod          string
	Container    string
	CPUUsageRate float64
	PodInfoMax   float64
	Node         string
}

func main() {
	// 模拟从Prometheus中获取的指标数据
	metricsData := []MetricData{
		{"cluster1", "namespace1", "pod1", "container1", 10.0, 42.0, "node1"},
		{"cluster1", "namespace1", "pod1", "container2", 20.0, 40.0, "node1"},
		{"cluster1", "namespace1", "pod1", "container3", 30.0, 38.0, "node1"},
		// 添加更多指标数据
	}

	// 计算CPU使用率的变化率
	cpuUsageRates := make(map[string]float64)
	for _, metric := range metricsData {
		key := fmt.Sprintf("%s/%s/%s/%s", metric.Cluster, metric.Namespace, metric.Pod, metric.Container)
		cpuUsageRates[key] = metric.CPUUsageRate
	}

	// 计算kube_pod_info指标的最大值
	podInfoMaxValues := make(map[string]float64)
	for _, metric := range metricsData {
		key := fmt.Sprintf("%s/%s/%s/%s", metric.Cluster, metric.Namespace, metric.Pod, metric.Node)
		if value, ok := podInfoMaxValues[key]; !ok || metric.PodInfoMax > value {
			podInfoMaxValues[key] = metric.PodInfoMax
		}
	}

	// 根据最大值对CPU使用率变化率进行排序
	sortedKeys := sortKeysByValueDesc(cpuUsageRates)
	topKey := sortedKeys[0]

	// 输出结果
	fmt.Printf("Top CPU Usage Rate Key: %s, Value: %f\n", topKey, cpuUsageRates[topKey])
	fmt.Printf("Corresponding PodInfoMax Value: %f\n", podInfoMaxValues[topKey])
}

// sortKeysByValueDesc 根据map的值降序排序keys
func sortKeysByValueDesc(m map[string]float64) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return keys
}
