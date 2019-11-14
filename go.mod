module terraform-provider-dsv

require (
	github.com/amigus/dsv-sdk-go v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform v0.12.14
)

replace github.com/amigus/dsv-sdk-go => ../dsv-sdk-go

go 1.13
