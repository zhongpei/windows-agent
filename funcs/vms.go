package funcs

import (
	"github.com/open-falcon/common/model"
	ps "github.com/zhongpei/go-powershell"
	"github.com/zhongpei/go-powershell/backend"
	"strconv"
	"strings"
	"fmt"
)

func RunPs(cmd string) (string, error) {
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		return "", err
	}

	defer shell.Exit()

	// ... and interact with it
	fmt.Println("cmd", cmd)

	stdout, _, err := shell.Execute(cmd)

	if err != nil {

		return "", err
	}
	fmt.Println("out", stdout, err)

	return stdout, err
}

func VMCountMetrics() []*model.MetricValue {
	var totalCount = 0
	var offCount = 0
	var runningCount = 0

	total, err := RunPs("$VMs=Get-VM;$VMs.Count")
	if err == nil {
		totalCount, _ = strconv.Atoi(strings.TrimSpace(total))
	}
	off, err := RunPs(" $VMs=Get-VM |where {$_.State -eq 'Off'}; $VMs.Count")
	if err == nil {
		offCount, _ = strconv.Atoi(strings.TrimSpace(off))
	}
	running, err := RunPs(" $VMs=Get-VM |where {$_.State -eq 'Running'}; $VMs.Count")
	if err == nil {
		runningCount, _ = strconv.Atoi(strings.TrimSpace(running))
	}
	fmt.Println("VMcount", totalCount, offCount, running)
	return []*model.MetricValue{
		GaugeValue("vms.count.total", totalCount),
		GaugeValue("vms.count.off", offCount),
		GaugeValue("vms.count.running", runningCount),
	}
}
func VMHostNumaNodeMetrics() (L []*model.MetricValue) {
	nodeIds, err := RunPs("$VMs=Get-VMHostNumaNode;$VMs.NodeId")
	if err != nil {
		return []*model.MetricValue{}
	}
	ids := strings.Split(nodeIds, "\n")

	fmt.Println("ids", ids)

	for _, id := range ids {
		i := strings.TrimSpace(id)
		if len(i) == 0 {
			continue
		}

		var memfreeGB = 0
		var memtotalGB = 0
		var cputotalCount = 0
		memfree, err := RunPs(fmt.Sprintf("$VMs=Get-VMHostNumaNode;$VMs[%s].MemoryAvailable;", i))
		if err == nil {
			memfreeGB, _ = strconv.Atoi(strings.TrimSpace(memfree))
		}
		memtotal, err := RunPs(fmt.Sprintf("$VMs=Get-VMHostNumaNode;$VMs[%s].MemoryTotal;", i))
		if err == nil {
			memtotalGB, _ = strconv.Atoi(strings.TrimSpace(memtotal))
		}
		cputotal, err := RunPs(fmt.Sprintf("$VMs=Get-VMHostNumaNode;$VMs[%s].ProcessorsAvailability.Count;", i))
		if err == nil {
			cputotalCount, _ = strconv.Atoi(strings.TrimSpace(cputotal))
		}

		nodename := fmt.Sprintf("node%s",i)
		L = append(L, GaugeValue("numa.memfree."+nodename, memfreeGB/1024))
		L = append(L, GaugeValue("numa.memtotal."+nodename, memtotalGB/1024))
		L = append(L, GaugeValue("numa.cputotal."+nodename, cputotalCount))
		fmt.Println("numa",id,i,nodename, memfreeGB/1024, memtotalGB/1024, cputotalCount)

	}

	fmt.Println("numa metrics len", L, len(L))
	return
}
