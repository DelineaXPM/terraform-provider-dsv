module github.com/thycotic/terraform-provider-dsv

require (
	github.com/hashicorp/terraform v0.12.14
	github.com/thycotic/dsv-sdk-go v1.0.1
	github.com/ulikunitz/xz v0.5.10 // indirect
)

// replace github.com/thycotic/dsv-sdk-go => ../dsv-sdk-go

go 1.13
