package dnsmasq

var dnsmasqConfTemplate = `# managed by interlock
port={{ .Port }}
`
