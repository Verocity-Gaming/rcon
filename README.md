# rcon
Bindings for the Hell Let Loose RCON tool written purely in Go.

This project is constructed to act like a third-party library for other applications. These might include a REST server, a command line interface, etc.


`import "github.com/verocity-gaming/rcon"`

# Client

## Connect to a server
```
c, err := rcon.New(addr, password)
if err != nil {
	panic(err)
}
defer c.Close()
```

# Maps

## Get the current map
```
m, err := c.Map()
if err != nil {
	panic(err)
}

println(m.Location, m.Type, m.Side)
```

## Set the current map
```
err = c.SetMap(rcon.MapCarentanOffensiveUS)
if err != nil {
	panic(err)
}
```

## Get the current map rotation
```
r, err := c.Rotation()
if err != nil {
	panic(err)
}

for _, m := range r {
	println(m.String())
}
```

## Players
```
p, err := c.Players()
if err != nil {
	panic(err)
}

for _, player := range p {
	println(player.String())
}
```

## VIPs
```
v, err := c.VIPs()
if err != nil {
	panic(err)
}

for _, vip := range v {
	println(vip.String())
}
```