module github.com/thycotic/terraform-provider-dsv

require (
	github.com/hashicorp/terraform v0.12.14
	github.com/thycotic/dsv-sdk-go v0.0.0-20200116184609-53e6e5a3ba69
)

replace github.com/thycotic/dsv-sdk-go => ../dsv-sdk-go

go 1.13
