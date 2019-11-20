package main

import (
	SDK "github.com/advwacloud/WISEPaaS.SCADA.Go.SDK"
	"fmt"
)

func main() {
	quit := make(chan bool)

	options := SDK.NewEdgeAgentOptions()
	options.ScadaID = "9eb2bbe4-6833-45ff-b884-297be549c5cc"
	options.ConnectType = SDK.ConnectType["DCCS"]
	options.DCCS.Key = "9ba5b0eace39c528dd6c095e15de2ere"
	options.DCCS.URL = "https://api-dccs.wise-paas.com/"

	agent := SDK.NewAgent(options)
	agent.SetOnConnectHandler(func(a SDK.Agent) {
		fmt.Println("connect successfully")
		config := generateConfig()
		action := SDK.Action["Delete"]
		_ = agent.UploadConfig(action, config)
	})
	agent.SetOnDisconnectHandler(func(a SDK.Agent) {
		fmt.Println("disconnect successfully")
	})
	agent.SetOnMessageReceiveHandler(func(res SDK.MessageReceivedEventArgs) {
		fmt.Println(res)
	})
	agent.Connect()

	<-quit
}

func generateConfig() SDK.EdgeConfig {
	scadaConfig := generateScadaConfig()
	edgeConfig := SDK.EdgeConfig{
		Scada: scadaConfig,
	}
	return edgeConfig
}

func generateScadaConfig() SDK.ScadaConfig {
	var scadaName = "Test_Scada"
	var deviceNum = 1

	scadaConfig := SDK.NewScadaConfig(scadaName)

	for idx := 0; idx < deviceNum; idx++ {
		config := generateDeviceConfig(idx + 1)
		scadaConfig.DeviceList = append(scadaConfig.DeviceList, config)
	}

	return scadaConfig
}

func generateDeviceConfig(idx int) SDK.DeviceConfig {
	var deviceID = fmt.Sprintf("%s%d", "Device", idx)
	var analogNum = 0
	var discreteNum = 0
	var textNum = 0

	deviceConfig := SDK.NewDeviceConfig(deviceID)

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
	return config
}

func generateDiscreteConfig(idx int) SDK.DiscreteTagConfig {
	var tagName = fmt.Sprintf("%s%d", "DTag", idx)
	config := SDK.NewDiscreteTagConfig(tagName)
	return config
}

func generateTextConfig(idx int) SDK.TextTagConfig {
	var tagName = fmt.Sprintf("%s%d", "TTag", idx)
	config := SDK.NewTextTagConfig(tagName)
	return config
}
