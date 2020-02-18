package main

import (
	SDK "github.com/advwacloud/WISEPaaS.DataHub.Edge.Go.SDK"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	quit := make(chan bool)

	options := SDK.NewEdgeAgentOptions()
	options.NodeID = "9eb2bbe4-6833-45ff-b884-297be549c5cc"
	options.ConnectType = SDK.ConnectType["DCCS"]
	options.DCCS.Key = "9ba5b0eace39c528dd6c095e15de2ere"
	options.DCCS.URL = "https://api-dccs.wise-paas.com/"

	interval := 1
	var timer chan bool = nil

	agent := SDK.NewAgent(options)
	agent.SetOnConnectHandler(func(a SDK.Agent) {
		fmt.Println("connect successfully")

		config := generateConfig()
		action := SDK.Action["Update"]
		_ = agent.UploadConfig(action, config)

		timer = setInterval(func() {
			data := generateData()
			_ = agent.SendData(data)
		}, interval, true)

	})
	agent.SetOnDisconnectHandler(func(a SDK.Agent) {
		fmt.Println("disconnect successfully")
	})
	agent.SetOnMessageReceiveHandler(func(res SDK.MessageReceivedEventArgs) {
		fmt.Println(res)
	})

	fmt.Println(agent.IsConnected())
	agent.Connect()
	fmt.Println(agent.IsConnected())
	/* agent.Disconnect() */
	fmt.Println(agent.IsConnected())

	<-quit
}

func generateData() SDK.WriteDataMessage {
	deviceNum := 1
	msg := SDK.WriteDataMessage{
		Timestamp: time.Now(),
	}

	for idx := 0; idx < deviceNum; idx++ {
		analogNum := 3
		discreteNum := 2
		textNum := 1
		device := SDK.Device{
			ID: fmt.Sprintf("%s%d", "Device", idx+1),
		}
		for num := 0; num < analogNum; num++ {
			v := make(map[string]interface{})
			for i := 0; i < 3; i++ {
				v[fmt.Sprintf("%d", i)] = rand.Float64()
			}
			t := SDK.Tag{
				Name:  fmt.Sprintf("%s%d", "ATag", num+1),
				Value: v,
			}
			device.TagList = append(device.TagList, t)
		}
		for num := 0; num < discreteNum; num++ {
			v := make(map[string]interface{})
			for i := 0; i < 2; i++ {
				v[fmt.Sprintf("%d", i)] = rand.Intn(7)
			}
			t := SDK.Tag{
				Name:  fmt.Sprintf("%s%d", "DTag", num+1),
				Value: v,
			}
			device.TagList = append(device.TagList, t)
		}
		for num := 0; num < textNum; num++ {
			v := make(map[string]interface{})
			for i := 0; i < 1; i++ {
				v[fmt.Sprintf("%d", i)] = fmt.Sprintf("%s%f", "str", rand.Float64())
			}
			t := SDK.Tag{
				Name:  fmt.Sprintf("%s%d", "TTag", num+1),
				Value: v,
			}
			device.TagList = append(device.TagList, t)
		}
		msg.DeviceList = append(msg.DeviceList, device)
	}
	fmt.Println(msg)
	return msg
}

func setInterval(someFunc func(), seconds int, async bool) chan bool {
	interval := time.Duration(seconds) * time.Second
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				close(clear)
				return
			}
		}
	}()
	return clear
}

func generateConfig() SDK.EdgeConfig {
	nodeConfig := generateNodeConfig()
	edgeConfig := SDK.EdgeConfig{
		Node: nodeConfig,
	}
	return edgeConfig
}

func generateNodeConfig() SDK.NodeConfig {
	var nodeName = "Test_Node"
	var deviceNum = 1

	nodeConfig := SDK.NewNodeConfig(nodeName)
	nodeConfig.SetDescription("For Test")
	nodeConfig.SetNodeType(SDK.EdgeType["Gateway"])

	for idx := 0; idx < deviceNum; idx++ {
		config := generateDeviceConfig(idx + 1)
		nodeConfig.DeviceList = append(nodeConfig.DeviceList, config)
	}

	return nodeConfig
}

func generateDeviceConfig(idx int) SDK.DeviceConfig {
	var deviceID = fmt.Sprintf("%s%d", "Device", idx)
	var analogNum = 3
	var discreteNum = 2
	var textNum = 1

	deviceConfig := SDK.NewDeviceConfig(deviceID)
	deviceConfig.SetDeviceName(fmt.Sprintf("%s%d", "Device", idx))
	deviceConfig.SetDeviceType("Smart Device")
	deviceConfig.SetDeviceDescription(fmt.Sprintf("%s %d", "Device ", idx))

	for idx := 0; idx < analogNum; idx++ {
		config := generateAnalogConfig(idx + 1)
		deviceConfig.AnalogTagList = append(deviceConfig.AnalogTagList, config)
	}
	for idx := 0; idx < discreteNum; idx++ {
		config := generateDiscreteConfig(idx + 1)
		deviceConfig.DiscreteTagList = append(deviceConfig.DiscreteTagList, config)
	}
	for idx := 0; idx < textNum; idx++ {
		config := generateTextConfig(idx + 1)
		deviceConfig.TextTagList = append(deviceConfig.TextTagList, config)
	}
	return deviceConfig
}

func generateAnalogConfig(idx int) SDK.AnalogTagConfig {
	var tagName = fmt.Sprintf("%s%d", "ATag", idx)

	config := SDK.NewAnaglogTagConfig(tagName)
	config.SetTagDescription(fmt.Sprintf("%s %d", "ATag", idx))
	config.SetTagReadOnly(false)
	config.SetTagArraySize(3)
	config.SetTagSpanHigh(1000.0)
	config.SetTagSpanLow(0.0)
	config.SetTagEngineerUnit("")
	config.SetTagIntegerDisplayFormat(4)
	config.SetTagFractionDisplayFormat(2)

	return config
}

func generateDiscreteConfig(idx int) SDK.DiscreteTagConfig {
	var tagName = fmt.Sprintf("%s%d", "DTag", idx)

	config := SDK.NewDiscreteTagConfig(tagName)
	config.SetTagDescription(fmt.Sprintf("%s %d", "DTag ", idx))
	config.SetTagArraySize(2)
	config.SetTagReadOnly(true)
	config.SetTagState0("No")
	config.SetTagState1("Yes")

	return config
}

func generateTextConfig(idx int) SDK.TextTagConfig {
	var tagName = fmt.Sprintf("%s%d", "TTag", idx)

	config := SDK.NewTextTagConfig(tagName)
	config.SetTagDescription(fmt.Sprintf("%s %d", "TTagx", idx))
	config.SetTagReadOnly(false)
	config.SetTagArraySize(1)

	return config
}
