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
	mu           sync.Mutex
	Proxies      []string
	used         []bool
	currentIndex int
	r            *rand.Rand
}

func NewManager(filename string) (*ProxyManager, error) {
	pm := &ProxyManager{
		Proxies:      make([]string, 0),
		currentIndex: 0,
		r:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if err := pm.LoadProxies(filename); err != nil {
		return nil, err
	}
	pm.shuffleProxies()
	pm.used = make([]bool, len(pm.Proxies))
	pm.currentIndex = 0
	return pm, nil
}

func (p *ProxyManager) LoadProxies(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}
	if scannerErr := scanner.Err(); scannerErr != nil {
		return scannerErr
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.Proxies = proxies
	p.currentIndex = 0
	p.used = make([]bool, len(proxies))
	return nil
}

func (p *ProxyManager) shuffleProxies() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.r.Shuffle(len(p.Proxies), func(i, j int) {
		p.Proxies[i], p.Proxies[j] = p.Proxies[j], p.Proxies[i]
	})
}

func (p *ProxyManager) AssignProxy() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := p.currentIndex; i < len(p.Proxies); i++ {
		if !p.used[i] {
			p.used[i] = true
			p.currentIndex = i + 1
			return formatProxy(p.Proxies[i]), nil
		}
	}
	for i := 0; i < p.currentIndex; i++ {
		if !p.used[i] {
			p.used[i] = true
			p.currentIndex = i + 1
			return formatProxy(p.Proxies[i]), nil
		}
	}
	index := p.r.Intn(len(p.Proxies))
	return formatProxy(p.Proxies[index]), nil
}

func (p *ProxyManager) NextProxy(current string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, raw := range p.Proxies {
		if formatProxy(raw) == current {
			p.used[i] = false
			break
		}
	}

	var candidate string
	for i := p.currentIndex; i < len(p.Proxies); i++ {
		formatted := formatProxy(p.Proxies[i])
		if !p.used[i] && formatted != current {
			p.used[i] = true
			candidate = formatted
			p.currentIndex = i + 1
			break
		}
	}
	if candidate == "" {
		for i := 0; i < p.currentIndex; i++ {
			formatted := formatProxy(p.Proxies[i])
			if !p.used[i] && formatted != current {
				p.used[i] = true
				candidate = formatted
				p.currentIndex = i + 1
				break
			}
		}
	}
	if candidate == "" {
		index := p.r.Intn(len(p.Proxies))
		candidate = formatProxy(p.Proxies[index])
	}
	return candidate, nil
}

func (p *ProxyManager) RandomProxy() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Proxies) == 0 {
		return "", errors.New("no proxies loaded")
	}
	index := p.r.Intn(len(p.Proxies))
	return formatProxy(p.Proxies[index]), nil
}

func (p *ProxyManager) GetProxies() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	proxiesCopy := make([]string, len(p.Proxies))
	for i, proxy := range p.Proxies {
		proxiesCopy[i] = formatProxy(proxy)
	}
	return proxiesCopy
}

func formatProxy(proxy string) string {
	parts := strings.Split(proxy, ":")
	if len(parts) == 4 {
		return fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])
	}
	return proxy
}
