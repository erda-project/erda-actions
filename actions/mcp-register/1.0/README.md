# MCP Register

## 简介

提供 MCP Server 注册到 MCP Proxy 的能力

## params

### mcp_proxy_access_key

MCP Proxy的Access Key

### mcp_proxy_url

MCP Proxy的访问地址

### release_id

制品ID

### service_info

dice Action 部署后返回的 服务信息，包含 k8s service 地址

## 例子

```yaml
  - stage:
      - mcp-register:
          alias: mcp-register
          description: 用于将服务作为 MCP Server 注册至 MCP Proxy
          version: "1.0"
          params:
            mcp_proxy_access_key: ${your access key}
            mcp_proxy_url: https://mcp-proxy.daily.terminus.io
            release_id: ${release:OUTPUT:releaseID}
            service_info: ${dice:OUTPUT:serviceInfo}
```
