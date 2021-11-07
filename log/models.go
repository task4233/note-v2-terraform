package log

import "github.com/hashicorp/terraform-plugin-framework/types"

type OrderLog struct {
	Log Log `tfsdk:"log"`
}

type Log struct {
	Body types.String `tfsdk:"body"`
}
