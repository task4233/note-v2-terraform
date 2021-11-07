package log

import "github.com/hashicorp/terraform-plugin-framework/types"

type OrderLog struct {
	Body types.String `tfsdk:"body"`
}
