# rcon
Bindings for the Hell Let Loose RCON tool written purely in Go.

# Connect to a server
```
c, err := rcon.New(addr, password)
if err != nil {
	panic(err)
}
defer c.Close()
```

# Get the current Map
```
m, err := c.Map()
if err != nil {
	panic(err)
}

println(m.Location, m.Type, m.Side)
```

# Set the current Map
```
err = c.SetMap(rcon.MapCarentanOffensiveUS)
if err != nil {
	panic(err)
}
```