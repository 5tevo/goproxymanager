package proxymanager

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type ProxyManager struct {
	Proxies      []string
	CurrentIndex int
}

func NewManager(filename string) (*ProxyManager, error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return nil, openErr
	}
	defer file.Close()
	manager := &ProxyManager{Proxies: []string{}, CurrentIndex: 0}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		manager.Proxies = append(manager.Proxies, scanner.Text())
	}
	return manager, scanner.Err()
}

func (p *ProxyManager) LoadProxies(filename string) error {
	p.Proxies = []string{}
	file, openErr := os.Open(filename)
	if openErr != nil {
		return openErr
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p.Proxies = append(p.Proxies, scanner.Text())
	}
	return scanner.Err()
}

func (p *ProxyManager) NextProxy() (string, error) {
	if len(p.Proxies) == 0 {
		return "", errors.New("ProxyManager.Proxies is empty, load proxies")
	}

	p.CurrentIndex++
	if p.CurrentIndex > len(p.Proxies)-1 {
		p.CurrentIndex = 0
	}
	return formatProxy(p.Proxies[p.CurrentIndex]), nil
}

func (p *ProxyManager) RandomProxy() (string, error) {
	if len(p.Proxies) == 0 {
		return "", errors.New("ProxyManager.Proxies is empty, load proxies")
	}

	// Create a new random source and generator
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Select a random proxy
	return formatProxy(p.Proxies[random.Intn(len(p.Proxies))]), nil
}

func formatProxy(proxy string) string {
	parts := strings.Split(proxy, ":")
	if len(parts) == 4 {
		return fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])
	}
	return proxy
}
