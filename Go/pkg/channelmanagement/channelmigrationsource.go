package channnelmanagement

import (
	"fmt"
	"io"
	"net/http"
)

func RunChannelManagemntSource(module, channelId, channelName *string) {

	// Build API
	baseURL := fmt.Sprintf("https://channelmanagement.nxg.revenuepremier.com/api/v1/channel/GetChannel?module=%s&channelId=%s&channelName=%s", *module, *channelId, *channelName)

	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("base URL failed:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
