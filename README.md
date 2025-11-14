# Go MCP Server

MCP(Model Context Protocol) 서버 구현을 위한 Go 기반 프로젝트입니다.

## 아키텍처

**Layered Architecture** 패턴을 사용합니다. 자세한 내용은 [ARCHITECTURE.md](./ARCHITECTURE.md)를 참조하세요.

```
Handler → Service → Repository → Model
```

## 디렉토리 구조

```
go_mcp_server/
├── cmd/server/              # 애플리케이션 진입점
├── internal/                # 비공개 애플리케이션 코드
│   ├── model/               # 데이터 모델
│   ├── repository/          # 데이터 접근 계층
│   ├── service/             # 비즈니스 로직 계층
│   ├── handler/             # HTTP 핸들러 계층
│   ├── router/              # 라우팅
│   ├── mcp/                 # MCP 서버 로직
│   └── infrastructure/      # 인프라 (Config, DB)
├── pkg/                     # 공용 라이브러리
└── test/                    # 테스트 유틸리티
```

## 빌드 및 실행

```bash
# 빌드
go build -o mcp-server ./cmd/server

# 실행
./mcp-server
```

## 환경 변수

`.env` 파일에 다음 변수를 설정하세요:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mcp_server
SERVER_PORT=8080
ALLOWED_ORIGINS=http://localhost
```

## API 엔드포인트

- `POST /mcp/tools/sync` - Tool 동기화
- `GET /mcp` - MCP SSE 연결
- `POST /mcp` - MCP 요청 처리

## 개발

### 테스트

```bash
go test ./...
```

### 새 기능 추가

1. `internal/model/` - 모델 정의
2. `internal/repository/` - 데이터 접근 인터페이스
3. `internal/service/` - 비즈니스 로직
4. `internal/handler/` - HTTP 핸들러
5. `internal/router/` - 라우트 등록
