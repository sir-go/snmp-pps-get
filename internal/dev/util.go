package dev

import (
	"strconv"
	"strings"
)

// "192.168.28.26:5-15" -> "192.168.28.26", "5-15"
// "192.168.28.26" -> "192.168.28.26", ""
func parseTarget(target string) (ip, ports string) {
	addr := strings.SplitN(target, ":", 3)
	if len(addr) < 2 {
		return target, ""
	}
	return addr[0], addr[1]
}

// PortIsInRange (13, "4,6,12-15") -> true
func PortIsInRange(pNum int, pRange string) bool {

	if pNum < 1 || len(pRange) < 1 {
		return false
	}

	// "4,6,12-15" -> "4", "6", "12-15"
	for _, part := range strings.Split(pRange, ",") {
		if !strings.ContainsRune(part, '-') {
			// "4", "6"
			if digit, err := strconv.Atoi(part); err == nil {
				if digit == pNum {
					return true
				}
			}
		} else {
			//"12-15" | "15-12" | "15-" | "-12"
			edges := strings.SplitN(part, "-", 3)

			e0, _ := strconv.Atoi(edges[0])
			e1, _ := strconv.Atoi(edges[1])

			if e0 == e1 {
				if e0 == 0 {
					// "_-_" ? -> next part
					continue
				}
				if e0 == pNum {
					// "15-15"; pNum == 15 ? else -> next part
					return true
				} else {
					continue
				}
			}

			// ........p..........e1
			if e0 == 0 {
				if pNum <= e1 {
					return true
				} else {
					continue
				}
			}

			// e0.......p..........
			if e1 == 0 {
				if pNum >= e0 {
					return true
				} else {
					continue
				}
			}

			// swap values if reversed "15-12" -> "12-15"
			if e0 > e1 {
				e0, e1 = e1, e0
			}

			// e0......p.........e1
			if pNum <= e1 && pNum >= e0 {
				return true
			}
		}
	}

	return false
}
