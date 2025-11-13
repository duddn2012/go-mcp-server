# Go MCP Server

DB 기반 MCP Tool 관리 및 동기화를 지원하는 Go HTTP 서버

## 아키텍처

```
HTTP Request
    ↓
Gin Router
    ↓
MCPHandler (Origin 검증)
    ↓
MCPService (Tool 동기화)
    ↓
ServerManager (Tool 등록/제거)
    ↓
MCP SDK (Streamable HTTP Handler)
```

### 레이어 구조

- **handlers/**: HTTP 엔드포인트 및 Origin 검증
- **services/**: DB Tool 조회 및 동기화 로직
- **mcp/**: MCP Server 및 Tool 관리
- **models/**: DB 모델 (Tool)
- **database/**: PostgreSQL 연결
- **config/**: 환경 설정
- **utils/**: JSON Schema 변환 유틸리티

## 기술 스택

- **Gin**: Web Framework
- **GORM**: ORM
- **PostgreSQL**: Database
- **MCP SDK**: Go MCP SDK v1.0.0

## 주요 기능

- DB에 저장된 Tool을 동적으로 로드 및 등록
- Tool 동기화 API (`/mcp/tools/sync`)
- MCP 프로토콜 지원 (Streamable HTTP)
- Origin 기반 CORS 검증


## 빌드
os 및 cpu 아키텍쳐는 배포 환경에 따라 가변적
MAC - GOOS=darwin GOARCH=amd64 go build -o common-mcp-server
Linux - GOOS=linux GOARCH=amd64 go build -o common-mcp-server

## 실행 방법
- 코드 기반
```bash
go run main.go
```

- 바이너리 실행
```bash
./common-mcp-server
```