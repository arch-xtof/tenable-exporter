package client

import (
	"fmt"
	"sort"
	"strings"
)

func cpeGuess(output string, cpes []string) string {
	if cpes == nil {
		return "no details"
	}
	if osCheck(cpes) {
		return getOs(cpes)
	}
	if len(cpes) == 1 {
		return singleCpe(cpes[0])
	}

	lower := strings.ToLower(output)
	if packageCheck(cpes) {
		return getPackage(lower, cpes)
	}

	return getApp(lower, cpes)
}

func singleCpe(cpe string) string {
	sp := strings.Split(cpe, ":")
	if sp[2] == sp[3] {
		return sp[2]
	}
	if sp[0] == "p-cpe" {
		return sp[len(sp)-1]
	}
	return fmt.Sprintf("%s %s", sp[2], sp[3])
}

func osCheck(cpes []string) bool {
	for _, cpe := range cpes {
		sp := strings.Split(cpe, ":")
		if sp[1] != "/o" {
			return false
		}
	}
	return true
}

func getOs(cpes []string) string {
	for _, cpe := range cpes {
		sp := strings.Split(cpe, ":")
		if strings.Contains(sp[3], "mac") {
			return "macos"
		}
		if strings.Contains(sp[3], "ubuntu") {
			return "ubuntu"
		}
		if strings.Contains(sp[3], "windows") {
			return "windows"
		}
	}
	return fmt.Sprintf("unknown os: %s", cpes[0])
}

func packageCheck(cpes []string) bool {
	for _, cpe := range cpes {
		if strings.Split(cpe, ":")[0] == "p-cpe" {
			return true
		}
	}
	return false
}

func getPackage(output string, cpes []string) string {
	if len(cpes) == 1 {
		return strings.Split(cpes[0], ":")[len(cpes[0])-1]
	}
	var pkgs []string
	for _, cpe := range cpes {
		if strings.Split(cpe, ":")[1] == "/o" {
			continue
		}
		if len(pkgs) > 10 {
			return "multiple packages"
		}
		sp := strings.Split(cpe, ":")
		pkgs = append(pkgs, sp[len(sp)-1])
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return len(cpes[i]) > len(cpes[j])
	})
	for i, pkg := range pkgs {
		if i > 1 && pkgs[i] == pkgs[i-1] {
			return pkg
		}
		if strings.Contains(output, pkg) || strings.Contains(output, strings.ReplaceAll(pkg, "_", " ")) {
			return pkg
		}
	}
	return fmt.Sprintf("undefined: %s", cpes[0])
}

func getApp(output string, cpes []string) string {
	sort.Slice(cpes, func(i, j int) bool {
		return len(strings.Split(cpes[i], ":")[3]) > len(strings.Split(cpes[j], ":")[3])
	})
	for _, cpe := range cpes {
		sp := strings.Split(cpe, ":")
		if sp[1] == "/o" {
			continue
		}
		if strings.Contains(output, sp[3]) || strings.Contains(output, strings.ReplaceAll(sp[3], "_", " ")) {
			if sp[2] == sp[3] {
				return sp[2]
			}
			return fmt.Sprintf("%s %s", sp[2], sp[3])
		}
	}

	return fmt.Sprintf("unidentifed: %s", cpes[0])
}
