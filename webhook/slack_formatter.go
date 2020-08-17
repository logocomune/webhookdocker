package webhook

var slackLabels = map[int]string{
	appTitle:      `:heart: Webhook Docker __APP_VERSION__ (Built: __APP_BUILT_DATE__) __NEW_LINE__:bug: <__GITHUB_URL__|GitHub Issues> :ship: <__DOCKER_HUB_URL__|Docker Hub> __NEW_LINE__`,
	startupTitle1: `:whale: Docker Version: *__DOCKER_VERSION__* __TAB__ApiVersion: __DOCKER_API_VERSION__ Node: ` + "`__NODE_NAME__`__NEW_LINE__",
	startupTitle2: `:hammer_and_wrench: Os: __OS__  __TAB__Kernel: __KERNEL_VERSION__`,

	groupTitle:             `:articulated_lorry: ` + "`__NAME__`__NODE_NAME_W_SPACES__" + ` __NEW_LINE__:cd: __IMAGE__ __NEW_LINE__:mantelpiece_clock: __TIME__ __NEW_LINE__`,
	groupFooter:            `:id: __ID__`,
	groupFooterWithInspect: `:id: <__INSPECT_URL__|__ID__ :mag:>`,

	containerDefault: `>:package: *__ACTION__* Container __NEW_LINE__`,
	containerExit0:   `>:package: *__ACTION__* Container __NEW_LINE__> __TAB__ :heavy_check_mark: Exit code: *__EXIT_CODE__* __NEW_LINE__`,
	containerDie:     `>:package: *__ACTION__* Container __NEW_LINE__> __TAB__ :exclamation: Exit code: *__EXIT_CODE__* __NEW_LINE__`,
	containerKill:    `>:package: *__ACTION__* Container __NEW_LINE__> __TAB__ :mega: Signal: *__SIGNAL__* __NEW_LINE__`,

	volumeMount:   `>:oil_drum: *__ACTION__* Volume __NEW_LINE__> __TAB__   Mount point: __VOLUME_DESTINATION__ __NEW_LINE__`,
	volumeUnmount: `>:oil_drum: *__ACTION__* Volume __NEW_LINE__`,

	networkDefault: `>:link: *__ACTION__* Network ` + "`__NETWORK_NAME__`" + `__NEW_LINE__`,
}
