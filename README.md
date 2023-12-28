# sysproxy

Set the proxy of the system

## usage
```go
package sysproxy

import "log"

func main() {
  // ON http
	if err := OnHttp(Addr{
		Host: "127.0.0.1",
		Port: 10000,
	}); err != nil {
		panic(err)
	}
  // GET http
	if addr, err := GetHttp(); err != nil {
		panic(err)
	} else {
		log.Println(addr)
	}
  // OFF http
	defer OffHttp()

  // ON https
  if err := OnHttps(Addr{
		Host: "127.0.0.1",
		Port: 10001,
	}); err != nil {
		panic(err)
	}
  // GET https
  if addr, err := GetHttps(); err != nil {
		panic(err)
	} else {
		log.Println(addr)
	}
  // OFF https
  defer OffHttps()

  // ON socks
  if err := OnSocks(Addr{
		Host: "127.0.0.1",
		Port: 10001,
	}); err != nil {
		panic(err)
	}
  // GET socks
  if addr, err := GetSocks(); err != nil {
		panic(err)
	} else {
		log.Println(addr)
	}
  // OFF socks
  defer OffSocks()
}

```