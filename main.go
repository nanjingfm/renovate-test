package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// 创建一个带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 示例 1：获取网页标题
	fmt.Println("示例 1：获取网页标题")
	title, err := GetPageTitle(ctx, "https://golang.org")
	if err != nil {
		log.Printf("获取网页标题失败: %v", err)
	} else {
		fmt.Printf("网页标题: %s\n", title)
	}

	// 示例 2：使用 HttpClient 获取网页
	fmt.Println("\n示例 2：使用 HttpClient 获取网页")
	client := NewHttpClient(30 * time.Second)
	resp, err := client.Get(ctx, "https://httpbin.org/json")
	if err != nil {
		log.Printf("HTTP 请求失败: %v", err)
	} else {
		fmt.Printf("HTTP 状态码: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// 示例 3：解析 HTML 内容
	fmt.Println("\n示例 3：解析 HTML 内容")
	htmlContent := `<html><head><title>测试页面</title></head><body><h1>Hello World</h1></body></html>`
	parser := NewHtmlParser()
	title, err = parser.ExtractTitle(htmlContent)
	if err != nil {
		log.Printf("解析 HTML 失败: %v", err)
	} else {
		fmt.Printf("解析出的标题: %s\n", title)
	}

	fmt.Println("\n所有示例执行完成！")
}
