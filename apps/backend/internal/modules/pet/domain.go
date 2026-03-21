package pet

type PetCatalogItem struct {
	PetID     int64  `json:"pet_id"`
	Name      string `json:"name"`
	Rarity    int    `json:"rarity"`
	Element   string `json:"element"`
	Locked    bool   `json:"locked"`
	Owned     bool   `json:"owned"`
	Power     int    `json:"power"`
	Portrait  string `json:"portrait"`
	StoryText string `json:"story_text"`
}

type PlayerPet struct {
	PlayerID int64 `json:"player_id"`
	PetID    int64 `json:"pet_id"`
	Level    int   `json:"level"`
	Star     int   `json:"star"`
	Power    int   `json:"power"`
}

type PetTeam struct {
	PlayerID int64   `json:"player_id"`
	PetIDs   []int64 `json:"pet_ids"`
}
