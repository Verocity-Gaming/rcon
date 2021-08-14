package rcon

import (
	"fmt"
	"strings"
)

// Map represents a playable landscape in HLL, having parts of the name constructed into a type.
type Map struct {
	Location string
	Type     string
	Side     string

	MapName
}

// MapName represents a server readable name for a map in HLL.
type MapName string

const (
	MapUtahBeachOffensiveUS          MapName = "utahbeach_offensive_us"
	MapFoyWarfare                    MapName = "foy_warfare"
	MapStMarieDuMontOffensiveUS      MapName = "stmariedumont_off_us"
	MapKurskOffensiveRussia          MapName = "kursk_offensive_rus"
	MapStMereEgliseWarfare           MapName = "stmereeglise_warfare"
	MapCarentanOffensiveUS           MapName = "carentan_offensive_us"
	MapHill400Warfare                MapName = "hill400_warfare"
	MapStMarieDuMontOffensiveGermany MapName = "stmariedumont_off_ger"
	MapHurtgenForestOffensiveGermany MapName = "hurtgenforest_offensive_ger"
	MapStalingradWarfare             MapName = "stalingrad_warfare"
	MapFoyOffensiveGermany           MapName = "foy_offensive_ger"
	MapUtahBeachOffensiveGermany     MapName = "utahbeach_offensive_ger"
	MapCarentanWarfare               MapName = "carentan_warfare"
	MapKurskWarfare                  MapName = "kursk_warfare"
	MapPurpleHeartLaneOffensiveUS    MapName = "purpleheartlane_offensive_us"
	MapStMereEgliseOffensiveUS       MapName = "stmereeglise_offensive_us"
	MapUtahBeachWarfare              MapName = "utahbeach_warfare"
	MapStalingradOffensiveGermany    MapName = "stalingrad_offensive_ger"
	MapHurtgenForestWarfare          MapName = "hurtgenforest_warfare_V2"
	MapStMereEgliseOffensiveGermany  MapName = "stmereeglise_offensive_ger"
	MapHill400OffensiveUS            MapName = "hill400_offensive_US"
	MapOmahaBeachOffensiveUS         MapName = "omahabeach_offensive_us"
	MapPurpleHeartLaneWarfare        MapName = "purpleheartlane_warfare"
	MapKurskOffensiveGermany         MapName = "kursk_offensive_ger"
	MapStalingradOffensiveRussia     MapName = "stalingrad_offensive_rus"
	MapStMarieDuMontWarfare          MapName = "stmariedumont_warfare"
	MapHurtgenForestOffensiveUS      MapName = "hurtgenforest_offensive_US"
)

// Map returns the current map in rotation for a Conn.
func (c *Conn) Map() (Map, error) {
	err := c.send("get map")
	if err != nil {
		return Map{}, err
	}

	name, err := c.read()
	if err != nil {
		return Map{}, err
	}

	m := MapName(name)

	return fromName(m), nil
}

// Rotation returns the current map rotation for a Conn.
func (c *Conn) Rotation() ([]Map, error) {
	err := c.send("rotlist")
	if err != nil {
		return nil, err
	}

	result, err := c.read()
	if err != nil {
		return nil, err
	}

	if result == "FAIL" {
		return nil, fmt.Errorf("failed to get map rotation")
	}

	maps := []Map{}
	for _, name := range strings.Split(result, "\n") {
		if name == "" {
			continue
		}

		m := MapName(name)

		maps = append(maps, fromName(m))
	}

	return maps, nil
}

// SetMap will change the current map in rotation for a Conn.
func (c *Conn) SetMap(n MapName) error {
	err := c.send("map " + n.String())
	if err != nil {
		return err
	}

	result, err := c.read()
	if err != nil {
		return err
	}

	if result == "FAIL" {
		return fmt.Errorf("failed to set map for %s", n)
	}

	return nil
}

// String returns a prettier standard string for a Map.
func (m Map) String() string {
	s := fmt.Sprintf("%s - %s", m.Location, m.Type)
	if m.Side != "" {
		s = fmt.Sprintf("%s (%s)", s, m.Side)
	}
	return s
}

// String converts a MapName into a standard string.
func (n MapName) String() string {
	return string(n)
}

func fromName(n MapName) Map {
	name := n.String()

	args := strings.Split(name, "_")
	if len(args) == 0 {
		return Map{}
	}

	m := Map{
		Location: location(args[0]),
		MapName:  MapName(name),
	}

	switch len(args) {
	case 2:
		m.Type = "Warfare"
	case 3:
		m.Type = "Offensive"
		m.Side = side(args[2])
	}

	return m
}

func location(name string) string {
	switch name {
	case "utahbeach":
		return "Utah Beach"
	case "omahabeach":
		return "Omaha Beach"
	case "foy":
		return "Foy"
	case "stmariedumont":
		return "St. Marie Du Mont"
	case "kursk":
		return "Kursk"
	case "stmereeglise":
		return "St. Mere Eglise"
	case "carentan":
		return "Carentan"
	case "hill400":
		return "Hill 400"
	case "hurtgenforest":
		return "Hurtgen Forest"
	case "stalingrad":
		return "Stalingrad"
	case "purpleheartlane":
		return "Purple Heart Lane"
	default:
		return strings.Title(name)
	}
}

func side(name string) string {
	switch name {
	case "us":
		return "United States"
	case "ger":
		return "Germany"
	case "rus":
		return "Russia"
	default:
		return ""
	}
}
