# proxymanager

A simple GoLang implementation of a proxy manager package. Can be initialized with a txt file filled with proxies.
When initialised it will create a shuffled list of the proxies and return it formatted http://username:password@ip:port aswell as
initiating an internal tracking state. From the ProxyManager object, AssignProxy() can be used when initialising tasks to return the 
first proxies which are shuffled and ready to be used. You can pass in a failing / failed proxy into NextProxy() for it to be released 
from being used and to get a sequential proxy / a proxy not being used. If task count > amount of proxies it will just get a random 
proxy whilst trying to get one not being used and not being the proxy just released. RandomProxy() will just return a random proxy 
formatted http://user:pass@ip:port. GetProxies() will return the full shuffledproxy list formatted http://user:pass@ip:port.

### Installation
```
go get github.com/5tevo/goproxymanager
```

### Usage

Initialize new object
```golang
package main 

import 
"fmt"

"github.com/5tevo/goproxymanager"

func main() {
    // Initialises the proxy manager by reading proxies.txt, shuffling them and returning a new ProxyManager instance. It takes proxies stored in the ip:port:user:pass format
    pm, err := proxymanager.NewManager("proxies.txt")
	if err != nil {
		fmt.Fatalf("Error initialising Proxy Manager: %v", err)
	}
    
    // Gets the initial proxies formatted http://username:password@ip:port and marks them as being used
    proxy, err := pm.AssignProxy()
    if err != nil {
        fmt.Printf("Error obtaining initial proxy: %v", err)
        return
    }

    // Gets next proxy in the manager & passes in the current proxy releasing it formatted http://username:password@ip:port
    proxy, err = pm.NextProxy(proxy)
    if err != nil{
        fmt.Printf("Error obtaining next proxy: %v", err)
    }
   
    // Gets random proxy in the manager formatted http://username:password@ip:port (will not assign its state as being used)
    randomProxy, err := pm.RandomProxy()
    if err != nil {
        fmt.Printf("Error obtaining random proxy: %v", err)
    }

    // Returns the full shuffled proxy list formatted http://username:password@ip:port
    // Useful for if you'd like to handle proxies yourself as soon as the application starts
    proxies := pm.GetProxies()
    if err != nil {
        fmt.Printf("Error obtaining proxy list: %v", err)
    }

}
```
