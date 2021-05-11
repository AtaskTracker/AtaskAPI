package dto

type Bearer struct {
	Token string `bson:"id_token" json:"id_token"`
}
