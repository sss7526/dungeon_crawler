package main

import (
	"fmt"
	gloss "github.com/charmbracelet/lipgloss"
)

type Renderable interface {
	View(theme Theme) string
}
// CharacterInfo represents core information about the player/character
type CharacterInfo struct {
	Name 		string
	Class		string
	Level		int
	Experience 	int 
	Gold 		int
}

func (info CharacterInfo) View(theme Theme) string {
	// sectionTitle := theme.MenuOptionStyle.Bold(true).Render("Character")
	content := gloss.JoinVertical(gloss.Left,
		theme.MenuOptionStyle.Bold(true).PaddingLeft(0).PaddingBottom(1).Underline(true).Render("Character"),
		renderKeyValue(theme, "Name", info.Name),
		renderKeyValue(theme, "Class", info.Class),
		renderKeyValue(theme, "Level", fmt.Sprintf("%d", info.Level)),
		renderKeyValue(theme, "Experience", fmt.Sprintf("%d", info.Experience)),
		renderKeyValue(theme, "Gold", fmt.Sprintf("%d", info.Gold)),
	)
	// return gloss.JoinVertical(gloss.Left, sectionTitle, theme.BorderStyle.Render(content))
	return theme.BorderStyle.Render(content)
}

// Attribute represents an individual attribute
type Attribute struct {
	Current 	int
	Max 		int
}

func (attr Attribute) View(theme Theme, name string) string {
	bar := theme.ProgressBar.ViewAs(float64(attr.Current) / float64(attr.Max))
	formattedName := theme.AttributeStyle.Width(10).Render(name+":")
	return gloss.NewStyle().
		Width(70).
		PaddingRight(2).
		PaddingBottom(1).
		Render(fmt.Sprintf("%-10s %s [%d / %d]", formattedName, bar, attr.Current, attr.Max))
}

// Attributes represents player specific attributes
type Attributes struct {
	Strength	Attribute
	Agility 	Attribute
	Intelect 	Attribute
	Endurance	Attribute
	Luck		Attribute
}

func (a Attributes) View(theme Theme) string {
	// sectionTitle := theme.MenuOptionStyle.Bold(true).Render("Attributes")
	content := gloss.JoinVertical(gloss.Left,
		theme.MenuOptionStyle.Bold(true).PaddingLeft(0).PaddingBottom(1).Underline(true).Render("Attributes"),
		a.Strength.View(theme, "Strength"),
		a.Agility.View(theme, "Agility"),
		a.Intelect.View(theme, "Intelect"),
		a.Endurance.View(theme, "Endurance"),
		a.Luck.View(theme, "Luck"),
	)
	// return gloss.JoinVertical(gloss.Left, sectionTitle, theme.BorderStyle.Width(80).Align(gloss.Left).Render(content))
	return theme.BorderStyle.Width(80).Align(gloss.Left).Render(content)
}



// Stats are derived or directly modifiable
type Stats struct {
	Health 		Attribute
	Mana 		Attribute
	Stamina		Attribute
	Damage 		int
	Defense		int
	CritRate	float64
}

type Equipment struct {
	Weapon		string
	Armor		string
	Trinkets 	[]string
}

func (e Equipment) View(theme Theme) string {
	sectionTitle := theme.MenuOptionStyle.Bold(true).PaddingLeft(0).PaddingBottom(1).Underline(true).Render("Equipment")
	var items []string
	items = append(items, sectionTitle)
	items = append(items, renderKeyValue(theme, "Weapon", e.Weapon))
	items = append(items, renderKeyValue(theme, "Armor", e.Armor))
	for _, trinket := range e.Trinkets {
		items = append(items, gloss.NewStyle().PaddingLeft(2).Render("- "+trinket))
	}

	content := gloss.JoinVertical(gloss.Left, items...)
	// return gloss.JoinVertical(gloss.Left, sectionTitle, theme.BorderStyle.Render(content))
	return theme.BorderStyle.Render(content)
}

type Player struct {
	Info  		CharacterInfo
	Attributes	Attributes
	Stats 		Stats
	Inventory	[]string
	Equipment	Equipment
	Gold		int
}

func (p Player) View(theme Theme) string {
	return gloss.JoinVertical(gloss.Left,
		p.Info.View(theme),
		p.Attributes.View(theme),
		p.Equipment.View(theme),
	)
}

func newAttribute(current, max int) Attribute {
	return Attribute{Current: current, Max: max}
}

func newTestPlayer() *Player {
	return &Player{
		Info: CharacterInfo {
			Name:		"Sman",
			Class:		"Warrior",
			Level: 		12,
			Experience: 5600,
		},
		Attributes: Attributes{
			Strength:	newAttribute(8, 10),
			Agility:	newAttribute(7, 10),
			Intelect: 	newAttribute(2, 10),
			Endurance:	newAttribute(9, 10),
			Luck:		newAttribute(4, 10),
		},
		Stats: Stats{
			Health:		newAttribute(85, 100),
			Mana:		newAttribute(30, 50),
			Stamina:	newAttribute(40, 50),
			Damage:		25,
			Defense: 	10,
			CritRate:	0.15,
		},
		Inventory: []string{
			"Health Potion",
			"Mana Potion",
			"Antidote",
		},
		Equipment:	Equipment{
			Weapon:		"Sword of Flame",
			Armor:		"Darksteel Shield",
			Trinkets:	[]string{"Ring of Vitality", "Amulet of Luck"},
		},
		Gold: 200,
	}
}

func renderKeyValue(theme Theme, key, value string) string {
	formattedKey := theme.AttributeStyle.Render(key+":")
	formattedValue := gloss.NewStyle().Render(value)
	return gloss.JoinHorizontal(gloss.Left, formattedKey, formattedValue)
}