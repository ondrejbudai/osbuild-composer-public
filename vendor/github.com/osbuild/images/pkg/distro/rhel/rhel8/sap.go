package rhel8

import (
	"github.com/osbuild/images/internal/common"
	"github.com/osbuild/images/pkg/distro"
	"github.com/osbuild/images/pkg/osbuild"
)

// sapImageConfig returns the SAP specific ImageConfig data
func sapImageConfig(rd distro.Distro) *distro.ImageConfig {
	ic := &distro.ImageConfig{
		SELinuxConfig: &osbuild.SELinuxConfigStageOptions{
			State: osbuild.SELinuxStatePermissive,
		},
		// RHBZ#1960617
		Tuned: osbuild.NewTunedStageOptions("sap-hana"),
		// RHBZ#1959979
		Tmpfilesd: []*osbuild.TmpfilesdStageOptions{
			osbuild.NewTmpfilesdStageOptions("sap.conf",
				[]osbuild.TmpfilesdConfigLine{
					{
						Type: "x",
						Path: "/tmp/.sap*",
					},
					{
						Type: "x",
						Path: "/tmp/.hdb*lock",
					},
					{
						Type: "x",
						Path: "/tmp/.trex*lock",
					},
				},
			),
		},
		// RHBZ#1959963
		PamLimitsConf: []*osbuild.PamLimitsConfStageOptions{
			osbuild.NewPamLimitsConfStageOptions("99-sap.conf",
				[]osbuild.PamLimitsConfigLine{
					{
						Domain: "@sapsys",
						Type:   osbuild.PamLimitsTypeHard,
						Item:   osbuild.PamLimitsItemNofile,
						Value:  osbuild.PamLimitsValueInt(1048576),
					},
					{
						Domain: "@sapsys",
						Type:   osbuild.PamLimitsTypeSoft,
						Item:   osbuild.PamLimitsItemNofile,
						Value:  osbuild.PamLimitsValueInt(1048576),
					},
					{
						Domain: "@dba",
						Type:   osbuild.PamLimitsTypeHard,
						Item:   osbuild.PamLimitsItemNofile,
						Value:  osbuild.PamLimitsValueInt(1048576),
					},
					{
						Domain: "@dba",
						Type:   osbuild.PamLimitsTypeSoft,
						Item:   osbuild.PamLimitsItemNofile,
						Value:  osbuild.PamLimitsValueInt(1048576),
					},
					{
						Domain: "@sapsys",
						Type:   osbuild.PamLimitsTypeHard,
						Item:   osbuild.PamLimitsItemNproc,
						Value:  osbuild.PamLimitsValueUnlimited,
					},
					{
						Domain: "@sapsys",
						Type:   osbuild.PamLimitsTypeSoft,
						Item:   osbuild.PamLimitsItemNproc,
						Value:  osbuild.PamLimitsValueUnlimited,
					},
					{
						Domain: "@dba",
						Type:   osbuild.PamLimitsTypeHard,
						Item:   osbuild.PamLimitsItemNproc,
						Value:  osbuild.PamLimitsValueUnlimited,
					},
					{
						Domain: "@dba",
						Type:   osbuild.PamLimitsTypeSoft,
						Item:   osbuild.PamLimitsItemNproc,
						Value:  osbuild.PamLimitsValueUnlimited,
					},
				},
			),
		},
		// RHBZ#1959962
		Sysctld: []*osbuild.SysctldStageOptions{
			osbuild.NewSysctldStageOptions("sap.conf",
				[]osbuild.SysctldConfigLine{
					{
						Key:   "kernel.pid_max",
						Value: "4194304",
					},
					{
						Key:   "vm.max_map_count",
						Value: "2147483647",
					},
				},
			),
		},
	}

	if common.VersionLessThan(rd.OsVersion(), "8.10") {
		// E4S/EUS
		if ic.DNFConfig == nil {
			ic.DNFConfig = &distro.DNFConfig{}
		}
		ic.DNFConfig.SetReleaseVerVar = common.ToPtr(true)
	}

	return ic
}
