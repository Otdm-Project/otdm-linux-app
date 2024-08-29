package commands

import (
	"otdm-package/src/utils"
)

func RunUp() {
	utils.CallRefresh()
	utils.CallWebsocket()
	utils.CallBoot()
	utils.CallNetwork()
}
