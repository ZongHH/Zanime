package random

import (
	"fmt"
	"math/rand/v2"
)

// 浏览器基础信息
var (
	browsers = []string{
		"Chrome", "Firefox", "Safari", "Edge", "Opera",
	}

	browserVersions = map[string][]string{
		"Chrome":  generateVersions(70, 120),  // Chrome 70-120
		"Firefox": generateVersions(60, 120),  // Firefox 60-120
		"Safari":  generateVersions(605, 620), // Safari 605-620
		"Edge":    generateVersions(80, 120),  // Edge 80-120
		"Opera":   generateVersions(60, 100),  // Opera 60-100
	}

	// 操作系统
	windows = []string{
		"Windows NT 10.0", "Windows NT 6.3", "Windows NT 6.2", "Windows NT 6.1",
	}

	mac = []string{
		"Macintosh; Intel Mac OS X 10_%d_%d",
		"Macintosh; Intel Mac OS X 10_%d",
	}

	linux = []string{
		"X11; Linux x86_64",
		"X11; Ubuntu; Linux x86_64",
		"X11; Fedora; Linux x86_64",
	}

	// 移动设备
	mobileDevices = []string{
		"iPhone; CPU iPhone OS %d_%d like Mac OS X",
		"iPad; CPU OS %d_%d like Mac OS X",
		"Linux; Android %d.%d",
	}

	// 处理器架构
	cpuTypes = []string{
		"x86_64", "Win64; x64", "amd64", "arm64",
	}
)

// generateVersions 生成版本号范围
func generateVersions(start, end int) []string {
	versions := make([]string, end-start+1)
	for i := range versions {
		versions[i] = fmt.Sprintf("%d.0.%d.%d",
			start+i,
			rand.IntN(5000)+1000,
			rand.IntN(200)+1,
		)
	}
	return versions
}

// 生成随机版本号
func randomVersion() string {
	return fmt.Sprintf("%d.%d",
		rand.IntN(15)+1,
		rand.IntN(10),
	)
}

// GetRandomUserAgent 返回一个随机生成的 UserAgent
func GetRandomUserAgent() string {
	// 随机选择是否生成移动设备UA
	if rand.IntN(100) < 30 { // 30%概率生成移动设备UA
		return generateMobileUA()
	}
	return generateDesktopUA()
}

// generateDesktopUA 生成桌面端 UserAgent
func generateDesktopUA() string {
	browser := browsers[rand.IntN(len(browsers))]
	version := browserVersions[browser][rand.IntN(len(browserVersions[browser]))]

	var platform string
	switch rand.IntN(3) {
	case 0:
		platform = windows[rand.IntN(len(windows))] + "; " + cpuTypes[rand.IntN(len(cpuTypes))]
	case 1:
		platform = fmt.Sprintf(mac[rand.IntN(len(mac))], rand.IntN(5)+10, rand.IntN(10))
	case 2:
		platform = linux[rand.IntN(len(linux))]
	}

	switch browser {
	case "Chrome", "Edge", "Opera":
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36%s",
			platform,
			version,
			getBrowserSuffix(browser, version),
		)
	case "Firefox":
		return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/20100101 Firefox/%s",
			platform,
			version,
			version,
		)
	case "Safari":
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/%s (KHTML, like Gecko) Version/%s Safari/%s",
			platform,
			version,
			randomVersion(),
			version,
		)
	default:
		return ""
	}
}

// generateMobileUA 生成移动端 UserAgent
func generateMobileUA() string {
	device := mobileDevices[rand.IntN(len(mobileDevices))]
	version := fmt.Sprintf(device, rand.IntN(7)+10, rand.IntN(10))

	browserVersion := browserVersions["Chrome"][rand.IntN(len(browserVersions["Chrome"]))]

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Mobile/15E148 Safari/604.1",
		version,
		browserVersion,
	)
}

// getBrowserSuffix 获取浏览器特定后缀
func getBrowserSuffix(browser, version string) string {
	switch browser {
	case "Edge":
		return " Edg/" + version
	case "Opera":
		return " OPR/" + version
	default:
		return ""
	}
}

// GetUserAgentChannel 返回一个能持续产生随机 UserAgent 的通道
func GetUserAgentChannel(count int) <-chan string {
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		for i := 0; i < count; i++ {
			ch <- GetRandomUserAgent()
		}
	}()
	return ch
}
