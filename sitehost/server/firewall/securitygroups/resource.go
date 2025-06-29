// Package securitygroups provides the functions to create a Security Group resource via SiteHost API.
package securitygroups

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server/firewall/securitygroups"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for security groups.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new security group.
func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	opts := securitygroups.AddRequest{
		Label: fmt.Sprint(d.Get("label")),
	}

	res, err := client.Add(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating security group: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error creating security group: %s", res.Msg)
	}

	// Set data
	d.SetId(res.Return.Name)
	if err := d.Set("name", res.Return.Name); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Security Group Name: %s", d.Id())

	if d.Get("rules_in") != nil || d.Get("rules_out") != nil {
		if err := updateResource(ctx, d, meta); err != nil {
			return err
		}
	}

	return nil
}

// updateResource is a function to update a security group.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)
	opts := securitygroups.UpdateRequest{
		Name: d.Id(),
		Params: securitygroups.ParamsOptions{
			Label:    fmt.Sprint(d.Get("label")),
			RulesIn:  processRules(d.Get("rules_in"), ruleIn),
			RulesOut: processRules(d.Get("rules_out"), ruleOut),
		},
	}

	response, err := client.Update(ctx, opts)
	if err != nil {
		return diag.Errorf("Error updating security group: %s", err)
	}

	if !response.Status {
		return diag.Errorf("Error updating security group: %s", response.Msg)
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
		return diag.FromErr(err)
	}

	return readResource(ctx, d, meta)
}

// processRules processes the rules for a security group.
func processRules(rulesRaw interface{}, direction string) []securitygroups.UpdateRequestRule {
	rules := make([]securitygroups.UpdateRequestRule, 0)
	rulesSlice, ok := rulesRaw.([]interface{})
	if !ok {
		return rules
	}

	for _, ruleRaw := range rulesSlice {
		ruleMap, ok := ruleRaw.(map[string]interface{})
		if !ok {
			continue
		}

		enabled := true
		if e, ok := ruleMap["enabled"].(bool); ok {
			enabled = e
		}

		var ip string
		switch direction {
		case ruleIn:
			ip = fmt.Sprint(ruleMap["src_ip"])
		case ruleOut:
			ip = fmt.Sprint(ruleMap["dest_ip"])
		}

		rules = append(rules, securitygroups.UpdateRequestRule{
			Enabled:         enabled,
			IP:              ip,
			Action:          fmt.Sprint(ruleMap["action"]),
			Protocol:        fmt.Sprint(ruleMap["protocol"]),
			DestinationPort: fmt.Sprint(ruleMap["dest_port"]),
		})
	}
	return rules
}

// deleteResource is a function to delete a security group.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	resp, err := client.Delete(ctx, securitygroups.DeleteRequest{
		Name: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error deleting security group: %s", err)
	}
	if !resp.Status {
		return diag.Errorf("Error deleting security group: %s", resp.Msg)
	}

	return nil
}

// readResource is a function to read a security group.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)
	resp, err := client.Get(ctx, securitygroups.GetRequest{
		Name: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error reading security group: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error reading security group: %s", resp.Msg)
	}

	if err := d.Set("label", resp.Return.Label); err != nil {
		return diag.FromErr(err)
	}

	// Set inbound rules
	if err := readRules(d, resp.Return.Rules.In, ruleIn); err != nil {
		return err
	}

	// Set outbound rules
	if err := readRules(d, resp.Return.Rules.Out, ruleOut); err != nil {
		return err
	}

	return nil
}

// readRules reads the rules for a security group.
func readRules(d *schema.ResourceData, rulesRaw []securitygroups.Rule, direction string) diag.Diagnostics {
	rules := make([]map[string]interface{}, len(rulesRaw))
	for i, ruleRaw := range rulesRaw {
		rule := map[string]interface{}{
			"enabled":   ruleRaw.Enabled,
			"action":    ruleRaw.Action,
			"protocol":  ruleRaw.Protocol,
			"dest_port": fmt.Sprint(ruleRaw.DestPort),
		}
		switch direction {
		case ruleIn:
			rule["src_ip"] = ruleRaw.SrcIP
		case ruleOut:
			rule["dest_ip"] = ruleRaw.DestIP
		}
		rules[i] = rule
	}

	if err := d.Set("rules_"+direction, rules); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
