function "tag" {
  params = [service]
  result = ["${APP_NAMESPACE}/${service}:${APP_VERSION}"]
}

function "tagdev" {
  params = [service]
  result = ["${APP_NAMESPACE}/${service}dev:${APP_VERSION}"]
}

group "dev" {
  targets = [
    "envoy",
    "redis",
    "uidev",
    "rpcdev",
    "websocketdev",
    "protobuilddev",
    "prototsbuilddev",
    "detectdev",
  ]
}

group "prod" {
  targets = [
    "envoy",
    "redis",
    "ui",
    "rpc",
    "websocket",
    "detect",
  ]
}

group "default" {
  targets = []
}

target "appargs" {
  args = {
    APP_NAMESPACE = APP_NAMESPACE
    APP_VERSION   = APP_VERSION
    BUILD_TIME    = timestamp()
  }
}

target "vcsargs" {
  args = {
    GITBRANCH = GITBRANCH
    GITCOMMIT = GITCOMMIT
    GITSHORT  = GITSHORT
  }
}

target "protobase" {
  inherits   = ["appargs", "vcsargs"]
  tags       = tag(SERVICE_PROTOBASE)
  context    = "."
  dockerfile = "build/base/${SERVICE_PROTOBASE}/Dockerfile"
  args = {
    ALPINE_VERSION            = ALPINE_VERSION
    GO_VERSION                = GO_VERSION
    PROTO_VERSION             = PROTO_VERSION
    PROTO_GEN_GO_VERSION      = PROTO_GEN_GO_VERSION
    PROTO_GEN_GO_GRPC_VERSION = PROTO_GEN_GO_GRPC_VERSION
  }
}

target "gobase" {
  inherits = ["appargs", "vcsargs"]
  tags     = tag(SERVICE_GOBASE)
  context  = "."
  contexts = {
    "${SERVICE_PROTOBASE}" = "target:${SERVICE_PROTOBASE}"
  }
  dockerfile      = "build/base/${SERVICE_GOBASE}/Dockerfile"
  no-cahce-filter = ["buildrpc", "buildwebsocket"]
  args            = {}
}

target "uibase" {
  inherits   = ["appargs"]
  tags       = tag(SERVICE_UIBASE)
  context    = "."
  dockerfile = "build/base/${SERVICE_UIBASE}/Dockerfile"
  args = {
    NODE_VERSION   = NODE_VERSION
    ALPINE_VERSION = ALPINE_VERSION
  }
}

target "rpc" {
  inherits = ["appargs"]
  tags     = tag(SERVICE_GRPC)
  context  = "."
  contexts = {
    "${SERVICE_GOBASE}" = "target:${SERVICE_GOBASE}"
  }
  dockerfile = "build/service/${SERVICE_GRPC}/Dockerfile"
  args = {
    SERVICE_GRPC_PORT = SERVICE_GRPC_PORT
  }
}

target "websocket" {
  inherits = ["appargs"]
  tags     = tag(SERVICE_WEBSOCKET)
  context  = "."
  contexts = {
    "${SERVICE_GOBASE}" = "target:${SERVICE_GOBASE}"
  }
  dockerfile = "build/service/${SERVICE_WEBSOCKET}/Dockerfile"
  args = {
    SERVICE_WEBSOCKET_PORT = SERVICE_WEBSOCKET_PORT
  }
}

target "detect" {
  inherits = ["appargs"]
  tags     = tag(SERVICE_DETECT)
  context  = "."
  contexts = {
    "${SERVICE_GOBASE}" = "target:${SERVICE_GOBASE}"
  }
  dockerfile = "build/service/${SERVICE_DETECT}/Dockerfile"
}

target "ui" {
  inherits = ["appargs"]
  tags     = tag(SERVICE_UI)
  context  = "."
  contexts = {
    "${SERVICE_UIBASE}" = "target:${SERVICE_UIBASE}"
  }
  dockerfile = "build/service/${SERVICE_UI}/Dockerfile"
  args = {
    ALPINE_VERSION  = ALPINE_VERSION
    NGINX_VERSION   = NGINX_VERSION
    SERVICE_UI_PORT = SERVICE_UI_PORT
  }
}

target "redis" {
  tags       = tag(SERVICE_REDIS)
  context    = "."
  dockerfile = "build/service/${SERVICE_REDIS}/Dockerfile"
  args = {
    REDIS_VERSION  = REDIS_VERSION
    ALPINE_VERSION = ALPINE_VERSION
  }
}

target "envoy" {
  tags       = tag(SERVICE_ENVOY)
  context    = "."
  dockerfile = "build/service/${SERVICE_ENVOY}/Dockerfile"
  args = {
    ALPINE_VERSION         = ALPINE_VERSION
    ENVOY_ADMIN_PORT       = ENVOY_ADMIN_PORT
    ENVOY_PUBLIC_PORT      = ENVOY_PUBLIC_PORT
    ENVOY_VERSION          = ENVOY_VERSION
    TLS_CERT_PATH          = TLS_CERT_PATH
    TLS_KEY_PATH           = TLS_KEY_PATH
    SERVICE_GRPC           = SERVICE_GRPC
    SERVICE_GRPC_PORT      = SERVICE_GRPC_PORT
    SERVICE_WEBSOCKET      = SERVICE_WEBSOCKET
    SERVICE_WEBSOCKET_PORT = SERVICE_WEBSOCKET_PORT
    SERVICE_UI             = SERVICE_UI
    SERVICE_UI_PORT        = SERVICE_UI_PORT
  }
}

target "detectdev" {
  tags       = tagdev(SERVICE_DETECT)
  context    = "."
  dockerfile = "build/dev/${SERVICE_DETECT}/Dockerfile"
  args = {
    ALPINE_VERSION = ALPINE_VERSION
    GO_VERSION     = GO_VERSION
  }
}

target "protobuilddev" {
  tags    = tagdev("protobuild")
  context = "."
  contexts = {
    "${SERVICE_PROTOBASE}" = "target:${SERVICE_PROTOBASE}"
  }
  dockerfile = "build/dev/protobuild/Dockerfile"
}

target "prototsbuilddev" {
  tags       = tagdev("prototsbuild")
  context    = "."
  dockerfile = "build/dev/prototsbuild/Dockerfile"
  args = {
    ALPINE_VERSION = ALPINE_VERSION
    NODE_VERSION   = NODE_VERSION
  }
}

target "rpcdev" {
  tags       = tagdev(SERVICE_GRPC)
  context    = "."
  dockerfile = "build/dev/${SERVICE_GRPC}/Dockerfile"
  args = {
    ALPINE_VERSION    = ALPINE_VERSION
    GO_VERSION        = GO_VERSION
    SERVICE_GRPC_PORT = SERVICE_GRPC_PORT
  }
}

target "uidev" {
  tags       = tagdev(SERVICE_UI)
  context    = "."
  dockerfile = "build/dev/${SERVICE_UI}/Dockerfile"
  args = {
    ALPINE_VERSION    = ALPINE_VERSION
    NODE_VERSION      = NODE_VERSION
    SERVICE_UI_PORT   = SERVICE_UI_PORT
    ENVOY_PUBLIC_PORT = ENVOY_PUBLIC_PORT
  }
}

target "websocketdev" {
  tags       = tagdev(SERVICE_WEBSOCKET)
  context    = "."
  dockerfile = "build/dev/${SERVICE_WEBSOCKET}/Dockerfile"
  args = {
    ALPINE_VERSION         = ALPINE_VERSION
    GO_VERSION             = GO_VERSION
    SERVICE_WEBSOCKET_PORT = SERVICE_WEBSOCKET_PORT
  }
}