package log

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/task4233/note-v2-terraform/client"
)

type resourceLogType struct{}

func (r resourceLogType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"log": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"body": {
						Type:     types.StringType,
						Required: true,
					},
				}),
			},
		},
	}, nil
}

type resourceLog struct {
	p provider
}

func (r resourceLogType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceLog{
		p: *(p.(*provider)),
	}, nil
}

func (r resourceLog) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	f, err := os.Open("~/test")
	if err != nil {
		return
	}
	defer f.Close()

	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	f.Write([]byte("0"))

	// Retrieve values from plan
	var plan Log
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Body.Null || plan.Body.Unknown {
		resp.Diagnostics.AddError(
			"Plan is not corrent value",
			"body must not be null or unknown",
		)
		return
	}

	f.Write([]byte("1"))

	// Generate API request body from plan
	log := client.OrderLog{
		Log: client.Log{
			Body: plan.Body.Value,
		},
	}

	f.Write([]byte("2"))

	gotLog, err := r.p.client.CreateLog(ctx, &log)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating log",
			fmt.Sprintf("Could not create log: %s", err.Error()),
		)
	}

	f.Write([]byte("3"))

	result := OrderLog{
		Log: Log{
			types.String{
				Unknown: false,
				Null:    false,
				Value:   gotLog.Body,
			},
		},
	}

	f.Write([]byte("4"))

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Fprintf(os.Stderr, "[Create]\n")
}

func (r resourceLog) Read(context.Context, tfsdk.ReadResourceRequest, *tfsdk.ReadResourceResponse) {
	fmt.Fprintf(os.Stderr, "[Read]\n")
}

func (r resourceLog) Update(context.Context, tfsdk.UpdateResourceRequest, *tfsdk.UpdateResourceResponse) {
	fmt.Fprintf(os.Stderr, "[Update]\n")
}

func (r resourceLog) Delete(context.Context, tfsdk.DeleteResourceRequest, *tfsdk.DeleteResourceResponse) {
	fmt.Fprintf(os.Stderr, "[Delete]\n")
}

func (r resourceLog) ImportState(context.Context, tfsdk.ImportResourceStateRequest, *tfsdk.ImportResourceStateResponse) {
	fmt.Fprintf(os.Stderr, "[ImportState]\n")
}
