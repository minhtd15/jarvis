package sports

type SportsResponse struct {
	SportId   int    `json:"sportId"`
	SportName string `json:"sportName"`
	SportUlr  string `json:"sportUrl"`
}
