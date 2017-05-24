package main

type Progression struct {
	Name  string
	Steps []int
}

var (
	Alternative    Progression = Progression{Name: "Alternative", Steps: []int{5, 3, 0, 4}}
	Canon                      = Progression{Name: "Canon", Steps: []int{0, 4, 5, 2, 3, 0, 3, 4}}
	Cliche                     = Progression{Name: "Cliché", Steps: []int{0, 4, 5, 3}}
	Cliche2                    = Progression{Name: "Cliché 2", Steps: []int{0, 5, 2, 6}}
	Creepy                     = Progression{Name: "Creepy", Steps: []int{0, 5, 3, 4}}
	Creepy2                    = Progression{Name: "Creepy 2", Steps: []int{0, 5, 1, 4}}
	Endless                    = Progression{Name: "Endless", Steps: []int{0, 5, 1, 3}}
	Energetic                  = Progression{Name: "Energetic", Steps: []int{0, 2, 3, 5}}
	Grungy                     = Progression{Name: "Grungy", Steps: []int{0, 3, 2, 5}}
	Memories                   = Progression{Name: "Memories", Steps: []int{0, 3, 0, 4}}
	Rebellious                 = Progression{Name: "Rebellious", Steps: []int{3, 0, 3, 4}}
	Sad                        = Progression{Name: "Sad", Steps: []int{0, 3, 4, 4}}
	Simple                     = Progression{Name: "Simple", Steps: []int{0, 3}}
	Simple2                    = Progression{Name: "Simple 2", Steps: []int{0, 4}}
	TwelveBarBlues             = Progression{Name: "Twelve Bar Blues", Steps: []int{0, 0, 0, 0, 3, 3, 0, 0, 4, 3, 0, 4}}
	Wistful                    = Progression{Name: "Wistful", Steps: []int{0, 0, 3, 5}}
)
