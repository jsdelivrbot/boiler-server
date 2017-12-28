package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/models"
	"fmt"
	"errors"
)

type AddressController struct {
	goazure.Controller
}

var AddrCtl *AddressController = &AddressController{}

func (addrCtl *AddressController) GetAddress(local *models.Location, address string) (*models.Address, error) {
	var err error
	var addr models.Address
	if local == nil {
		addr= models.Address{
			MyAddress: models.MyAddress{
				Location: LocalCtl.GetLocationWithId(0),
				Address: "",
			},
		}
		err = errors.New("Location & Address Can not be nil!\n Set Default")
	} else {
		addr = models.Address{
			MyAddress: models.MyAddress{
				Location: local,
				Address: address,
			},
		}
	}

	err = DataCtl.AddData(&addr, false, "Location", "Address")
	if err != nil {
		fmt.Println("Read Address Error: ", err)
	}

	return &addr, err
}
