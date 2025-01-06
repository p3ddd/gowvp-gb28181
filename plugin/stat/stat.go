package stat

import (
	"log/slog"
	"time"

	"github.com/ixugo/goweb/pkg/orm"
	"github.com/shirou/gopsutil/net"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type UsageStat struct {
	Name      string `json:"name"`       // 路径
	Unit      string `json:"unit"`       // 单位
	Size      string `json:"size"`       // 大小
	FreeSpace string `json:"free_space"` // 可用空间
	Used      string `json:"used"`       // 已使用
	Percent   string `json:"percent"`    // 比率
	Threshold string `json:"threshold"`  // 阈值
}

const (
	TopQueneCap = 30
)

var (
	memData = NewCircleQueue(TopQueneCap)
	cpuData = NewCircleQueue(TopQueneCap)
	// netUpData         = NewCircleQueue(TopQueneCap)
	netData           = NewCircleQueue(TopQueneCap)
	currentMem        float64
	currentCPU        float64
	currentMainDisk   uint64
	totalMainDisk     uint64
	currentKernelDisk float64
	totalKernelDisk   uint64
)

func GetCurrentMem() float64 {
	return currentMem
}

func GetCurrentCPU() float64 {
	return currentCPU
}

func GetCurrentMainDisk() uint64 {
	return currentMainDisk
}

func GetTotalMainDisk() uint64 {
	return totalMainDisk
}

func GetCurrentKernelDisk() float64 {
	return currentKernelDisk
}

func GetTotalKernelDisk() uint64 {
	return totalKernelDisk
}

func GetMemData() []PercentData {
	return memData.Range()
}

func GetCPUData() []PercentData {
	return cpuData.Range()
}

func GetNetData() []PercentData {
	return netData.Range()
}

func LoadTop(path string, fn func(map[string]any)) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		// cpu
		now := orm.Now()
		cpu, err := cpu.Percent(0, false)
		if err != nil && err.Error() != `not implemented yet` {
			slog.Error("LoadTop cpu", "err", err)
		}
		if len(cpu) > 0 {
			cpuData.Push(PercentData{Time: now, Used: cpu[0]})
		}

		// memory
		mem, err := mem.VirtualMemory()
		if err != nil && err.Error() != `not implemented yet` {
			slog.Error("LoadTop VirtualMemory", "err", err)
		}
		if mem != nil {
			memData.Push(PercentData{Time: now, Used: mem.UsedPercent})
		}

		// net flow
		n1, _ := net.IOCounters(false)
		time.Sleep(1000 * time.Millisecond)
		n2, _ := net.IOCounters(false)
		if len(n1) > 0 && len(n2) > 0 {
			netData.Push(PercentData{
				Time: now, Up: float64(n2[0].BytesSent-n1[0].BytesSent) * 8,
				Down: float64(n2[0].BytesRecv-n1[0].BytesRecv) * 8,
			})
		}
		// 当前值统计
		if mem != nil {
			currentMem = mem.UsedPercent
			if len(cpu) > 0 {
				currentCPU = cpu[0]
			}
		}

		if diskres, err := disk.Usage(path); err == nil {
			currentMainDisk = diskres.Used
			totalMainDisk = diskres.Total
		}
		fn(map[string]any{
			"mem": memData.Last(),
			"cpu": cpuData.Last(),
			"net": netData.Last(),
			// "netup":   netUpData.Last(),
			// "netdown": netDownData.Last(),
			"disk": []map[string]any{
				{
					"name":  path,
					"used":  currentMainDisk,
					"total": totalMainDisk,
				},
			},
		})

		ticker.Reset(200 * time.Millisecond)
	}
}
