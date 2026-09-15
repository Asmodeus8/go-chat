package server

import (
	"fmt"
	"strings"
)

func handleCommand(h *Hub, c *Client, line string) bool {
	parts := strings.Fields(line)
	cmd := strings.ToLower(parts[0])
	switch cmd {
	case "/help":
		c.Send("commands: /users /rooms /join <room> /msg <user> <text> /quit")
	case "/users":
		c.Send("users: " + strings.Join(h.Users(c.Room), ", "))
	case "/rooms":
		c.Send("rooms: " + strings.Join(h.Rooms(), ", "))
	case "/join":
		if len(parts) != 2 { c.Send("usage: /join <room>"); break }
		old := c.Room
		h.Broadcast(old, "* "+c.Name+" left #"+old, c)
		h.Move(c, parts[1])
		c.Send("joined #" + c.Room)
		h.Broadcast(c.Room, "* "+c.Name+" joined #"+c.Room, c)
	case "/msg":
		if len(parts) < 3 { c.Send("usage: /msg <user> <text>"); break }
		text := strings.Join(parts[2:], " ")
		if !h.Direct(parts[1], fmt.Sprintf("[DM] %s: %s", c.Name, text)) { c.Send("error: user not found") } else { c.Send("[DM -> "+parts[1]+"] "+text) }
	case "/quit":
		c.Send("bye")
		return false
	default:
		c.Send("error: unknown command; use /help")
	}
	return true
}
