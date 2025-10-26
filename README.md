# Go MCP Server

## 프로젝트 목표
Go MCP SDK 기반의 HTTP MCP 서버
- SSE (Server-Sent Events) 프로토콜로 MCP 통신
- DB 기반 Tool 동적 관리 (미래 확장)

---

## 아키텍처

```
HTTP Request
    ↓
Gin Router
    ↓
MCPHandler
    ↓
ServerManager (SSEHandler)
    ↓
MCP SDK
    ↓
Tool Handler
```

### 레이어 구조

- **MCP Layer** (`mcp/`): MCP Server, SSEHandler, Tool 등록
- **Service Layer** (`services/`): DB Tool 관리 (미래 확장)
- **Handler Layer** (`handlers/`): HTTP 엔드포인트
- **Model Layer** (`models/`): DB 모델 (미래 확장)

---

## 기술 스택

- Web Framework: **Gin**
- ORM: **GORM**
- Database: **PostgreSQL**
- MCP SDK: **Go MCP SDK v1.0.0**

---

## 디렉토리 구조

```
go_mcp_server/
├── mcp/               # MCP Server + SSEHandler
├── services/          # DB Tool 관리
├── handlers/          # HTTP 엔드포인트
├── models/            # DB 모델
├── database/          # DB 연결
├── config/            # 환경 설정
└── main.go            # 진입점
```

---

## 사용 방법

### 서버 실행
```bash
go run main.go
```

### MCP 프로토콜

**1. 세션 생성 (GET)**
```bash
curl -N -H "Accept: text/event-stream" http://localhost:8080/mcp
```

**2. Tool 호출 (POST)**
```bash
curl -X POST "http://localhost:8080/mcp?sessionid=<SESSION_ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "greet",
      "arguments": {"name": "World"}
    }
  }'
```

---

## 작업 상태

### 완료
- MCP Server 및 SSEHandler 구현
- HTTP 엔드포인트
- `greet` tool 구현
- DB 모델 및 연결

### 미래 확장
- DB 기반 동적 Tool 등록
- Tool CRUD API
- `echo` / `api_call` 타입 Tool
- DB Pool 사용 여부 확인