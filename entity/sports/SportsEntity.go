package sports

type SportsEntity struct {
	SportId   int    `db:"SPORT_ID"`
	SportName string `db:"SPORT_NAME"`
	SportUlr  string `db:"SPORT_URL"`
}
