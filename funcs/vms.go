package funcs

import (
	"github.com/zhongpei/windows-agent/g"
	"github.com/open-falcon/common/model"
	ps "github.com/zhongpei/go-powershell"
	"github.com/zhongpei/go-powershell/backend"

	"strconv"
)

func RunPs(cmd string )(string,error){
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		return "",err
	}

	defer shell.Exit()

	// ... and interact with it
	stdout, _, err := shell.Execute(cmd)
	if err != nil {

		return "",err
	}

	return stdout,err
}

func VMCountMetrics()[]*model.MetricValue {
	var totalCount = 0
	var offCount = 0
	var runningCount = 0

	total,err := RunPs("$VMs=Get-VM;$VMs.Count")
	if err==nil{
		totalCount,_ = strconv.Atoi(total)
	}
	off ,err := RunPs(" $VMs=Get-VM |where {$_.State -eq 'Off'}; $VMs.Count")
	if err==nil{
		offCount,_ = strconv.Atoi(off)
	}
	running ,err := RunPs(" $VMs=Get-VM |where {$_.State -eq 'Running'}; $VMs.Count")
	if err==nil{
		runningCount,_ = strconv.Atoi(running)
	}

	return []*model.MetricValue{
		CounterValue("vms.count.total", totalCount),
		CounterValue("vms.count.off", offCount),
		CounterValue("vms.count.running", runningCount),
	}
}
