package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	"github.com/zarldev/goenums/examples/validation"
)

func main() {
	fmt.Println("=== YAML 序列化和反序列化示例 ===")

	// 创建一个枚举实例
	status := validation.StringStatuses.StringActive
	fmt.Printf("原始枚举值: %s (值: %d)\n", status.String(), status.Val())

	// YAML 序列化
	yamlData, err := yaml.Marshal(status)
	if err != nil {
		log.Fatalf("YAML 序列化失败: %v", err)
	}
	fmt.Printf("YAML 序列化结果: %s", yamlData)

	// YAML 反序列化
	var newStatus validation.StringStatus
	err = yaml.Unmarshal(yamlData, &newStatus)
	if err != nil {
		log.Fatalf("YAML 反序列化失败: %v", err)
	}
	fmt.Printf("反序列化后的枚举值: %s (值: %d)\n", newStatus.String(), newStatus.Val())

	// 验证往返转换
	if status.String() == newStatus.String() && status.Val() == newStatus.Val() {
		fmt.Println("✅ YAML 往返转换成功！")
	} else {
		fmt.Println("❌ YAML 往返转换失败")
	}

	// 测试复杂的 YAML 结构
	fmt.Println("\n=== 复杂 YAML 结构示例 ===")

	data := struct {
		Name   string                    `yaml:"name"`
		Status validation.StringStatus   `yaml:"status"`
		Items  []validation.StringStatus `yaml:"items"`
	}{
		Name:   "示例配置",
		Status: validation.StringStatuses.StringActive,
		Items: []validation.StringStatus{
			validation.StringStatuses.StringActive,
			validation.StringStatuses.StringInactive,
		},
	}

	complexYaml, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("复杂 YAML 序列化失败: %v", err)
	}
	fmt.Printf("复杂 YAML 序列化结果:\n%s", complexYaml)

	// 反序列化复杂结构
	var newData struct {
		Name   string                    `yaml:"name"`
		Status validation.StringStatus   `yaml:"status"`
		Items  []validation.StringStatus `yaml:"items"`
	}

	err = yaml.Unmarshal(complexYaml, &newData)
	if err != nil {
		log.Fatalf("复杂 YAML 反序列化失败: %v", err)
	}

	fmt.Printf("反序列化后的复杂结构:\n")
	fmt.Printf("  名称: %s\n", newData.Name)
	fmt.Printf("  状态: %s\n", newData.Status.String())
	fmt.Printf("  项目列表:\n")
	for i, item := range newData.Items {
		fmt.Printf("    [%d] %s\n", i, item.String())
	}

	fmt.Println("✅ 复杂 YAML 结构处理成功！")
}
