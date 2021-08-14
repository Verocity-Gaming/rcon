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

# Get the current map
```
m, err := c.Map()
if err != nil {
	panic(err)
}

println(m.Location, m.Type, m.Side)
```

# Set the current map
```
err = c.SetMap(rcon.MapCarentanOffensiveUS)
if err != nil {
	panic(err)
}
```

# Get the current map rotation
```
r, err := c.Rotation()
if err != nil {
	panic(err)
}

for _, m := range r {
	println(m.String())
}
```