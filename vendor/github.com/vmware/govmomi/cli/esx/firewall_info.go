/*
Copyright (c) 2015-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package esx

import "context"

type FirewallInfo struct {
	Loaded        bool   `json:"loaded"`
	Enabled       bool   `json:"enabled"`
	DefaultAction string `json:"defaultAction"`
}

// GetFirewallInfo via 'esxcli network firewall get'
// The HostFirewallSystem type does not expose this data.
// This helper can be useful in particular to determine if the firewall is enabled or disabled.
func (x *Executor) GetFirewallInfo(ctx context.Context) (*FirewallInfo, error) {
	res, err := x.Run(ctx, []string{"network", "firewall", "get"})
	if err != nil {
		return nil, err
	}

	info := &FirewallInfo{
		Loaded:        res.Values[0]["Loaded"][0] == "true",
		Enabled:       res.Values[0]["Enabled"][0] == "true",
		DefaultAction: res.Values[0]["DefaultAction"][0],
	}

	return info, nil
}
