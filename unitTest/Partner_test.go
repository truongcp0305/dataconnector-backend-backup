package unittest

import (
	"data-connector/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (es *Suite) TestTrue() {
	es.True(true)
}

type mockPartnerRepo struct {
	mock.Mock
}

func (m *mockPartnerRepo) GetListPartner() ([]map[string]string, error) {
	args := m.Called()
	result := args.Get(0)
	return result.([]map[string]string), args.Error(1)
}

func TestGetListPartner(t *testing.T) {
	mockdb := new(mockPartnerRepo)
	var dataReturn []map[string]string
	mockdb.On("GetListPartner").Return(append(dataReturn, map[string]string{
		"namePartner": "Misa",
		"config":      "data",
	}), nil)

	//  testSevice := Partner()

}

func TestGetListPartner_pass(t *testing.T) {
	partnerModel := new(model.Partner)
	_, err := partnerModel.GetListPartner()
	if err != nil {
		t.Errorf("cannot get List Partner : %s", err)
	} else {
		t.Logf("oke")
	}
}
