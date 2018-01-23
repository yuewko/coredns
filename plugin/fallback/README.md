# fallback

This plugin will send queries to an alternate set of upstreams if the plugin
chain returns specific error messages.

```
{
	log
	fallback NXDOMAIN . 192.168.1.1:53
	proxy . 8.8.8.8
}
```
