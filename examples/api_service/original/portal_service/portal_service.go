package portal_service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := getUser(id)
		maskedName := desensitizeUserName(name)
		c.String(http.StatusOK, maskedName)
	})
	_ = r.Run(":50080")
}

// desensitizeUserName 用户名称脱敏
func desensitizeUserName(name string) string {
	runes := []rune(name)

	for i := 1; i < len(runes); i++ {
		runes[i] = '*'
	}

	return string(runes)
}

// getUser 调用其他服务，根据用户ID获取用户名称
func getUser(id string) string {
	client := http.Client{}

	url := fmt.Sprintf("http://127.0.0.1:50081/users/%s", id)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}

	return string(body)
}
