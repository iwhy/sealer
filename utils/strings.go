// Copyright © 2021 Alibaba Group Holding Ltd.

package utils

import (
	"bytes"
	"net"
	"sort"
	"strings"
)

func NotIn(key string, slice []string) bool {
	for _, s := range slice {
		if key == s {
			return false
		}
	}
	return true
}

func NotInIPList(key string, slice []string) bool {
	for _, s := range slice {
		if s == "" {
			continue
		}
		if key == strings.Split(s, ":")[0] {
			return false
		}
	}
	return true
}

func ReduceIPList(src, dst []string) []string {
	var ipList []string
	for _, ip := range src {
		if !NotIn(ip, dst) {
			ipList = append(ipList, ip)
		}
	}
	return ipList
}

func AppendIPList(src, dst []string) []string {
	for _, ip := range dst {
		if NotIn(ip, src) {
			src = append(src, ip)
		}
	}
	return src
}

func SortIPList(iplist []string) {
	realIPs := make([]net.IP, 0, len(iplist))
	for _, ip := range iplist {
		realIPs = append(realIPs, net.ParseIP(ip))
	}

	sort.Slice(realIPs, func(i, j int) bool {
		return bytes.Compare(realIPs[i], realIPs[j]) < 0
	})

	for i := range realIPs {
		iplist[i] = realIPs[i].String()
	}
}
