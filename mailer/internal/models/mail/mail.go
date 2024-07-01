package mail

// TODO: figure out better place to place models.
// It's not right place I guess. Also need to think how can
// I possibly avoid dublication in naming.

type Mail struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}
