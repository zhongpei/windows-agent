package funcs

import (
	"github.com/zhongpei/windows-agent/g"
	"github.com/open-falcon/common/model"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				SwapMemMetrics,
				DeviceMetrics,
				DiskIOMetrics,
				TcpipMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				PortMetrics,
				ProcMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				iisMetrics,
				mssqlMetrics,
			},
			Interval: interval,
		},
	}
}
