SPK_NAME = npc
SPK_VERS = $(VERSION)
SPK_REV = 1
SPK_ICON = src/npc.png

DEPENDS = cross/busybox cross/$(SPK_NAME)
SPK_CONFLICT = npc
#UNSUPPORTED_ARCHES = $(PPC_ARCHES)

MAINTAINER = cnlh
DESCRIPTION = nps是一款轻量级、高性能、功能强大的内网穿透代理服务器。目前支持tcp、udp流量转发，可支持任何tcp、udp上层协议（访问内网网站、本地支付接口调试、ssh访问、远程桌面，内网dns解析等等……），此外还支持内网http代理、内网socks5代理、p2p等，并带有功能强大的web管理端。
DISPLAY_NAME = npc
HELPURL = https://github.com/cnlh/nps
#LICENSE = ISC
#CHANGELOG = "Update to 24.0"

STARTABLE = yes
#BETA = 1

WIZARDS_DIR = src/wizard/
DSM_UI_DIR = ui
#CONF_DIR = src/conf
#FWPORTS = src/${SPK_NAME}.sc
SERVICE_USER = auto
SSS_SCRIPT = src/service-setup.sh
#SERVICE_PORT = 53
#SERVICE_PORT_TITLE = NPC
#NO_SERVICE_SHORTCUT = true
#SERVICE_PORT_ALL_USERS = false
#STARTSTOP_RESTART_SERVICE = cron
SERVICE_TYPE = legacy
SERVICE_URL = /webman/3rdparty/npc/index.cgi
INSTALL_PREFIX = /usr/local/$(SPK_NAME)
PRE_COPY_TARGET = dnscrypt_extra_install
BUSYBOX_CONFIG = usrmng daemon
ENV += BUSYBOX_CONFIG="$(BUSYBOX_CONFIG)"
include ../../mk/spksrc.spk.mk

ifeq ($(findstring $(ARCH),$(ARM5_ARCHES)),$(ARCH))
GO_ARCH = arm
ENV += GOARM=5
endif
ifeq ($(findstring $(ARCH),$(ARM7_ARCHES)),$(ARCH))
GO_ARCH = arm
ENV += GOARM=7
endif
ifeq ($(findstring $(ARCH),$(ARM8_ARCHES)),$(ARCH))
GO_ARCH = arm64
endif
ifeq ($(findstring $(ARCH),$(x86_ARCHES)),$(ARCH))
GO_ARCH = 386
endif
ifeq ($(findstring $(ARCH),$(x64_ARCHES)),$(ARCH))
GO_ARCH = amd64
endif
# Go has no ppc32 support for qoriq https://github.com/golang/go/issues/22885
#ifeq ($(findstring $(ARCH),$(PPC_ARCHES)),$(ARCH))
#GO_ARCH = ppc32
#endif
ifeq ($(GO_ARCH),)
$(error Unsupported ARCH $(ARCH))
endif

CONF_ARGS = -ldflags="-s -w"
ENV += GOPATH="$(WORK_DIR)" \
	GOOS=linux \
	GOARCH="$(GO_ARCH)" \
	CGO_ENABLED=0 \
	PATH=$(WORK_DIR)/../../../native/go/work-native/go/bin/:$$PATH

.PHONY: dnscrypt_extra_install
dnscrypt_extra_install: copy_msg
	@echo "===> Copying files"
	mkdir -p $(STAGING_DIR)/conf/
	cp -rf $(CURDIR)/src/conf/* $(STAGING_DIR)/conf/
#	install -m 755 $(CURDIR)/src/update-blocklist.sh $(STAGING_DIR)/blocklist/update-blocklist.sh
	@echo "===> UI - Copying files"
	mkdir -p $(STAGING_DIR)/$(DSM_UI_DIR)/
	cp -rf $(CURDIR)/src/$(DSM_UI_DIR)/* $(STAGING_DIR)/$(DSM_UI_DIR)/
	@echo "===> UI - Compiling cgi"
	env $(ENV) go build $(CONF_ARGS) -o $(STAGING_DIR)/$(DSM_UI_DIR)/index.cgi $(CURDIR)/src/$(DSM_UI_DIR)/cgi.go
