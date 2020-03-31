package webhook

var webExLabels = map[int]string{
	startupTitle1: `🐳 Docker Version: **__DOCKER_VERSION__** __TAB__ApiVersion: __DOCKER_API_VERSION__ Node: ` + "`__NODE_NAME__`__NEW_LINE____NEW_LINE__",
	startupTitle2: `__NEW_LINE__🛠️ Os: __OS__  __TAB__Kernel: __KERNEL_VERSION__`,

	groupTitle:  `__NEW_LINE__🚛 ` + "`__NAME__`__NODE_NAME__" + ` __NEW_LINE__💿 __IMAGE__    __NEW_LINE__🕰️ __TIME__   __NEW_LINE__`,
	groupFooter: `__NEW_LINE__🆔 __ID__`,

	containerDefault: `>📦 **__ACTION__** Container   __NEW_LINE__  `,
	containerDie:     `>📦 **__ACTION__** Container   __NEW_LINE__   >>❗ Exit code: **__EXIT_CODE__**   __NEW_LINE__   ` ,
	containerKill:    `>📦 **__ACTION__** Container   __NEW_LINE__   >>📣 Signal: **__SIGNAL__**    __NEW_LINE__   `,

	volumeMount:   `>🛢️ **__ACTION__** Volume   __NEW_LINE__>>  Mount point: __VOLUME_DESTINATION__    __NEW_LINE__ `,
	volumeUnmount: `>🛢️ **__ACTION__** Volume   __NEW_LINE__ `,

	networkDefault: `>🔗 **__ACTION__** Network ` + "`__NETWORK_NAME__`" + `  __NEW_LINE__`,
}
