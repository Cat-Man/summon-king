package growth

type Wallet struct {
	PlayerID       int64 `json:"player_id"`
	SpiritPower    int64 `json:"spirit_power"`
	SpiritFreeWash int   `json:"spirit_free_wash"`
	BoneLevel      int   `json:"bone_level"`
	SoulPieces     int   `json:"soul_pieces"`
	ManorPlots     int   `json:"manor_plots"`
}

type WashOption struct {
	SpiritID int64 `json:"spirit_id"`
}

type Bone struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Spirit struct {
	Name      string `json:"name"`
	WashCount int    `json:"wash_count"`
}

type Soul struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

type ManorPlot struct {
	PlotID int64  `json:"plot_id"`
	State  string `json:"state"`
}

type ManorHarvestResult struct {
	Message string      `json:"message"`
	Plots   []ManorPlot `json:"plots"`
}
