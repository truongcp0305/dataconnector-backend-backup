package model

// type Repository struct {
// }
// type RespositoryInterface interface {
// 	GetListPartner() ([]map[string]string, error)
// }

// var (
// 	repository Repository
// )

// func NewPartnerModel(repo Repository) {
// 	repository = repo
// 	return
// }

type Partner struct {
	Id          string `json:"id" db:"id" param:"id" form:"id" type:"text" primary:"true"`
	NamePartner string `json:"name_partner" db:"name_partner" param:"name_partner" form:"name_partner" type:"text"`
	Config      string `json:"config" db:"config" param:"config" form:"config" type:"text"`
}

func GetTableName() string {
	return "partner"
}

type PartnerModelInterface interface {
	InitByArray(data map[string]string)
	GetListPartner() ([]map[string]string, error)
}

func (partner *Partner) InitByArray(data map[string]string) {
	partner.Id = data["id"]
	partner.NamePartner = data["name"]
	partner.Config = data["detail"]
}

func (partner Partner) GetListPartner() ([]map[string]string, error) {
	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM partner")
	if err != nil {
		return nil, err
	} else {
		value := PackageData(rows)
		return value, nil
	}
}
