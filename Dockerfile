# ========================================
# Stage 1: Build
# ========================================
FROM golang:1.25.2-alpine AS builder

# 빌드에 필요한 도구 설치
RUN apk add --no-cache git ca-certificates tzdata

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 먼저 복사 (캐싱 최적화)
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드 (크기 최적화)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/mcp-server \
    .

# ========================================
# Stage 2: Runtime
# ========================================
FROM alpine:latest

# 보안 업데이트 및 필수 패키지 설치
RUN apk --no-cache add ca-certificates tzdata

# 비 root 유저 생성
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 작업 디렉토리 설정
WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/mcp-server .

# 환경변수 설정 (필요시 .env 파일이나 AWS Secrets Manager 사용)
# ENV 변수들은 런타임에 주입하는 것이 보안상 좋습니다

# 비 root 유저로 실행
RUN chown -R appuser:appuser /app
USER appuser

# 헬스체크 추가 (선택사항)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 포트 노출 (기본값, 실제 포트는 환경변수로 설정)
EXPOSE 8080

# 실행
CMD ["./mcp-server"]
