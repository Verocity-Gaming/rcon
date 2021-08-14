# rcon
Go bindings for the HLL RCON tool

# Connect to a server
```
c, err := rcon.New("176.57.165.8:38016", "c78d881")
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