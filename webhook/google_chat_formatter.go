package webhook

var googleChatLabel = map[int]string{
	appTitle:      `❤️ Webhook Docker __APP_VERSION__ (Built: __APP_BUILT_DATE__)    __NEW_LINE__🗂️ ` + "`__GITHUB_URL__ ` " + `    __NEW_LINE__🚢 __DOCKER_HUB_URL__    __NEW_LINE__`,
	startupTitle1: `🐳 Docker Version: *__DOCKER_VERSION__* __TAB__ApiVersion: __DOCKER_API_VERSION__  __NEW_LINE__🌐 Node: ` + "`__NODE_NAME__` __NEW_LINE__",
	startupTitle2: `🛠️ Os: __OS__  __TAB__Kernel: __KERNEL_VERSION__`,

	groupTitle:             `__NEW_LINE__🚛 ` + "`__NAME__`__NODE_NAME__" + ` __NEW_LINE__💿 __IMAGE__    __NEW_LINE__🕰️ __TIME__   __NEW_LINE__`,
	groupFooter:            `__NEW_LINE__🆔 __ID__`,
	groupFooterWithInspect: `🔍 __INSPECT_URL__ __NEW_LINE__🆔 __ID__`,

	containerDefault: `- 📦 *__ACTION__* Container   __NEW_LINE__  `,
	containerExit0:   `- 📦 *__ACTION__* Container   __NEW_LINE__ __TAB__ -  ✔️ Exit code: *__EXIT_CODE__*   __NEW_LINE__   `,
	containerDie:     `- 📦 *__ACTION__* Container   __NEW_LINE__ __TAB__ -  ❗ Exit code: *__EXIT_CODE__*   __NEW_LINE__   `,
	containerKill:    `- 📦 *__ACTION__* Container   __NEW_LINE__ __TAB__ -  📣 Signal: *__SIGNAL__*    __NEW_LINE__   `,

	volumeMount:   `- 🛢️ *__ACTION__* Volume   __NEW_LINE__ __TAB__-  Mount point: __VOLUME_DESTINATION__    __NEW_LINE__ `,
	volumeUnmount: `- 🛢️ *__ACTION__* Volume   __NEW_LINE__ `,

	networkDefault: `- 🔗 *__ACTION__* Network ` + "`__NETWORK_NAME__`" + `  __NEW_LINE__`,
}
