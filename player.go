package rcon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Name string
	ID64 string
}

type Ban struct {
	Player
	Admin
	Reason string
	time.Duration
	time.Time
}

func (c *Conn) BannedTemporarily() ([]Ban, error) {
	result, err := c.send("get", "tempbans")
	if err != nil {
		return nil, fmt.Errorf("failed to get permanent bans: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse temporary ban information")
	}

	admins, err := c.Admins()
	if err != nil {
		return nil, fmt.Errorf("failed to parse temporary ban admins: %v", err)
	}

	bans := []Ban{}

	for _, ban := range args[1 : len(args)-1] {
		b, err := parseBanned(ban, admins)
		if err != nil {
			return nil, err
		}

		bans = append(bans, b)
	}

	return bans, nil
}

func (b Ban) String() string {
	return fmt.Sprintf("%s [%s] from: %s until %s by: %s", b.Player.String(), b.Reason, b.Time.Format(time.Stamp), time.Now().Add(b.Duration).Format(time.Stamp), b.Admin)
}

func (c *Conn) BannedPermanently() ([]Ban, error) {
	result, err := c.send("get", "permabans")
	if err != nil {
		return nil, fmt.Errorf("failed to get permanent bans: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse temporary ban information")
	}

	admins, err := c.Admins()
	if err != nil {
		return nil, fmt.Errorf("failed to parse temporary ban admins: %v", err)
	}

	bans := []Ban{}

	for _, ban := range args[1 : len(args)-1] {
		b, err := parseBanned(ban, admins)
		if err != nil {
			return nil, err
		}

		bans = append(bans, b)
	}

	return bans, nil
}

// BanPermanently will remove an active player and block server access indefinitely.
func (c *Conn) BanPermanently(p Player, reason, admin string) error {
	_, err := c.send("permaban", q(p.ID64), q(reason), q(admin))
	if err != nil {
		return fmt.Errorf("failed to set permanent ban %s: %v", p, err)
	}

	return nil
}

// BanRemove will remove a Player's temp or perma ban and re-allow server access.
func (c *Conn) BanRemove(p Player) error {
	_, err := c.send("pardontempban", q(p.ID64))
	if err != nil {
		if err != ErrResultFailed {
			return fmt.Errorf("failed to remove ban for %s: %v", p, err)
		}

		_, err := c.send("pardonpermaban", q(p.ID64))
		if err != nil {
			return fmt.Errorf("failed to remove ban for %s: %v", p, err)
		}
	}

	return nil
}

// BanTemporarily will remove an active player and block server access for the specified hours.
func (c *Conn) BanTemporarily(p Player, hours int, reason, admin string) error {
	_, err := c.send("tempban", q(p.ID64), strconv.Itoa(hours), q(reason), q(admin))
	if err != nil {
		return fmt.Errorf("failed to set temporary ban %s: %v", p, err)
	}

	return nil
}

// Kick will remove an active player.
func (c *Conn) Kick(p Player, reason string) error {
	_, err := c.send("kick", q(p.Name), q(reason))
	if err != nil {
		return fmt.Errorf("failed to kick %s: %v", p, err)
	}

	return nil
}

// Punish will punish an active player.
func (c *Conn) Punish(p Player, reason string) error {
	_, err := c.send("punish", q(p.Name), q(reason))
	if err != nil {
		return fmt.Errorf("failed to punish %s: %v", p, err)
	}

	return nil
}

// Player return a Player for a given username.
func (c *Conn) Player(username string) (Player, error) {
	result, err := c.send("playerinfo", username)
	if err != nil {
		return Player{}, fmt.Errorf("failed to get player information for %s: %v", username, err)
	}

	args := strings.Split(result, "\n")
	if len(args) < 2 {
		return Player{}, fmt.Errorf("invalid player information for %s", username)
	}

	name := strings.Split(args[0], ":")
	id := strings.Split(args[1], ":")

	if len(name) < 2 || len(id) < 2 {
		return Player{}, fmt.Errorf("invalid player information for %s", username)
	}

	return Player{
		Name: strings.Trim(name[1], " "),
		ID64: strings.Trim(id[1], " "),
	}, nil
}

// Players returns all active Players.
func (c *Conn) Players() ([]Player, error) {
	result, err := c.send("get", "playerids")
	if err != nil {
		return nil, fmt.Errorf("failed to get information for all players: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for all players")
	}

	players := []Player{}

	for _, p := range args[1 : len(args)-1] {
		nameID := strings.Split(p, ":")

		players = append(players, Player{
			Name: strings.Trim(nameID[0], " "),
			ID64: strings.Trim(nameID[1], " "),
		})
	}

	return players, nil
}

func (c *Conn) SetSwitchTeamNow(p Player) error {
	_, err := c.send("switchteamnow", p.ID64)
	if err != nil {
		return fmt.Errorf("failed to set switch player now for %s: %v", p.String(), err)
	}

	return nil
}

func (c *Conn) SetSwitchTeamOnDeath(p Player) error {
	_, err := c.send("switchteamondeath", p.ID64)
	if err != nil {
		return fmt.Errorf("failed to set switch player on death for %s: %v", p.String(), err)
	}

	return nil
}

func (p Player) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.ID64)
}

var matchBannedName = regexp.MustCompile(`: nickname "(.*?)" banned"`)
var matchBanned = regexp.MustCompile(`(.*?) : nickname "(.*?)" banned for (.*?) hours on (.*?) for "(.*?)" by admin "(.*?)"`)

// parseBanned will return a slice of items that are in quotes from a string.
func parseBanned(s string, admins []Admin) (Ban, error) {
	b := Ban{}

	matches := matchBanned.FindAllStringSubmatch(s, -1)
	if len(matches) != 1 {
		return b, fmt.Errorf("failed to parse temporary ban information")
	}

	match := matches[0]

	b.ID64 = match[1]
	b.Name = match[2]

	hours, err := strconv.Atoi(match[3])
	if err != nil {
		return Ban{}, fmt.Errorf("failed to parse temporary ban hours: %v", err)
	}

	b.Duration = time.Duration(hours) * time.Hour

	b.Time, err = time.Parse("2006.01.02-15.04.05", match[4])
	if err != nil {
		return Ban{}, fmt.Errorf("failed to parse temporary ban time: %v", err)
	}

	b.Reason = match[5]

	b.Admin = unknownAdmin

	for i := range admins {
		if admins[i].Name == match[6] || admins[i].ID64 == match[6] {
			b.Admin = admins[i]
			break
		}
	}

	return b, nil
}
