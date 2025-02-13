package proxymanager

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type ProxyManager struct {
	pool    chan string
	proxies []string
	r       *rand.Rand
	mu      sync.Mutex
}

func NewManager(filename string) (*ProxyManager, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var raw []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			raw = append(raw, proxy)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(raw), func(i, j int) {
		raw[i], raw[j] = raw[j], raw[i]
	})

	formatted := make([]string, len(raw))
	for i, proxy := range raw {
		formatted[i] = formatProxy(proxy)
	}

	pool := make(chan string, len(formatted))
	for _, proxy := range formatted {
		pool <- proxy
	}

	return &ProxyManager{
		pool:    pool,
		proxies: formatted,
		r:       r,
	}, nil
}

func (pm *ProxyManager) AssignProxy() (string, error) {
	select {
	case proxy := <-pm.pool:
		return proxy, nil
	default:
		return pm.RandomProxy()
	}
}

func (pm *ProxyManager) NextProxy(current string) (string, error) {
	select {
	case pm.pool <- current:
	default:
	}

	var newProxy string
	select {
	case newProxy = <-pm.pool:
	default:
		var err error
		newProxy, err = pm.RandomProxy()
		if err != nil {
			return "", err
		}
	}

	if newProxy == current && len(pm.proxies) > 1 {
		var err error
		newProxy, err = pm.RandomProxy()
		if err != nil {
			return "", err
		}
	}
	return newProxy, nil
}

func (pm *ProxyManager) RandomProxy() (string, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if len(pm.proxies) == 0 {
		return "", errors.New("no proxies available")
	}
	return pm.proxies[pm.r.Intn(len(pm.proxies))], nil
}

func (pm *ProxyManager) GetProxies() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	copied := make([]string, len(pm.proxies))
	copy(copied, pm.proxies)
	return copied
}

func formatProxy(proxy string) string {
	parts := strings.Split(proxy, ":")
	if len(parts) == 4 {
		return fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])
	}
	return proxy
}
