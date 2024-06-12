.PHONY: $(OUT_FILES) 

OUT_DIR := tmp
TPL_FILES = $(wildcard template/*.tpl.yaml)
OUT_FILES = $(patsubst template/%.tpl.yaml,$(OUT_DIR)/%.yaml, $(TPL_FILES))

build: buildmain

buildmain: $(OUT_FILES) 

$(OUT_FILES): cleanmain
	@mkdir -p $(OUT_DIR)
	$(eval TPL:=$(patsubst $(OUT_DIR)/%.yaml,template/%.tpl.yaml, $@))
	@echo $(SERVICE_GRPC_PORT)
	@echo generating $@ from $(TPL)...
	@sed -e "s~%ENVOY_PUBLIC_PORT%~${ENVOY_PUBLIC_PORT}~gi" \
    -e "s~%TLS_CERT_PATH%~${TLS_CERT_PATH}~gi" \
    -e "s~%TLS_KEY_PATH%~${TLS_KEY_PATH}~gi" \
    -e "s~%SERVICE_GRPC%~${SERVICE_GRPC}~gi" \
    -e "s~%SERVICE_GRPC_PORT%~${SERVICE_GRPC_PORT}~gi" \
    -e "s~%SERVICE_WEBSOCKET%~${SERVICE_WEBSOCKET}~gi" \
    -e "s~%SERVICE_WEBSOCKET_PORT%~${SERVICE_WEBSOCKET_PORT}~gi" \
    -e "s~%SERVICE_UI%~${SERVICE_UI}~gi" \
    -e "s~%SERVICE_UI_PORT%~${SERVICE_UI_PORT}~gi" $(TPL) > $@
	
clean: cleanmain

cleanmain:
	@echo cleaning up...
	@rm -rf $(OUT_DIR)
