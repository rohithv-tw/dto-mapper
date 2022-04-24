package dto

type TargetDto struct {
	RandomFieldA string   `map_from:"name.first"`
	RandomFieldB string   `map_from:"name.middle"`
	RandomFieldC string   `map_from:"name.last"`
	RandomFieldD []string `map_from:"friends.#.first"`
	RandomFieldE int      `map_from:"age"`
	// RandomFieldF time.Time `map_from:"created_at"`
}
